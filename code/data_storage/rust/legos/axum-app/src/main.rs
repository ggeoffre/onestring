// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

mod data;

use axum::{Router, body::Bytes, response::IntoResponse};
use axum::response::Response;
use axum::http::{header, HeaderValue, StatusCode};
use std::net::SocketAddr;
use tokio::net::TcpListener;
use std::env;
use data::sensor_data_json_helper::{json_to_csv};
use data::redis_data_access::RedisDataAccess;
use data::sensor_data_access_trait::SensorDataAccess;
use data::mongo_data_access::MongoDataAccess;
use data::cassandra_data_access::CassandraDataAccess;
use data::mysql_data_access::MySQLDataAccess;
use data::postgres_data_access::PostgresDataAccess;

fn get_data_access() -> Box<dyn SensorDataAccess + Send> {
    let data_access_type = env::var("DATA_ACCESS").unwrap_or_else(|_| "redis".to_string());

    match data_access_type.as_str() {
        "redis" => Box::new(RedisDataAccess::new()),
        "mongo" => Box::new(MongoDataAccess::new()),
        "cassandra" => Box::new(CassandraDataAccess::new()),
        "mysql" => Box::new(MySQLDataAccess::new()),
        "postgres" => Box::new(PostgresDataAccess::new()),
        _ => panic!("Unsupported DATA_ACCESS type: {}", data_access_type),
    }
}

#[tokio::main]
async fn main() {
    // Build our application with the external handler function
    let app = Router::new()
        .route("/", axum::routing::get(root_handler))
        .route("/echo", axum::routing::post(echo_handler))
        .route("/log", axum::routing::post(log_handler))
        .route("/report", axum::routing::get(report_handler))
        .route("/purge", axum::routing::post(purge_handler))
        .route("/purge", axum::routing::get(purge_handler));

    // Listen on a specified address
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    let listener = TcpListener::bind(&addr).await.unwrap();

    println!("Listening on http://{}", addr);

    // Serve the application
    axum::serve(listener, app).await.unwrap();
}

// The external handler function
pub async fn root_handler() -> impl IntoResponse {
    axum::response::Json(serde_json::json!({"message": "axum api server is running"}))
}

pub async fn echo_handler(body: Bytes) -> impl IntoResponse {
    let json_result = serde_json::from_slice::<serde_json::Value>(&body);
    if let Ok(json) = &json_result {
        println!("{}", json);
    }
    match json_result {
        Ok(json) => axum::response::Json(json),
        Err(_) => axum::response::Json(serde_json::json!({"error": "Invalid JSON"})),
    }
}

pub async fn log_handler(body: Bytes) -> impl IntoResponse {
    let json_result = serde_json::from_slice::<serde_json::Value>(&body);
    match json_result {
        Ok(json) => {
            let json_data = json.to_string();
            let sensor_data_access = get_data_access();

            match sensor_data_access.log_sensor_data(&json_data).await {
                Ok(result) => {
                    match result {
                        Ok(_) => axum::response::Json(serde_json::json!({
                            "message": "Data logged successfully"
                        })),
                        Err(e) => axum::response::Json(serde_json::json!({
                            "error": format!("Failed to log sensor data: {}", e)
                        })),
                    }
                }
                Err(e) => axum::response::Json(serde_json::json!({
                    "error": format!("Task join error: {}", e)
                })),
            }
        }
        Err(_) => axum::response::Json(serde_json::json!({"error": "Invalid JSON"})),
    }
}

pub async fn report_handler() -> Response {
    let sensor_data_access = get_data_access();

    match sensor_data_access.fetch_sensor_data().await {
        Ok(result) => {
            match result {
                Ok(json_strings) => {
                    if json_strings.is_empty() {
                        return Response::builder()
                            .status(StatusCode::OK)
                            .header(header::CONTENT_TYPE, "text/csv")
                            .body(axum::body::Body::from("No data available"))
                            .unwrap();
                    }
                    let json_str = &json_strings[0];
                    match json_to_csv(json_str) {
                        Ok(csv) => Response::builder()
                            .status(StatusCode::OK)
                            .header(header::CONTENT_TYPE, "text/csv")
                            .header(
                                header::CONTENT_DISPOSITION,
                                HeaderValue::from_str("attachment; filename=\"report.csv\"").unwrap(),
                            )
                            .body(axum::body::Body::from(csv))
                            .unwrap(),
                        Err(e) => Response::builder()
                            .status(StatusCode::INTERNAL_SERVER_ERROR)
                            .body(axum::body::Body::from(format!("Failed to convert to CSV: {}", e)))
                            .unwrap(),
                    }
                }
                Err(e) => Response::builder()
                    .status(StatusCode::INTERNAL_SERVER_ERROR)
                    .body(axum::body::Body::from(format!("Failed to fetch sensor data: {}", e)))
                    .unwrap(),
            }
        }
        Err(e) => Response::builder()
            .status(StatusCode::INTERNAL_SERVER_ERROR)
            .body(axum::body::Body::from(format!("Task join error: {}", e)))
            .unwrap(),
    }
}

pub async fn purge_handler() -> impl IntoResponse {
    let sensor_data_access = get_data_access();

    match sensor_data_access.purge_sensor_data().await {
        Ok(result) => {
            match result {
                Ok(_) => axum::response::Json(serde_json::json!({
                    "message": "purged"
                })),
                Err(e) => axum::response::Json(serde_json::json!({
                    "error": format!("Failed to purge sensor data: {}", e)
                })),
            }
        }
        Err(e) => axum::response::Json(serde_json::json!({
            "error": format!("Task join error: {}", e)
        })),
    }
}
