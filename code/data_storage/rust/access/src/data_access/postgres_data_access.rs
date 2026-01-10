// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use crate::data_access::sensor_data_access_trait::SensorDataAccess;
use std::error::Error;
use tokio::task;

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
            println!("Logging sensor data to Postgres: {}", json_owned);
            Ok(())
        })
    }

    fn fetch_sensor_data(&self) -> task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Fetching sensor data from Postgres");
            Ok(vec![
                r#"{"sensor":"sound","value":75}"#.to_string(),
            ])
        })
    }

    fn purge_sensor_data(&self) -> task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>> {
        task::spawn(async move {
            println!("Purging sensor data from Postgres");
            Ok(())
        })
    }
}
