// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.database;

import java.util.*;
import org.bson.Document;
import org.json.*;

public class SensorDataJsonHelper {

    // Default sensor data for convenience
    public static final String DEFAULT_SENSOR_DATA_JSON_STRING =
        "{\"recorded\": 1768237200, \"location\": \"den\", \"sensor\": \"bmp280\", \"measurement\": \"temperature\", \"units\": \"C\", \"value\": 22.3}";

    /**
     * Generates a BSON Document for sensor data.
     * If parameters are not provided, uses defaults and randomizes 'value' and 'recorded'.
     */
    public static Document getSensorDataDocument() {
        return getSensorDataDocument(
            "den",
            "bmp280",
            "temperature",
            "C",
            22.4 + new Random().nextDouble() * (32.1 - 22.4), // random value
            System.currentTimeMillis() / 1000 // current time in seconds
        );
    }

    /**
     * Generates a BSON Document for sensor data with parameters.
     */
    public static Document getSensorDataDocument(
        String location,
        String sensor,
        String measurement,
        String units,
        double value,
        long recorded
    ) {
        Document sensorDoc = new Document();
        sensorDoc.put("recorded", recorded);
        sensorDoc.put("location", location);
        sensorDoc.put("sensor", sensor);
        sensorDoc.put("measurement", measurement);
        sensorDoc.put("units", units);
        // Round to one decimal place
        sensorDoc.put("value", Math.round(value * 10.0) / 10.0);
        return sensorDoc;
    }

    /**
     * Returns a JSONObject representing sensor data.
     * Uses default/randomized values.
     */
    public static JSONObject getSensorDataJSONObject() {
        Document document = getSensorDataDocument();
        return new JSONObject(document.toJson());
    }

    /**
     * Returns a JSONObject representing sensor data with parameters.
     */
    public static JSONObject getSensorDataJSONObject(
        String location,
        String sensor,
        String measurement,
        String units,
        double value,
        long recorded
    ) {
        Document document = getSensorDataDocument(
            location,
            sensor,
            measurement,
            units,
            value,
            recorded
        );
        return new JSONObject(document.toJson());
    }

    /**
     * Returns a JSON string representing sensor data.
     * Uses default/randomized values.
     */
    public static String getSensorDataJSONString() {
        Document document = getSensorDataDocument();
        return document.toJson();
    }

    /**
     * Returns a JSON string representing sensor data with parameters.
     */
    public static String getSensorDataJSONString(
        String location,
        String sensor,
        String measurement,
        String units,
        double value,
        long recorded
    ) {
        Document document = getSensorDataDocument(
            location,
            sensor,
            measurement,
            units,
            value,
            recorded
        );
        return document.toJson();
    }

    // Custom exception for JSON utility operations
    public static class JsonUtilityException extends RuntimeException {

        public JsonUtilityException(String message, Throwable cause) {
            super(message, cause);
        }
    }

    /**
     * Converts an array of JSON strings to a CSV format.
     * Collects all unique headers and aligns rows accordingly.
     */
    public static String jsonArrayToCsv(String[] jsonStrings) {
        if (jsonStrings == null || jsonStrings.length == 0) {
            return "";
        }

        try {
            // Collect all unique headers
            Set<String> headerSet = new LinkedHashSet<>();
            List<JSONObject> jsonObjects = new ArrayList<>();

            for (String jsonString : jsonStrings) {
                if (jsonString == null || jsonString.isEmpty()) {
                    continue;
                }
                JSONObject jsonObject = new JSONObject(jsonString);
                jsonObjects.add(jsonObject);
                headerSet.addAll(jsonObject.keySet());
            }

            List<String> headers = new ArrayList<>(headerSet);

            StringBuilder csvBuilder = new StringBuilder();
            csvBuilder.append(String.join(",", headers)).append("\n");

            // Build rows aligned with headers
            for (JSONObject jsonObject : jsonObjects) {
                List<String> row = new ArrayList<>();
                for (String header : headers) {
                    row.add(jsonObject.optString(header, ""));
                }
                csvBuilder.append(String.join(",", row)).append("\n");
            }

            return csvBuilder.toString();
        } catch (Exception e) {
            throw new JsonUtilityException(
                "Error converting JSON array to CSV",
                e
            );
        }
    }
}
