// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use axum::{Router, body::Bytes, response::IntoResponse};
use axum::response::Response;
use axum::http::{header, HeaderValue, StatusCode};
use std::net::SocketAddr;
use tokio::net::TcpListener;

const SENSOR_DATA: &str = r#"{
    "recorded": 1756655999,
    "location": "den",
    "sensor": "bmp280",
    "measurement": "temperature",
    "units": "C",
    "value": 22.3
}"#;

fn json_to_csv(json_str: &str) -> Result<String, Box<dyn std::error::Error>> {
    let json_value: serde_json::Value = serde_json::from_str(json_str)?;

    if let serde_json::Value::Object(map) = json_value {
        let mut csv_string = String::new();

        // Write the header row using the keys
        let headers: Vec<&str> = map.keys().map(|k| k.as_str()).collect();
        csv_string.push_str(&headers.join(","));
        csv_string.push('\n');

        // Write the values row
        let values: Vec<String> = map.values().map(|v| v.to_string()).collect();
        csv_string.push_str(&values.join(","));
        csv_string.push('\n');

        Ok(csv_string)
    } else {
        Err("Input JSON must be an object".into())
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
    if let Ok(json) = &json_result {
        println!("{}", json);
    }
    match json_result {
        Ok(json) => axum::response::Json(json),
        Err(_) => axum::response::Json(serde_json::json!({"error": "Invalid JSON"})),
    }
}

pub async fn report_handler() -> impl IntoResponse {
    let csv_data = match json_to_csv(SENSOR_DATA) {
        Ok(data) => data,
        Err(err) => {
            eprintln!("Error converting JSON to CSV: {}", err);
            return (StatusCode::INTERNAL_SERVER_ERROR, "Failed to generate report").into_response();
        }
    };
    let response: Response<axum::body::Body> = Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "text/csv")
        .header(
            header::CONTENT_DISPOSITION,
            HeaderValue::from_str("attachment; filename=\"report.csv\"").unwrap(),
        )
        .body(axum::body::Body::from(csv_data))
        .unwrap();

    response.into_response()
}

pub async fn purge_handler() -> impl IntoResponse {
    axum::response::Json(serde_json::json!({"message": "Data purged"}))
}
