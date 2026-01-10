// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

use serde_json::{Value};
use log::{error};

pub fn json_to_csv(json_str: &str) -> Result<String, Box<dyn std::error::Error>> {
    let json_value: serde_json::Value = serde_json::from_str(json_str)?;

    if let serde_json::Value::Object(map) = json_value {
        let mut csv_string = String::new();

        // Write the header row using the keys
        let headers: Vec<&str> = map.keys().map(|k| k.as_str()).collect();
        csv_string.push_str(&headers.join(","));
        csv_string.push('\n');

        // Write the values row
        let values: Vec<String> = map.values().map(|v| v.to_string()).collect();
        csv_string.push_str(&values.join(","));
        csv_string.push('\n');

        Ok(csv_string)
    } else {
        Err("Input JSON must be an object".into())
    }
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
