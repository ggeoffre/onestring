// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use cdrs_tokio::cluster::session::{Session, TcpSessionBuilder, SessionBuilder};
use cdrs_tokio::cluster::NodeTcpConfigBuilder;
use cdrs_tokio::load_balancing::RoundRobinLoadBalancingStrategy;
use cdrs_tokio::transport::TransportTcp;
use cdrs_tokio::cluster::TcpConnectionManager;
use cdrs_tokio::query_values;
use cdrs_tokio::types::IntoRustByName;
use log::info;
use crate::sensor_data_json_helper::validate_sensor_json;

const CASSANDRA_SERVER_IP: &str = "localhost";
const CASSANDRA_SERVER_PORT: u16 = 9042;
const KEYSPACE_NAME: &str = "sensor_data_db";
const TABLE_NAME: &str = "sensor_data";

pub async fn get_session() -> Result<Session<TransportTcp, TcpConnectionManager, RoundRobinLoadBalancingStrategy<TransportTcp, TcpConnectionManager>>, String> {
    let cluster_config = NodeTcpConfigBuilder::new()
        .with_contact_point(format!("{}:{}", CASSANDRA_SERVER_IP, CASSANDRA_SERVER_PORT).into())
        .build()
        .await
        .map_err(|e| format!("Cluster config error: {}", e))?;

    let lb = RoundRobinLoadBalancingStrategy::new();
    let session = TcpSessionBuilder::new(lb, cluster_config)
        .build()
        .await
        .map_err(|e| format!("Session build error: {}", e))?;

    Ok(session)
}

pub async fn create_keyspace_and_table() -> Result<(), String> {
    let session = get_session().await?;
    let create_ks = format!(
        "CREATE KEYSPACE IF NOT EXISTS {} WITH REPLICATION = {{ 'class' : 'SimpleStrategy', 'replication_factor' : 1 }};",
        KEYSPACE_NAME
    );
    session.query(create_ks).await.map_err(|e| format!("Keyspace create error: {}", e))?;

    let use_ks = format!("USE {};", KEYSPACE_NAME);
    session.query(use_ks).await.map_err(|e| format!("Use keyspace error: {}", e))?;

    let create_table = format!(
        "CREATE TABLE IF NOT EXISTS {} (location TEXT, recorded BIGINT, sensor TEXT, measurement TEXT, units TEXT, value DOUBLE, PRIMARY KEY ((location), recorded, sensor));",
        TABLE_NAME
    );
    session.query(create_table).await.map_err(|e| format!("Table create error: {}", e))?;
    info!("Keyspace and table created successfully.");
    Ok(())
}

pub async fn insert_json_row(json_data: &str) -> Result<(), String> {
    let parsed = validate_sensor_json(json_data)?;
    let location = parsed["location"].as_str().unwrap();
    let recorded = parsed["recorded"].as_i64().unwrap();
    let sensor = parsed["sensor"].as_str().unwrap();
    let measurement = parsed["measurement"].as_str().unwrap();
    let units = parsed["units"].as_str().unwrap();
    let value = parsed["value"].as_f64().unwrap();

    let session = get_session().await?;
    let insert_query = format!(
        "INSERT INTO {}.{} (location, recorded, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?);",
        KEYSPACE_NAME, TABLE_NAME
    );
    let prepared = session.prepare(&insert_query).await.map_err(|e| format!("Prepared statement error: {}", e))?;
    let values = query_values!(location, recorded, sensor, measurement, units, value);
    session.exec_with_values(&prepared, values).await.map_err(|e| format!("Query execution error: {}", e))?;
    info!("Data inserted successfully.");
    Ok(())
}

pub async fn fetch_table_as_csv() -> Result<String, String> {
    let session = get_session().await?;
    let select_query = format!(
        "SELECT location, recorded, sensor, measurement, units, value FROM {}.{};",
        KEYSPACE_NAME, TABLE_NAME
    );
    let response = session.query(select_query).await.map_err(|e| format!("Query execution error: {}", e))?;
    let rows = response.response_body().map_err(|e| format!("Response body error: {}", e))?.into_rows().ok_or("No rows found")?;
    let mut csv_data = String::from("location,recorded,sensor,measurement,units,value\n");
    for row in rows {
        let location: String = row.get_by_name("location").unwrap().unwrap();
        let recorded: i64 = row.get_by_name("recorded").unwrap().unwrap();
        let sensor: String = row.get_by_name("sensor").unwrap().unwrap();
        let measurement: String = row.get_by_name("measurement").unwrap().unwrap();
        let units: String = row.get_by_name("units").unwrap().unwrap();
        let value: f64 = row.get_by_name("value").unwrap().unwrap();
        csv_data.push_str(&format!("{},{},{},{},{},{}\n", location, recorded, sensor, measurement, units, value));
    }
    Ok(csv_data)
}

pub async fn delete_row(location: &str, recorded: i64, sensor: &str) -> Result<(), String> {
    let session = get_session().await?;
    let delete_query = format!(
        "DELETE FROM {}.{} WHERE location = ? AND recorded = ? AND sensor = ?;",
        KEYSPACE_NAME, TABLE_NAME
    );
    let prepared = session.prepare(&delete_query).await.map_err(|e| format!("Prepared statement error: {}", e))?;
    let values = query_values!(location, recorded, sensor);
    session.exec_with_values(&prepared, values).await.map_err(|e| format!("Delete execution error: {}", e))?;
    info!("Row deleted successfully.");
    Ok(())
}
