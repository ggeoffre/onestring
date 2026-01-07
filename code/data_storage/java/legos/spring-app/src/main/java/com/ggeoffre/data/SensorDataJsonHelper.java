// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.data;

import java.util.*;
import org.bson.Document;
import org.json.*;

public class SensorDataJsonHelper {

    public static final String SENSOR_DATA_JSON_STRING =
        "{\"recorded\": 1768237200, \"location\": \"den\", \"sensor\": \"bmp280\", \"measurement\": \"temperature\", \"units\": \"C\", \"value\": 22.3}";

    public static Document getSendorDataDocument(String sensorJsonString) {
        return Document.parse(sensorJsonString);
    }

    public static JSONObject getSensorDataJSONObject(String sensorJsonString) {
        return new JSONObject(sensorJsonString);
    }

    public static String getSensorDataJSONString(Document sensorDocument) {
        return sensorDocument.toJson();
    }

    // Custom exception for JSON utility operations
    public static class JsonUtilityException extends RuntimeException {

        public JsonUtilityException(String message, Throwable cause) {
            super(message, cause);
        }
    }

    // Converts an array of JSON strings to a CSV format.
    public static String jsonArrayToCsv(String[] jsonStrings) {
        if (jsonStrings == null || jsonStrings.length == 0) {
            return "";
        }

        try {
            List<String> headers = new ArrayList<>();
            List<List<String>> rows = new ArrayList<>();

            for (String jsonString : jsonStrings) {
                if (jsonString == null || jsonString.isEmpty()) {
                    continue;
                }
                JSONObject jsonObject = new JSONObject(jsonString);
                List<String> row = new ArrayList<>();

                for (String key : jsonObject.keySet()) {
                    if (!headers.contains(key)) {
                        headers.add(key);
                    }
                    row.add(jsonObject.optString(key, ""));
                }
                rows.add(row);
            }

            StringBuilder csvBuilder = new StringBuilder();
            csvBuilder.append(String.join(",", headers)).append("\n");

            for (List<String> row : rows) {
                List<String> orderedRow = new ArrayList<>();
                for (String header : headers) {
                    int index = headers.indexOf(header);
                    orderedRow.add(index < row.size() ? row.get(index) : "");
                }
                csvBuilder.append(String.join(",", orderedRow)).append("\n");
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
