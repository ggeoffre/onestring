# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import datetime
import decimal
import json
import random
import time

# Constants for sensor data generation
sample_timestamp = 1768570200
sample_location = "den"
sample_sensor = "bmp280"
sample_measurement = "temperature"
sample_units = "C"
sample_value = 22.3

# Temperature range for random data generation
min_temperature = 22.4
max_temperature = 32.1
temperature_precision = 1

# Sensor data defaults
sensor_data = {
    "recorded": sample_timestamp,
    "location": sample_location,
    "sensor": sample_sensor,
    "measurement": sample_measurement,
    "units": sample_units,
    "value": sample_value,
}

# Sensor data string and dictionary
sensor_data_string = json.dumps(sensor_data)
sensor_data_dict = json.loads(sensor_data_string)


def generate_random_sensor_data():
    data = sensor_data.copy()
    data["recorded"] = int(time.time())
    data["value"] = round(
        random.uniform(min_temperature, max_temperature), temperature_precision
    )
    return data


def create_insert_data_tuple(data):
    """Prepares a tuple for SQL/CQL insertion."""
    if not isinstance(data, dict):
        raise TypeError("Input must be a dictionary")

    # Return directly; no unreachable returns after raise
    return (
        int(data.get("recorded", 0)),
        data.get("location"),
        data.get("sensor"),
        data.get("measurement"),
        data.get("units"),
        float(data.get("value", 0.0)),
    )


def json_default(obj):
    """Handles non-standard JSON types."""
    if isinstance(obj, decimal.Decimal):
        return float(obj)
    if isinstance(obj, (datetime.datetime, datetime.date)):
        return obj.isoformat()
    if isinstance(obj, ObjectId):
        return str(obj)
    raise TypeError(f"Object of type {obj.__class__.__name__} is not serializable")


def json_list_to_csv(json_string_list):
    headers = list(json_string_list[0].keys())
    csv_lines = [",".join(headers)]

    for item in json_string_list:
        row = []
        for key in headers:
            val = str(item.get(key, ""))
            val = val.replace('"', '""')
            if "," in val or '"' in val:
                val = f'"{val}"'
            row.append(val)
        csv_lines.append(",".join(row))

    return "\n".join(csv_lines)
