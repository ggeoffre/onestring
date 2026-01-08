// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// SensorDataMySQLJSON is for unmarshaling JSON where recorded might be a string
type SensorDataMySQLJSON struct {
	Recorded    interface{} `json:"recorded"` // Accept both string and number
	Location    string      `json:"location"`
	Sensor      string      `json:"sensor"`
	Measurement string      `json:"measurement"`
	Units       string      `json:"units"`
	Value       float64     `json:"value"`
}

type SensorDataMySQL struct {
	Recorded    int64   `json:"recorded" mysql:"recorded"`
	Location    string  `json:"location" mysql:"location"`
	Sensor      string  `json:"sensor" mysql:"sensor"`
	Measurement string  `json:"measurement" mysql:"measurement"`
	Units       string  `json:"units" mysql:"units"`
	Value       float64 `json:"value" mysql:"value"`
}

// MySQLDataAccess implements the SensorDataAccess interface for MySQL
type MySQLDataAccess struct {
	DB *sql.DB
}

// NewMySQLDataAccess creates a new MySQLDataAccess instance and connects
func NewMySQLDataAccess() (*MySQLDataAccess, error) {
	// Initialize MySQL connection to default 'mysql' database
	db, err := sql.Open("mysql", "root:@tcp(192.168.1.60:3306)/mysql?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create database if it doesn't exist
	databaseName := "sensor_data_db"
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", databaseName)
	if _, err := db.Exec(query); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Use the created database
	if _, err := db.Exec("USE sensor_data_db"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to select database: %w", err)
	}

	fmt.Println("Connected to MySQL and selected DB")

	// Create table if it doesn't exist
	tableQuery := `CREATE TABLE IF NOT EXISTS sensor_data (
		recorded BIGINT,
		location VARCHAR(255),
		sensor VARCHAR(255),
		measurement VARCHAR(255),
		units VARCHAR(255),
		value DECIMAL(10, 2)
	)`
	if _, err := db.Exec(tableQuery); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Table created or already exists")

	return &MySQLDataAccess{DB: db}, nil
}

// Close closes the database connection
func (m *MySQLDataAccess) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

// parseRecorded converts the recorded field to int64 whether it's a string or number
func parseRecordedMySQL(recorded interface{}) (int64, error) {
	switch v := recorded.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("recorded field has unexpected type: %T", v)
	}
}

func (m *MySQLDataAccess) LogSensorData(jsonData string) error {
	var jsonStruct SensorDataMySQLJSON
	if err := json.Unmarshal([]byte(jsonData), &jsonStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}

	recorded, err := parseRecordedMySQL(jsonStruct.Recorded)
	if err != nil {
		return fmt.Errorf("failed to parse recorded field: %w", err)
	}

	sensorData := SensorDataMySQL{
		Recorded:    recorded,
		Location:    jsonStruct.Location,
		Sensor:      jsonStruct.Sensor,
		Measurement: jsonStruct.Measurement,
		Units:       jsonStruct.Units,
		Value:       jsonStruct.Value,
	}

	query := `INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err = m.DB.Exec(query,
		sensorData.Recorded,
		sensorData.Location,
		sensorData.Sensor,
		sensorData.Measurement,
		sensorData.Units,
		sensorData.Value)
	if err != nil {
		return fmt.Errorf("failed to insert into MySQL: %w", err)
	}

	return nil
}

func (m *MySQLDataAccess) FetchSensorData() ([]string, error) {
	rows, err := m.DB.Query(`SELECT recorded, location, sensor, measurement, units, value FROM sensor_data`)
	if err != nil {
		return nil, fmt.Errorf("failed to query data from MySQL: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var sensorData SensorDataMySQL
		err := rows.Scan(
			&sensorData.Recorded,
			&sensorData.Location,
			&sensorData.Sensor,
			&sensorData.Measurement,
			&sensorData.Units,
			&sensorData.Value,
		)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)
			continue
		}

		jsonBytes, err := json.Marshal(sensorData)
		if err != nil {
			fmt.Printf("failed to marshal SensorDataMySQL to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (m *MySQLDataAccess) PurgeSensorData() error {
	_, err := m.DB.Exec("DELETE FROM sensor_data")
	if err != nil {
		return fmt.Errorf("failed to delete data from MySQL: %w", err)
	}

	return nil
}
