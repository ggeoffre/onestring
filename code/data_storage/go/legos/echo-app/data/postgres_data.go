// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// SensorDataPostgresJSON is for unmarshaling JSON where recorded might be a string
type SensorDataPostgresJSON struct {
	Recorded    interface{} `json:"recorded"` // Accept both string and number
	Location    string      `json:"location"`
	Sensor      string      `json:"sensor"`
	Measurement string      `json:"measurement"`
	Units       string      `json:"units"`
	Value       float64     `json:"value"`
}

type SensorDataPostgres struct {
	Recorded    int64   `json:"recorded" postgres:"recorded"`
	Location    string  `json:"location" postgres:"location"`
	Sensor      string  `json:"sensor" postgres:"sensor"`
	Measurement string  `json:"measurement" postgres:"measurement"`
	Units       string  `json:"units" postgres:"units"`
	Value       float64 `json:"value" postgres:"value"`
}

// PostgresDataAccess implements the SensorDataAccess interface for PostgreSQL
type PostgresDataAccess struct {
	DB *sql.DB
}

// NewPostgresDataAccess creates a new PostgresDataAccess instance and connects
func NewPostgresDataAccess() (*PostgresDataAccess, error) {
	// Connect to default 'postgres' database to check/create target DB
	db, err := sql.Open("pgx", "postgres://postgres:postgres@192.168.1.60:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open default database: %w", err)
	}

	// Create database if it doesn't exist
	databaseName := "sensor_data_db"
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", databaseName)
	var exists int

	err = db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		db.Close()
		return nil, fmt.Errorf("failed to check for database existence: %w", err)
	}

	if exists == 0 {
		_, err = db.Exec(`CREATE DATABASE sensor_data_db`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		fmt.Println("Database created")
	}
	db.Close()

	// Connect to the specific 'sensor_data_db'
	targetDB, err := sql.Open("pgx", "postgres://postgres:postgres@192.168.1.60:5432/sensor_data_db?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open target database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Create table if it doesn't exist
	tableQuery := `CREATE TABLE IF NOT EXISTS sensor_data (
		recorded BIGINT NOT NULL,
		location VARCHAR(255) NOT NULL,
		sensor VARCHAR(255) NOT NULL,
		measurement VARCHAR(255) NOT NULL,
		units VARCHAR(255) NOT NULL,
		value NUMERIC(10, 2) NOT NULL
	)`
	if _, err := targetDB.Exec(tableQuery); err != nil {
		targetDB.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	fmt.Println("Table created or already exists")

	return &PostgresDataAccess{DB: targetDB}, nil
}

// Close closes the database connection
func (p *PostgresDataAccess) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

// parseRecorded converts the recorded field to int64 whether it's a string or number
func parseRecordedPostgres(recorded interface{}) (int64, error) {
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

func (p *PostgresDataAccess) LogSensorData(jsonData string) error {
	var jsonStruct SensorDataPostgresJSON
	if err := json.Unmarshal([]byte(jsonData), &jsonStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}

	recorded, err := parseRecordedPostgres(jsonStruct.Recorded)
	if err != nil {
		return fmt.Errorf("failed to parse recorded field: %w", err)
	}

	sensorData := SensorDataPostgres{
		Recorded:    recorded,
		Location:    jsonStruct.Location,
		Sensor:      jsonStruct.Sensor,
		Measurement: jsonStruct.Measurement,
		Units:       jsonStruct.Units,
		Value:       jsonStruct.Value,
	}

	query := `INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = p.DB.Exec(query,
		sensorData.Recorded,
		sensorData.Location,
		sensorData.Sensor,
		sensorData.Measurement,
		sensorData.Units,
		sensorData.Value)
	if err != nil {
		return fmt.Errorf("failed to insert into PostgreSQL: %w", err)
	}

	return nil
}

func (p *PostgresDataAccess) FetchSensorData() ([]string, error) {
	rows, err := p.DB.Query(`SELECT recorded, location, sensor, measurement, units, value FROM sensor_data`)
	if err != nil {
		return nil, fmt.Errorf("failed to query data from PostgreSQL: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var sensorData SensorDataPostgres
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
			fmt.Printf("failed to marshal SensorDataPostgres to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (p *PostgresDataAccess) PurgeSensorData() error {
	_, err := p.DB.Exec("DELETE FROM sensor_data")
	if err != nil {
		return fmt.Errorf("failed to delete data from PostgreSQL: %w", err)
	}

	return nil
}
