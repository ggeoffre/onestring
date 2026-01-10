// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use crate::data::sensor_data_access_trait::SensorDataAccess;
use crate::data::sensor_data_json_helper::validate_sensor_json;
use std::error::Error;
use tokio::task;
use sqlx::postgres::PgPoolOptions;
use sqlx::{PgPool, Row};

const DATABASE_URL: &str = "postgres://postgres:@192.168.1.60:5432/sensor_data_db";

async fn get_pool() -> Result<PgPool, Box<dyn Error + Send + Sync>> {
    PgPoolOptions::new()
        .max_connections(5)
        .connect(DATABASE_URL)
        .await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Pool error: {}", e))) })
}

async fn setup_database() -> Result<(), Box<dyn Error + Send + Sync>> {
    let pool = get_pool().await?;
    let create_table_query = r#"
        CREATE TABLE IF NOT EXISTS sensor_data (
            id SERIAL PRIMARY KEY,
            recorded BIGINT NOT NULL,
            location VARCHAR NOT NULL,
            sensor VARCHAR NOT NULL,
            measurement VARCHAR NOT NULL,
            units VARCHAR NOT NULL,
            value DOUBLE PRECISION NOT NULL
        )
    "#;
    sqlx::query(create_table_query)
        .execute(&pool)
        .await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Create table error: {}", e))) })?;

    println!("PostgreSQL database and table setup completed.");
    Ok(())
}

pub struct PostgresDataAccess;

impl PostgresDataAccess {
    pub fn new() -> Self {
        PostgresDataAccess
    }
}

impl SensorDataAccess for PostgresDataAccess {
    fn log_sensor_data<'a>(&'a self, json_data: &'a str) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        let json_owned = json_data.to_string();
        task::spawn(async move {
            // Ensure database and table exist
            setup_database().await?;

            let parsed = validate_sensor_json(&json_owned)
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, e)) })?;

            // Handle recorded as either integer or string
            let recorded = parsed["recorded"].as_i64()
                .or_else(|| parsed["recorded"].as_str().and_then(|s| s.parse::<i64>().ok()))
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'recorded' field")) })?;

            let location = parsed["location"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'location' field")) })?;

            let sensor = parsed["sensor"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'sensor' field")) })?;

            let measurement = parsed["measurement"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'measurement' field")) })?;

            let units = parsed["units"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'units' field")) })?;

            // Handle value as either float or integer or string
            let value = parsed["value"].as_f64()
                .or_else(|| parsed["value"].as_i64().map(|i| i as f64))
                .or_else(|| parsed["value"].as_str().and_then(|s| s.parse::<f64>().ok()))
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'value' field")) })?;

            let pool = get_pool().await?;
            let query = "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES ($1, $2, $3, $4, $5, $6)";
            sqlx::query(query)
                .bind(recorded)
                .bind(location)
                .bind(sensor)
                .bind(measurement)
                .bind(units)
                .bind(value)
                .execute(&pool)
                .await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Insert error: {}", e))) })?;

            println!("Logging sensor data to Postgres: {}", json_owned);
            Ok(())
        })
    }

    fn fetch_sensor_data(&self) -> task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Fetching sensor data from Postgres");

            let pool = get_pool().await?;
            let rows = sqlx::query("SELECT recorded, location, sensor, measurement, units, CAST(value AS DOUBLE PRECISION) as value FROM sensor_data")
                .fetch_all(&pool)
                .await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Fetch error: {}", e))) })?;

            let mut json_strings: Vec<String> = Vec::new();
            for row in rows {
                let recorded: i64 = row.get("recorded");
                let location: String = row.get("location");
                let sensor: String = row.get("sensor");
                let measurement: String = row.get("measurement");
                let units: String = row.get("units");
                let value: f64 = row.get("value");

                let json_str = format!(
                    r#"{{"recorded":{},"location":"{}","sensor":"{}","measurement":"{}","units":"{}","value":{}}}"#,
                    recorded, location, sensor, measurement, units, value
                );
                json_strings.push(json_str);
            }

            Ok(json_strings)
        })
    }

    fn purge_sensor_data(&self) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Purging sensor data from Postgres");

            let pool = get_pool().await?;
            sqlx::query("DELETE FROM sensor_data")
                .execute(&pool)
                .await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Delete error: {}", e))) })?;

            println!("Postgres sensor data purged successfully.");
            Ok(())
        })
    }
}
