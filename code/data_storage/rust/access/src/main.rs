// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

mod data_access;

use data_access::redis_data_access::RedisDataAccess;
use data_access::sensor_data_access_trait::SensorDataAccess;
use data_access::mongo_data_access::MongoDataAccess;
use data_access::cassandra_data_access::CassandraDataAccess;
use data_access::mysql_data_access::MySQLDataAccess;
use data_access::postgres_data_access::PostgresDataAccess;
use std::env;

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

    // Dynamically get the data access implementation
    let sensor_data_access = get_data_access();

    // Use the SensorDataAccess trait methods with error handling
    match sensor_data_access
        .log_sensor_data("{\"sensor\":\"temperature\",\"value\":22.3}")
        .await
    {
        Ok(_) => {},
        Err(e) => eprintln!("Failed to log sensor data: {}", e),
    }

    let vec: Vec<String>;
    match sensor_data_access.fetch_sensor_data().await {
        Ok(data) => {
            vec = data.expect("REASON");
        }
        Err(e) => {
            eprintln!("Failed to fetch sensor data: {}", e);
            vec = Vec::new();
        },
    }
    println!("Fetched sensor data (vec): {:?}", vec);

    match sensor_data_access.purge_sensor_data().await {
        Ok(_) => {},
        Err(e) => eprintln!("Failed to purge sensor data: {}", e),
    }
}
