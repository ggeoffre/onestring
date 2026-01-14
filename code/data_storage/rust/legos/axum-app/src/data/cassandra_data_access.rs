// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use crate::data::sensor_data_access_trait::SensorDataAccess;
use crate::data::sensor_data_json_helper::validate_sensor_json;
use std::error::Error;
use tokio::task;
use cdrs_tokio::cluster::session::{Session, TcpSessionBuilder, SessionBuilder};
use cdrs_tokio::cluster::NodeTcpConfigBuilder;
use cdrs_tokio::load_balancing::RoundRobinLoadBalancingStrategy;
use cdrs_tokio::transport::TransportTcp;
use cdrs_tokio::cluster::TcpConnectionManager;
use cdrs_tokio::query_values;
use cdrs_tokio::types::IntoRustByName;

const CASSANDRA_SERVER_IP: &str = "localhost";
const CASSANDRA_SERVER_PORT: u16 = 9042;
const KEYSPACE_NAME: &str = "sensor_data_db";
const TABLE_NAME: &str = "sensor_data";

async fn get_session() -> Result<Session<TransportTcp, TcpConnectionManager, RoundRobinLoadBalancingStrategy<TransportTcp, TcpConnectionManager>>, Box<dyn Error + Send + Sync>> {
    let cluster_config = NodeTcpConfigBuilder::new()
        .with_contact_point(format!("{}:{}", CASSANDRA_SERVER_IP, CASSANDRA_SERVER_PORT).into())
        .build()
        .await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Cluster config error: {}", e))) })?;

    let lb = RoundRobinLoadBalancingStrategy::new();
    let session = TcpSessionBuilder::new(lb, cluster_config)
        .build()
        .await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Session build error: {}", e))) })?;

    Ok(session)
}

async fn create_keyspace_and_table() -> Result<(), Box<dyn Error + Send + Sync>> {
    let session = get_session().await?;
    let create_ks = format!(
        "CREATE KEYSPACE IF NOT EXISTS {} WITH REPLICATION = {{ 'class' : 'SimpleStrategy', 'replication_factor' : 1 }};",
        KEYSPACE_NAME
    );
    session.query(create_ks).await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Keyspace create error: {}", e))) })?;

    let use_ks = format!("USE {};", KEYSPACE_NAME);
    session.query(use_ks).await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Use keyspace error: {}", e))) })?;

    let create_table = format!(
        "CREATE TABLE IF NOT EXISTS {} (location TEXT, recorded BIGINT, sensor TEXT, measurement TEXT, units TEXT, value DOUBLE, PRIMARY KEY ((location), recorded, sensor));",
        TABLE_NAME
    );
    session.query(create_table).await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Table create error: {}", e))) })?;

    println!("Cassandra keyspace and table created successfully.");
    Ok(())
}

pub struct CassandraDataAccess;

impl CassandraDataAccess {
    pub fn new() -> Self {
        CassandraDataAccess
    }
}

impl SensorDataAccess for CassandraDataAccess {
    fn log_sensor_data<'a>(&'a self, json_data: &'a str) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        let json_owned = json_data.to_string();
        task::spawn(async move {
            // Ensure keyspace and table exist
            create_keyspace_and_table().await?;

            let parsed = validate_sensor_json(&json_owned)
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, e)) })?;

            let location = parsed["location"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'location' field")) })?
                .to_string();

            // Handle recorded as either integer or string
            let recorded = parsed["recorded"].as_i64()
                .or_else(|| parsed["recorded"].as_str().and_then(|s| s.parse::<i64>().ok()))
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'recorded' field")) })?;

            let sensor = parsed["sensor"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'sensor' field")) })?
                .to_string();

            let measurement = parsed["measurement"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'measurement' field")) })?
                .to_string();

            let units = parsed["units"].as_str()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'units' field")) })?
                .to_string();

            // Handle value as either float or integer or string
            let value = parsed["value"].as_f64()
                .or_else(|| parsed["value"].as_i64().map(|i| i as f64))
                .or_else(|| parsed["value"].as_str().and_then(|s| s.parse::<f64>().ok()))
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "Missing or invalid 'value' field")) })?;

            let session = get_session().await?;
            let insert_query = format!(
                "INSERT INTO {}.{} (location, recorded, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?);",
                KEYSPACE_NAME, TABLE_NAME
            );
            let prepared = session.prepare(&insert_query).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Prepared statement error: {}", e))) })?;
            let values = query_values!(location, recorded, sensor, measurement, units, value);
            session.exec_with_values(&prepared, values).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Query execution error: {}", e))) })?;

            println!("Logging sensor data to Cassandra: {}", json_owned);
            Ok(())
        })
    }

    fn fetch_sensor_data(&self) -> task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Fetching sensor data from Cassandra");

            let session = get_session().await?;
            let select_query = format!(
                "SELECT location, recorded, sensor, measurement, units, value FROM {}.{};",
                KEYSPACE_NAME, TABLE_NAME
            );
            let response = session.query(select_query).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Query execution error: {}", e))) })?;
            let rows = response.response_body()
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Response body error: {}", e))) })?
                .into_rows()
                .ok_or_else(|| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, "No rows found")) })?;

            let mut json_strings: Vec<String> = Vec::new();
            for row in rows {
                let location: String = row.get_by_name("location").unwrap().unwrap();
                let recorded: i64 = row.get_by_name("recorded").unwrap().unwrap();
                let sensor: String = row.get_by_name("sensor").unwrap().unwrap();
                let measurement: String = row.get_by_name("measurement").unwrap().unwrap();
                let units: String = row.get_by_name("units").unwrap().unwrap();
                let value: f64 = row.get_by_name("value").unwrap().unwrap();

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
            println!("Purging sensor data from Cassandra");

            let session = get_session().await?;
            let truncate_query = format!("TRUNCATE {}.{};", KEYSPACE_NAME, TABLE_NAME);
            session.query(truncate_query).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Truncate error: {}", e))) })?;

            println!("Cassandra sensor data purged successfully.");
            Ok(())
        })
    }
}
