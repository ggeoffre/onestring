// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"sort"
	"strings"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// JSON String Literal
const SENSOR_DATA_JSON_STRING = `{
	"recorded": 1768237200,
	"location": "den",
	"sensor": "bmp280",
	"measurement": "temperature",
	"units": "C",
	"value": 22.3
}`

// JSON Object Helper Function
func GetSensorDataJsonObject() map[string]interface{} {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(SENSOR_DATA_JSON_STRING), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling SENSOR_DATA_JSON_STRING: %v\n", err)
		return nil
	}
	return data
}

// JSON Array Helper Function
func GetSensorDataJsonArray() []interface{} {
	jsonArrayString := "[" + SENSOR_DATA_JSON_STRING + "]"
	var data []interface{}
	err := json.Unmarshal([]byte(jsonArrayString), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON array string: %v\n", err)
		return nil
	}
	return data
}

// Create a new map to hold the sensor data with random values for keys recorded and value
func GetRandomSensorDataJsonString() []byte {
	// Create a new map to hold the sensor data
	var sensorData = GetSensorDataJsonObject()

	// Generate a random seed to use for generating random values
	// Note: rand.Seed is deprecated in newer Go versions (1.20+) as global random is seeded automatically,
	// but kept here for compatibility with older environments if needed.
	rand.Seed(time.Now().UnixNano())
	sensorData["recorded"] = time.Now().Unix()
	var value float64 = 22.4 + rand.Float64()*(32.1-22.4)
	value = math.Round(value*float64(math.Pow10(1))) / float64(math.Pow10(1))
	sensorData["value"] = value

	// Convert the map to a byte slice using json.Marshal
	b, _ := json.Marshal(sensorData)

	return b
}

// RemoveKeyValuePair removes a key-value pair from a map
func RemoveKeyValuePair(mapVar map[string]string, keyToRemove string) {
	if _, ok := mapVar[keyToRemove]; ok {
		delete(mapVar, keyToRemove)
	}
}

// JsonToCSV converts a JSON array string to CSV format
// Input: []string where each string is a JSON object from FetchSensorData
// Output: CSV string with headers and rows
func JsonToCSV(jsonArray []string) (string, error) {
	if len(jsonArray) == 0 {
		return "", nil
	}

	// Parse all JSON objects to get their fields
	var records []map[string]interface{}
	var allKeys map[string]bool = make(map[string]bool)

	for _, jsonStr := range jsonArray {
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &record); err != nil {
			return "", fmt.Errorf("failed to parse JSON: %w", err)
		}
		records = append(records, record)

		// Collect all unique keys
		for key := range record {
			allKeys[key] = true
		}
	}

	// Sort keys for consistent column order
	var keys []string
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build CSV
	var csvBuilder strings.Builder

	// Write header row
	csvBuilder.WriteString(strings.Join(keys, ","))
	csvBuilder.WriteString("\n")

	// Write data rows
	for _, record := range records {
		var values []string
		for _, key := range keys {
			value := record[key]

			// Special handling for numeric values to avoid scientific notation
			if floatVal, ok := value.(float64); ok {
				// Check if it's actually an integer (no decimal part)
				if floatVal == float64(int64(floatVal)) {
					values = append(values, fmt.Sprintf("%d", int64(floatVal)))
				} else {
					values = append(values, fmt.Sprintf("%v", value))
				}
			} else {
				values = append(values, fmt.Sprintf("%v", value))
			}
		}
		csvBuilder.WriteString(strings.Join(values, ","))
		csvBuilder.WriteString("\n")
	}

	return csvBuilder.String(), nil
}
