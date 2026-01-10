// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

mod cassandra_data;
mod mongo_data;
mod mysql_data;
mod postgres_data;
mod redis_data;
mod sensor_data_json_helper;

use sensor_data_json_helper::create_sample_json;
use log::{info, error};

#[tokio::main]
async fn main() {
    env_logger::init();

    // --- Cassandra ---
    info!("--- Cassandra Workflow ---");
    if let Err(e) = cassandra_data::create_keyspace_and_table().await {
        error!("Cassandra setup error: {}", e);
    } else {
        info!("Cassandra keyspace and table created successfully.");
    }
    let cass_json = create_sample_json();
    if let Err(e) = cassandra_data::insert_json_row(&cass_json).await {
        error!("Cassandra insert error: {}", e);
    } else {
        info!("Cassandra row inserted successfully.");
    }
    match cassandra_data::fetch_table_as_csv().await {
        Ok(csv) => info!("Cassandra CSV:\n{}", csv),
        Err(e) => error!("Cassandra fetch error: {}", e),
    }
    if let Err(e) = cassandra_data::delete_row("den", 1756655999, "bmp280").await {
        error!("Cassandra delete error: {}", e);
    } else {
        info!("Cassandra row deleted successfully.");
    }

    // --- MongoDB ---
    info!("--- MongoDB Workflow ---");
    let mongo_json = create_sample_json();
    if let Err(e) = mongo_data::insert_json_into_collection(&mongo_json).await {
        error!("MongoDB insert error: {}", e);
    } else {
        info!("MongoDB document inserted successfully.");
    }
    match mongo_data::fetch_collection_as_csv().await {
        Ok(csv) => info!("MongoDB CSV:\n{}", csv),
        Err(e) => error!("MongoDB fetch error: {}", e),
    }
    if let Err(e) = mongo_data::delete_row(1756655999, "den", "bmp280").await {
        error!("MongoDB delete error: {}", e);
    } else {
        info!("MongoDB document deleted successfully.");
    }

    // --- MySQL ---
    info!("--- MySQL Workflow ---");
    if let Err(e) = mysql_data::setup_database().await {
        error!("MySQL setup error: {}", e);
    } else {
        info!("MySQL database setup successfully.");
    }
    let mysql_json = create_sample_json();
    if let Err(e) = mysql_data::insert_json_into_table(&mysql_json).await {
        error!("MySQL insert error: {}", e);
    } else {
        info!("MySQL row inserted successfully.");
    }
    match mysql_data::fetch_data_as_csv().await {
        Ok(csv) => info!("MySQL CSV:\n{}", csv),
        Err(e) => error!("MySQL fetch error: {}", e),
    }
    if let Err(e) = mysql_data::delete_row(1756655999, "den", "bmp280").await {
        error!("MySQL delete error: {}", e);
    } else {
        info!("MySQL row deleted successfully.");
    }

    // --- PostgreSQL ---
    info!("--- PostgreSQL Workflow ---");
    if let Err(e) = postgres_data::setup_database().await {
        error!("Postgres setup error: {}", e);
    } else {
        info!("PostgreSQL database setup successfully.");
    }
    let pg_json = create_sample_json();
    if let Err(e) = postgres_data::insert_json_into_table(&pg_json).await {
        error!("Postgres insert error: {}", e);
    } else {
        info!("PostgreSQL row inserted successfully.");
    }
    match postgres_data::fetch_data_as_csv().await {
        Ok(csv) => info!("Postgres CSV:\n{}", csv),
        Err(e) => error!("Postgres fetch error: {}", e),
    }
    if let Err(e) = postgres_data::delete_row(1756655999, "den", "bmp280").await {
        error!("Postgres delete error: {}", e);
    } else {
        info!("PostgreSQL row deleted successfully.");
    }

    // --- Redis ---
    info!("--- Redis Workflow ---");
    let redis_key = "sensor:bmp280:single";
    let redis_json = create_sample_json();
    if let Err(e) = redis_data::set_json(redis_key, &redis_json).await {
        error!("Redis SET error: {}", e);
    } else {
        info!("Redis key-value pair set successfully.");
    }
    match redis_data::get_json(redis_key).await {
        Ok(val) => info!("Redis GET value: {}", val),
        Err(e) => error!("Redis GET error: {}", e),
    }
    if let Err(e) = redis_data::delete_key(redis_key).await {
        error!("Redis DEL error: {}", e);
    } else {
        info!("Redis key deleted successfully.");
    }
    // List example
    let list_key = "sensor:bmp280:list";
    for _ in 0..3 {
        let json = create_sample_json();
        if let Err(e) = redis_data::lpush_json(list_key, &json).await {
            error!("Redis LPUSH error: {}", e);
        } else {
            info!("Redis LPUSH executed successfully.");
        }
    }
    match redis_data::lrange_json(list_key).await {
        Ok(values) => match redis_data::json_array_to_csv(&values) {
            Ok(csv) => info!("Redis List as CSV:\n{}", csv),
            Err(e) => error!("Redis List CSV error: {}", e),
        },
        Err(e) => error!("Redis LRANGE error: {}", e),
    }
    if let Err(e) = redis_data::delete_key(list_key).await {
        error!("Redis DEL error: {}", e);
    } else {
        info!("Redis list key deleted successfully.");
    }
    // Set example
    let set_key = "sensor:bmp280:set";
    for _ in 0..3 {
        let json = create_sample_json();
        if let Err(e) = redis_data::sadd_json(set_key, &json).await {
            error!("Redis SADD error: {}", e);
        } else {
            info!("Redis SADD executed successfully.");
        }
    }
    match redis_data::smembers_json(set_key).await {
        Ok(values) => match redis_data::json_array_to_csv(&values) {
            Ok(csv) => info!("Redis Set as CSV:\n{}", csv),
            Err(e) => error!("Redis Set CSV error: {}", e),
        },
        Err(e) => error!("Redis SMEMBERS error: {}", e),
    }
    if let Err(e) = redis_data::delete_key(set_key).await {
        error!("Redis DEL error: {}", e);
    } else {
        info!("Redis set key deleted successfully.");
    }
}
