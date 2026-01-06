/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 * Copyright (C) 2025-2026 ggeoffre, LLC
 */

#include <WiFi.h>
#include <Wire.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <Adafruit_BMP280.h>
#include <NTPClient.h>
#include <WiFiUdp.h>

// --- Configuration ---
const char* ssid     = "YOUR_WIFI_SSID";
const char* password = "YOUR_WIFI_PASSWORD";
const char* api_url  = "http://your-server-ip:port/api/endpoint";

// --- Global Objects ---
Adafruit_BMP280 bmp;
WiFiUDP ntpUDP;
NTPClient timeClient(ntpUDP, "pool.ntp.org", 0); // UTC
WiFiClient client; // Explicit client for HTTP operations

void setup() {
  Serial.begin(115200);
  // Wait up to 3 seconds for Serial Monitor to open
  unsigned long startWait = millis();
  while (!Serial && millis() - startWait < 3000);

  // 1. Initialize I2C for ESP32-S3 (GPIO 8=SDA, GPIO 9=SCL)
  Wire.setPins(8, 9);
  if (!bmp.begin(0x77)) {
      Serial.println("Error: Could not find BMP280 sensor!");
      while (1);
  }

  // 2. Connect to WiFi with timeout
  Serial.printf("Connecting to WiFi: %s\n", ssid);
  WiFi.begin(ssid, password);

  int retry_count = 0;
  while (WiFi.status() != WL_CONNECTED && retry_count < 20) {
    delay(500);
    Serial.print(".");
    retry_count++;
  }

  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("\nConnected to WiFi");
    Serial.print("IP Address: ");
    Serial.println(WiFi.localIP());
  } else {
    Serial.println("\nWiFi connection failed. Check credentials.");
  }

  // 3. Setup Time Sync
  timeClient.begin();
  if (WiFi.status() == WL_CONNECTED) {
    timeClient.update();
  }
}

void loop() {
  // Ensure we have a fresh timestamp
  if (WiFi.status() == WL_CONNECTED) {
    timeClient.update();
  }

  float temperature = bmp.readTemperature();

  // 4. Construct Payload (Onestring JSON format)
  JsonDocument doc;
  doc["recorded"]    = timeClient.getEpochTime();
  doc["location"]    = "den";
  doc["sensor"]      = "bmp280";
  doc["measurement"] = "temperature";
  doc["units"]       = "C";
  // Force 1 decimal place string to match CircuitPython style
  doc["value"]       = serialized(String(temperature, 1));

  String jsonString;
  serializeJson(doc, jsonString);
  Serial.print("Payload: ");
  Serial.println(jsonString);

  // 5. POST data to API using recommended WiFiClient method
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;

    // Pass the client reference and URL
    if (http.begin(client, api_url)) {
      http.addHeader("Content-Type", "application/json");

      int httpResponseCode = http.POST(jsonString);

      if (httpResponseCode > 0) {
        Serial.printf("HTTP Status: %d\n", httpResponseCode);
        // Optional: Read response from server
        // String response = http.getString();
        // Serial.println(response);
      } else {
        Serial.printf("POST failed, error: %s\n", http.errorToString(httpResponseCode).c_str());
      }

      http.end(); // Release resources
    }
  } else {
    Serial.println("WiFi disconnected. Skipping POST.");
    // Optional: attempt to reconnect
    // WiFi.begin(ssid, password);
  }

  // 6. Wait exactly 61 seconds
  delay(61000);
}
