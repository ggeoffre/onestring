// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use redis::AsyncCommands;
use serde_json::Value;
use log::{info};
use crate::sensor_data_json_helper::{validate_sensor_json};

/// Get a Redis async connection.
async fn get_redis_connection() -> Result<redis::aio::MultiplexedConnection, String> {
    let client = redis::Client::open("redis://192.168.1.60/")
        .map_err(|e| format!("Redis client error: {}", e))?;
    client.get_multiplexed_async_connection()
        .await
        .map_err(|e| format!("Redis connection error: {}", e))
}

/// Set a JSON string at a key (overwrites any existing value).
pub async fn set_json(key: &str, json_str: &str) -> Result<(), String> {
    let parsed = validate_sensor_json(json_str)?;
    let cleaned_json = serde_json::to_string(&parsed)
        .map_err(|e| format!("JSON serialization error: {}", e))?;
    let mut con = get_redis_connection().await?;
    con.set::<_, _, ()>(key, cleaned_json).await
        .map_err(|e| format!("Redis SET error: {}", e))?;
    info!("Set JSON at key '{}'", key);
    Ok(())
}

/// Get a JSON string from a key.
pub async fn get_json(key: &str) -> Result<String, String> {
    let mut con = get_redis_connection().await?;
    let value: String = con.get(key).await
        .map_err(|e| format!("Redis GET error: {}", e))?;
    Ok(value)
}

/// Delete a key (row-level delete).
pub async fn delete_key(key: &str) -> Result<(), String> {
    let mut con = get_redis_connection().await?;
    con.del::<_, ()>(key).await
        .map_err(|e| format!("Redis DEL error: {}", e))?;
    info!("Deleted key '{}'", key);
    Ok(())
}

/// Push a JSON string to a Redis list.
pub async fn lpush_json(key: &str, json_str: &str) -> Result<(), String> {
    let parsed = validate_sensor_json(json_str)?;
    let cleaned_json = serde_json::to_string(&parsed)
        .map_err(|e| format!("JSON serialization error: {}", e))?;
    let mut con = get_redis_connection().await?;
    con.lpush::<_, _, ()>(key, cleaned_json).await
        .map_err(|e| format!("Redis LPUSH error: {}", e))?;
    info!("LPUSH JSON to list '{}'", key);
    Ok(())
}

/// Get all JSON strings from a Redis list as a Vec<String>.
pub async fn lrange_json(key: &str) -> Result<Vec<String>, String> {
    let mut con = get_redis_connection().await?;
    let values: Vec<String> = con.lrange(key, 0, -1).await
        .map_err(|e| format!("Redis LRANGE error: {}", e))?;
    Ok(values)
}

/// Add a JSON string to a Redis set.
pub async fn sadd_json(key: &str, json_str: &str) -> Result<(), String> {
    let parsed = validate_sensor_json(json_str)?;
    let cleaned_json = serde_json::to_string(&parsed)
        .map_err(|e| format!("JSON serialization error: {}", e))?;
    let mut con = get_redis_connection().await?;
    con.sadd::<_, _, ()>(key, cleaned_json).await
        .map_err(|e| format!("Redis SADD error: {}", e))?;
    info!("SADD JSON to set '{}'", key);
    Ok(())
}

/// Get all JSON strings from a Redis set as a Vec<String>.
pub async fn smembers_json(key: &str) -> Result<Vec<String>, String> {
    let mut con = get_redis_connection().await?;
    let members: Vec<String> = con.smembers(key).await
        .map_err(|e| format!("Redis SMEMBERS error: {}", e))?;
    Ok(members)
}

/// Convert a Vec of JSON strings to CSV (header inferred from first object).
pub fn json_array_to_csv(json_array: &[String]) -> Result<String, String> {
    if json_array.is_empty() {
        return Ok(String::new());
    }
    let mut csv_output = String::new();
    let mut headers_written = false;
    for json_str in json_array {
        let parsed: Value = serde_json::from_str(json_str)
            .map_err(|e| format!("JSON parse error: {}", e))?;
        if let Value::Object(map) = parsed {
            if !headers_written {
                let headers: Vec<&String> = map.keys().collect();
                csv_output.push_str(&headers.iter().map(|s| s.as_str()).collect::<Vec<&str>>().join(","));
                csv_output.push('\n');
                headers_written = true;
            }
            let row: Vec<String> = map.values()
                .map(|v| match v {
                    Value::String(s) => s.clone(),
                    _ => v.to_string(),
                })
                .collect();
            csv_output.push_str(&row.join(","));
            csv_output.push('\n');
        } else {
            return Err("Each JSON string must represent an object".to_string());
        }
    }
    Ok(csv_output)
}
