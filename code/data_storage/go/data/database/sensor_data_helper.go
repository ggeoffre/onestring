// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
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
func getRandomSensorDataJsonString() []byte {
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

// removeKeyValuePair removes a key-value pair from a map
func removeKeyValuePair(mapVar map[string]string, keyToRemove string) {
	if _, ok := mapVar[keyToRemove]; ok {
		delete(mapVar, keyToRemove)
	}
}

// jsonToCSV converts a JSON string to a CSV string
func jsonToCSV(jsonData string) (string, error) {

	// Parse the JSON String Literal into a JSON Object (or Map)
	var records []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &records); err != nil {
		return "", fmt.Errorf("failed to unmarshal json for csv: %w", err)
	}
	if len(records) == 0 {
		return "", fmt.Errorf("no records found for csv")
	}

	// Collect all unique keys from the map for the header
	headerMap := make(map[string]struct{})
	for _, rec := range records {
		for k := range rec {
			headerMap[k] = struct{}{}
		}
	}

	// Sort headers for consistent order
	headers := make([]string, 0, len(headerMap))
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	// Create a CSV writer using a bytes.Buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write csv headers: %w", err)
	}

	// Write each record
	for _, rec := range records {
		row := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := rec[h]; ok && v != nil {
				// Handle each value based on its type
				switch val := v.(type) {
				case float64:
					row[i] = strconv.FormatFloat(val, 'f', -1, 64)
				case int, int64, int32:
					row[i] = fmt.Sprintf("%d", val)
				default:
					row[i] = fmt.Sprintf("%v", val)
				}
			} else {
				row[i] = ""
			}
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write csv row: %w", err)
		}
	}

	// Flush the writer to ensure all data is written to the buffer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("writer flush error: %w", err)
	}

	// Return the CSV string from the buffer
	return buf.String(), nil
}
