// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

mod data;

use actix_web::{get, App, HttpServer, Responder};
use actix_web::{post, web, HttpResponse};
use actix_web::{route, HttpRequest};
use std::env;
use data::sensor_data_json_helper;
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

// This function is the handler for GET requests on the root path "/".
#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().json(serde_json::json!({
        "message": "actix api server is running"
    }))
}

// This function is the handler for POST requests on the "/echo" path.
#[post("/echo")]
async fn echo(req_body: web::Json<serde_json::Value>) -> HttpResponse {
    HttpResponse::Ok().json(req_body.into_inner())
}

// This function is the handler for POST requests on the "/log" path.
#[post("/log")]
async fn log(req_body: web::Json<serde_json::Value>) -> HttpResponse {
    let json_data = req_body.into_inner().to_string();
    let sensor_data_access = get_data_access();

    match sensor_data_access.log_sensor_data(&json_data).await {
        Ok(result) => {
            match result {
                Ok(_) => HttpResponse::Ok().json(serde_json::json!({
                    "message": "Data logged successfully"
                })),
                Err(e) => HttpResponse::InternalServerError()
                    .body(format!("Failed to log sensor data: {}", e)),
            }
        }
        Err(e) => HttpResponse::InternalServerError()
            .body(format!("Task join error: {}", e)),
    }
}

// This function is the handler for GET requests on the root path "/report".
#[get("/report")]
async fn report() -> impl Responder {
    let sensor_data_access = get_data_access();

    match sensor_data_access.fetch_sensor_data().await {
        Ok(result) => {
            match result {
                Ok(json_strings) => {
                    // Combine all JSON strings and convert to CSV
                    if json_strings.is_empty() {
                        return HttpResponse::Ok()
                            .content_type("text/csv")
                            .body("No data available");
                    }
                    // Convert first JSON string to CSV (or combine multiple)
                    let json_str = &json_strings[0];
                    match sensor_data_json_helper::json_to_csv(json_str) {
                        Ok(csv) => HttpResponse::Ok()
                            .content_type("text/csv")
                            .body(csv),
                        Err(e) => HttpResponse::InternalServerError()
                            .body(format!("Failed to convert to CSV: {}", e)),
                    }
                }
                Err(e) => HttpResponse::InternalServerError()
                    .body(format!("Failed to fetch sensor data: {}", e)),
            }
        }
        Err(e) => HttpResponse::InternalServerError()
            .body(format!("Task join error: {}", e)),
    }
}

// This function handles purge GET and POST requests on the "/purge" path.
#[route("/purge", method = "GET", method = "POST")]
async fn purge(req: HttpRequest, _body: web::Bytes) -> impl Responder {
    let sensor_data_access = get_data_access();

    match req.method() {
        &actix_web::http::Method::GET | &actix_web::http::Method::POST => {
            match sensor_data_access.purge_sensor_data().await {
                Ok(result) => {
                    match result {
                        Ok(_) => HttpResponse::Ok().json(serde_json::json!({
                            "message": "purged"
                        })),
                        Err(e) => HttpResponse::InternalServerError()
                            .body(format!("Failed to purge sensor data: {}", e)),
                    }
                }
                Err(e) => HttpResponse::InternalServerError()
                    .body(format!("Task join error: {}", e)),
            }
        }
        _ => HttpResponse::MethodNotAllowed().finish(),
    }
}

// The #[actix_web::main] macro sets up an async runtime for your main function.
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Create a new HttpServer.
    HttpServer::new(|| {
        // Create a new App instance and register the `hello` service.
        App::new().service(hello)
            .service(echo)
            .service(log)
            .service(report)
            .service(purge)
    })
    // Bind the server to the local address "127.0.0.1" and port 8080.
    .bind(("0.0.0.0", 8080))?
    // Start the server and wait for it to complete.
    .run()
    .await
}
