// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use std::error::Error;

pub trait SensorDataAccess: Send + Sync {
    fn log_sensor_data<'a>(&'a self, json_data: &'a str) -> tokio::task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>>;
    fn fetch_sensor_data(&self) -> tokio::task::JoinHandle<Result<Vec<String>, Box<dyn Error + Send + Sync>>>;
    fn purge_sensor_data(&self) -> tokio::task::JoinHandle<Result<(), Box<dyn Error + Send + Sync>>>;
}
