# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

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
