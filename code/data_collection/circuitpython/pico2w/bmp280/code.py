# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os
import ssl
import time

import adafruit_bmp280
import adafruit_ntp
import adafruit_requests
import board
import busio
import rtc
import socketpool
import wifi

# 1. Initialize I2C for Pico 2 W (Standard GP4/GP5)
i2c = busio.I2C(board.GP5, board.GP4)
bmp280 = adafruit_bmp280.Adafruit_BMP280_I2C(i2c)

# 2. Connect to WiFi using settings.toml variables
print("Connecting to WiFi...")
wifi.radio.connect(
    os.getenv("CIRCUITPY_WIFI_SSID"), os.getenv("CIRCUITPY_WIFI_PASSWORD")
)
print("Connected to", os.getenv("CIRCUITPY_WIFI_SSID"))

# 3. Setup Network and Time
pool = socketpool.SocketPool(wifi.radio)
requests = adafruit_requests.Session(pool, ssl.create_default_context())
ntp = adafruit_ntp.NTP(pool, tz_offset=0)

# Sync system RTC with NTP time once at startup
rtc.RTC().datetime = ntp.datetime
api_url = os.getenv("API_URL")

while True:
    temperature = bmp280.temperature

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
        response = requests.post(api_url, json=payload)
        print(f"Status: {response.status_code}")
        response.close()
    except Exception as e:
        print(f"Request failed: {e}")

    time.sleep(61)
