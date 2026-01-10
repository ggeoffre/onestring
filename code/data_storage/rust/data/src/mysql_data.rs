// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use sqlx::mysql::MySqlPoolOptions;
use sqlx::{MySqlPool, Row};
use log::info;
use crate::sensor_data_json_helper::validate_sensor_json;

const DATABASE_URL: &str = "mysql://root:@192.168.1.60:3306/sensor_data_db";

async fn get_pool() -> Result<MySqlPool, sqlx::Error> {
    MySqlPoolOptions::new().max_connections(5).connect(DATABASE_URL).await
}

pub async fn setup_database() -> Result<(), String> {
    let pool = get_pool().await.map_err(|e| format!("Pool error: {}", e))?;
    sqlx::query("CREATE DATABASE IF NOT EXISTS sensor_data_db").execute(&pool).await.map_err(|e| format!("Create DB error: {}", e))?;
    sqlx::query(r#"CREATE TABLE IF NOT EXISTS sensor_data (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        recorded BIGINT NOT NULL,
        location VARCHAR(255) NOT NULL,
        sensor VARCHAR(255) NOT NULL,
        measurement VARCHAR(255) NOT NULL,
        units VARCHAR(50) NOT NULL,
        value DOUBLE NOT NULL
    )"#).execute(&pool).await.map_err(|e| format!("Create table error: {}", e))?;
    info!("Database and table setup completed.");
    Ok(())
}

pub async fn insert_json_into_table(json_str: &str) -> Result<(), String> {
    let parsed = validate_sensor_json(json_str)?;
    let pool = get_pool().await.map_err(|e| format!("Pool error: {}", e))?;
    let query = "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)";
    sqlx::query(query)
        .bind(parsed["recorded"].as_i64().unwrap())
        .bind(parsed["location"].as_str().unwrap())
        .bind(parsed["sensor"].as_str().unwrap())
        .bind(parsed["measurement"].as_str().unwrap())
        .bind(parsed["units"].as_str().unwrap())
        .bind(parsed["value"].as_f64().unwrap())
        .execute(&pool).await.map_err(|e| format!("Insert error: {}", e))?;
    info!("Data inserted successfully.");
    Ok(())
}

pub async fn fetch_data_as_csv() -> Result<String, String> {
    let pool = get_pool().await.map_err(|e| format!("Pool error: {}", e))?;
    let rows = sqlx::query("SELECT recorded, location, sensor, measurement, units, CAST(value AS DOUBLE) as value FROM sensor_data")
        .fetch_all(&pool).await.map_err(|e| format!("Fetch error: {}", e))?;
    let mut csv_data = String::from("recorded,location,sensor,measurement,units,value\n");
    for row in rows {
        csv_data.push_str(&format!("{},{},{},{},{},{}\n",
            row.get::<i64, _>("recorded"),
            row.get::<String, _>("location"),
            row.get::<String, _>("sensor"),
            row.get::<String, _>("measurement"),
            row.get::<String, _>("units"),
            row.get::<f64, _>("value"),
        ));
    }
    Ok(csv_data)
}

pub async fn delete_row(recorded: i64, location: &str, sensor: &str) -> Result<(), String> {
    let pool = get_pool().await.map_err(|e| format!("Pool error: {}", e))?;
    let query = "DELETE FROM sensor_data WHERE recorded = ? AND location = ? AND sensor = ?";
    sqlx::query(query)
        .bind(recorded)
        .bind(location)
        .bind(sensor)
        .execute(&pool).await.map_err(|e| format!("Delete error: {}", e))?;
    info!("Row deleted successfully.");
    Ok(())
}
