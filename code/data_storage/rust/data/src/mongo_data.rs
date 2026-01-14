// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use mongodb::{options::ClientOptions, Client, bson::doc, bson::Document};
use futures::stream::StreamExt;
use log::info;
use crate::sensor_data_json_helper::validate_sensor_json;

async fn get_sensor_data_collection() -> Result<mongodb::Collection<Document>, String> {
    let client_uri = "mongodb://localhost:27017";
    let options = ClientOptions::parse(client_uri).await.map_err(|e| format!("ClientOptions error: {}", e))?;
    let client = Client::with_options(options).map_err(|e| format!("Client error: {}", e))?;
    Ok(client.database("sensor_data_db").collection("sensor_data"))
}

pub async fn insert_json_into_collection(json_data: &str) -> Result<(), String> {
    let collection = get_sensor_data_collection().await?;
    let parsed_json = validate_sensor_json(json_data)?;
    let bson_doc = mongodb::bson::to_document(&parsed_json).map_err(|e| format!("BSON conversion error: {}", e))?;
    collection.insert_one(bson_doc, None).await.map_err(|e| format!("Insert error: {}", e))?;
    info!("JSON inserted successfully!");
    Ok(())
}

pub async fn fetch_collection_as_csv() -> Result<String, String> {
    let collection = get_sensor_data_collection().await?;
    let mut cursor = collection.find(None, None).await.map_err(|e| format!("Find error: {}", e))?;
    let mut csv_data = String::from("recorded,location,sensor,measurement,units,value\n");
    while let Some(result) = cursor.next().await {
        let doc = result.map_err(|e| format!("Cursor error: {}", e))?;
        let recorded = doc.get_i64("recorded").unwrap_or_default();
        let location = doc.get_str("location").unwrap_or_default();
        let sensor = doc.get_str("sensor").unwrap_or_default();
        let measurement = doc.get_str("measurement").unwrap_or_default();
        let units = doc.get_str("units").unwrap_or_default();
        let value = doc.get_f64("value").unwrap_or_default();
        csv_data.push_str(&format!("{},{},{},{},{},{}\n", recorded, location, sensor, measurement, units, value));
    }
    Ok(csv_data)
}

pub async fn delete_row(recorded: i64, location: &str, sensor: &str) -> Result<(), String> {
    let collection = get_sensor_data_collection().await?;
    let filter = doc! { "recorded": recorded, "location": location, "sensor": sensor };
    collection.delete_one(filter, None).await.map_err(|e| format!("Delete error: {}", e))?;
    info!("Row deleted successfully.");
    Ok(())
}
