// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use crate::data::sensor_data_access_trait::SensorDataAccess;
use crate::data::sensor_data_json_helper::validate_sensor_json;
use std::error::Error;
use tokio::task;
use redis::AsyncCommands;

const REDIS_URL: &str = "redis://192.168.1.60/";
const REDIS_LIST_KEY: &str = "sensor_data";

async fn get_redis_connection() -> Result<redis::aio::MultiplexedConnection, Box<dyn Error + Send + Sync>> {
    let client = redis::Client::open(REDIS_URL)
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Redis client error: {}", e))) })?;
    client.get_multiplexed_async_connection()
        .await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Redis connection error: {}", e))) })
}

pub struct RedisDataAccess;

impl RedisDataAccess {
    pub fn new() -> Self {
        RedisDataAccess
    }
}

impl SensorDataAccess for RedisDataAccess {
    fn log_sensor_data<'a>(&'a self, json_data: &'a str) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        let json_owned = json_data.to_string();
        task::spawn(async move {
            let parsed = validate_sensor_json(&json_owned)
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, e)) })?;

            let cleaned_json = serde_json::to_string(&parsed)
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("JSON serialization error: {}", e))) })?;

            let mut con = get_redis_connection().await?;
            con.lpush::<_, _, ()>(REDIS_LIST_KEY, &cleaned_json).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Redis LPUSH error: {}", e))) })?;

            println!("Logging sensor data to Redis: {}", json_owned);
            Ok(())
        })
    }

    fn fetch_sensor_data(&self) -> task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Fetching sensor data from Redis");

            let mut con = get_redis_connection().await?;
            let values: Vec<String> = con.lrange(REDIS_LIST_KEY, 0, -1).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Redis LRANGE error: {}", e))) })?;

            Ok(values)
        })
    }

    fn purge_sensor_data(&self) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Purging sensor data from Redis");

            let mut con = get_redis_connection().await?;
            con.del::<_, ()>(REDIS_LIST_KEY).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Redis DEL error: {}", e))) })?;

            println!("Redis sensor data purged successfully.");
            Ok(())
        })
    }
}
