# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os
import time

import bmp280
import machine
import network
import ntptime
import urequests
from machine import I2C, Pin

WIFI_SSID = "your_wifi_ssid"
WIFI_PASSWORD = "your_wifi_password"
API_URL = "your-api-endpoint.com"

# 1. Initialize I2C for Pico 2 W (Standard GP4/GP5)
i2c = I2C(0, scl=Pin(5), sda=Pin(4), freq=400000)
bmp_sensor = bmp280.BMP280(i2c, addr=0x77)

# 2. Connect to WiFi using variables
print("Connecting to WiFi...")
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
wlan.connect(WIFI_SSID, WIFI_PASSWORD)
while not wlan.isconnected():
    time.sleep(1)
print("Connected to", WIFI_SSID)

# 3. Setup Network and Time
try:
    ntptime.settime()
except Exception as e:
    print(f"NTP sync failed: {e}")

while True:
    temperature = bmp_sensor.temperature

    # Payload structure as requested
    payload = {
        "recorded": int(time.time()),  # Unix timestamp integer
        "location": "den",
        "sensor": "bmp280",
        "measurement": "temperature",
        "units": "C",
        "value": round(temperature, 1),
    }

    print(f"Payload: {payload}")

    try:
        response = urequests.post(API_URL, json=payload)
        print(f"Status: {response.status_code}")
        response.close()
    except Exception as e:
        print(f"Request failed: {e}")

    time.sleep(61)
