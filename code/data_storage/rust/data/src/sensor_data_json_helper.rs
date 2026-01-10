// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use serde_json::{Value, json};
use log::{error};

pub fn create_sample_json() -> String {
    json!({
        "recorded": 1756655999,
        "location": "den",
        "sensor": "bmp280",
        "measurement": "temperature",
        "units": "C",
        "value": 22.3
    }).to_string()
}

pub fn validate_sensor_json(json_str: &str) -> Result<Value, String> {
    let parsed: Value = serde_json::from_str(json_str)
        .map_err(|e| format!("JSON parse error: {}", e))?;
    let required = ["recorded", "location", "sensor", "measurement", "units", "value"];
    for &field in &required {
        if !parsed.get(field).is_some() {
            error!("Missing field: {}", field);
            return Err(format!("Missing field: {}", field));
        }
    }
    Ok(parsed)
}
