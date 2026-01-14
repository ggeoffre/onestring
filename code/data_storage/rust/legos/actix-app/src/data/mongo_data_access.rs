// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use crate::data::sensor_data_access_trait::SensorDataAccess;
use crate::data::sensor_data_json_helper::validate_sensor_json;
use std::error::Error;
use tokio::task;
use mongodb::{options::ClientOptions, Client, bson::doc, bson::Document};
use futures::stream::StreamExt;

const MONGO_URI: &str = "mongodb://localhost:27017";
const DATABASE_NAME: &str = "sensor_data_db";
const COLLECTION_NAME: &str = "sensor_data";

async fn get_sensor_data_collection() -> Result<mongodb::Collection<Document>, Box<dyn Error + Send + Sync>> {
    let options = ClientOptions::parse(MONGO_URI).await
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("ClientOptions error: {}", e))) })?;
    let client = Client::with_options(options)
        .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Client error: {}", e))) })?;
    Ok(client.database(DATABASE_NAME).collection(COLLECTION_NAME))
}

pub struct MongoDataAccess;

impl MongoDataAccess {
    pub fn new() -> Self {
        MongoDataAccess
    }
}

impl SensorDataAccess for MongoDataAccess {
    fn log_sensor_data<'a>(&'a self, json_data: &'a str) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        let json_owned = json_data.to_string();
        task::spawn(async move {
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

            let collection = get_sensor_data_collection().await?;
            let bson_doc = doc! {
                "recorded": recorded,
                "location": location,
                "sensor": sensor,
                "measurement": measurement,
                "units": units,
                "value": value
            };
            collection.insert_one(bson_doc, None).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Insert error: {}", e))) })?;

            println!("Logging sensor data to Mongo: {}", json_owned);
            Ok(())
        })
    }

    fn fetch_sensor_data(&self) -> task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Fetching sensor data from Mongo");

            let collection = get_sensor_data_collection().await?;
            let mut cursor = collection.find(None, None).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Find error: {}", e))) })?;

            let mut json_strings: Vec<String> = Vec::new();
            while let Some(result) = cursor.next().await {
                let doc = result
                    .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Cursor error: {}", e))) })?;

                let recorded = doc.get_i64("recorded").unwrap_or_default();
                let location = doc.get_str("location").unwrap_or_default();
                let sensor = doc.get_str("sensor").unwrap_or_default();
                let measurement = doc.get_str("measurement").unwrap_or_default();
                let units = doc.get_str("units").unwrap_or_default();
                let value = doc.get_f64("value").unwrap_or_default();

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
            println!("Purging sensor data from Mongo");

            let collection = get_sensor_data_collection().await?;
            collection.delete_many(doc! {}, None).await
                .map_err(|e| -> Box<dyn Error + Send + Sync> { Box::new(std::io::Error::new(std::io::ErrorKind::Other, format!("Delete error: {}", e))) })?;

            println!("Mongo sensor data purged successfully.");
            Ok(())
        })
    }
}
