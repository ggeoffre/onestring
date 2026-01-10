// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use actix_web::{get, App, HttpServer, Responder};
use actix_web::{post, web, HttpResponse};
use actix_web::{route, HttpRequest};

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
    HttpResponse::Ok().json(req_body.into_inner())
}

// This function is the handler for GET requests on the root path "/report".
#[get("/report")]
async fn report() -> impl Responder {
    let cvs_string = match json_to_csv(SENSOR_DATA) {
        Ok(csv) => csv,
        Err(err) => {
            return HttpResponse::InternalServerError()
                .body(format!("Failed to generate CSV: {}", err));
        }
    };
    HttpResponse::Ok()
        .content_type("text/csv")
        .body(cvs_string)
}

// This function handles purge GET and POST requests on the "/purge" path.
#[route("/purge", method = "GET", method = "POST")]
async fn purge(req: HttpRequest, _body: web::Bytes) -> impl Responder {
    match req.method() {
        &actix_web::http::Method::GET => HttpResponse::Ok().json(serde_json::json!({
            "message": "purged"
        })),
        &actix_web::http::Method::POST => HttpResponse::Ok().json(serde_json::json!({
            "message": "purged"
        })),
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
