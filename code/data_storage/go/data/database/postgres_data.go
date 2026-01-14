// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SensorDataPostgres struct {
	Recorded    int64   `postgres:"recorded"`
	Location    string  `postgres:"location"`
	Sensor      string  `postgres:"sensor"`
	Measurement string  `postgres:"measurement"`
	Units       string  `postgres:"units"`
	Value       float64 `postgres:"value"`
}

// PostgresData struct for Postgres
type PostgresData struct {
	DB *sql.DB
}

// NewPostgresData creates a new PostgresData instance and connects
func NewPostgresData() (*PostgresData, error) {
	// 1. Connect to default 'postgres' database to check/create target DB
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open default database: %w", err)
	}

	databaseName := "sensor_data_db"
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", databaseName)
	var exists int

	err = db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		db.Close()
		return nil, fmt.Errorf("failed to check for database existence: %w", err)
	}

	if exists == 0 {
		_, err = db.Exec(`CREATE DATABASE sensor_data_db;`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create database 'sensor_data_db': %w", err)
		}
		fmt.Println("Database created successfully!")
	}
	db.Close()

	// 2. Connect to the specific 'sensor_data_db'
	targetDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/sensor_data_db?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open target database: %w", err)
	}

	fmt.Println("Connected to Postgres")
	return &PostgresData{DB: targetDB}, nil
}

// Close closes the database connection
func (c *PostgresData) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func (c *PostgresData) RawPostgresData() {
	// Create the sensor_data table if it doesn't exist.
	var exists int
	err := c.DB.QueryRow(`SELECT 1 FROM pg_tables WHERE tablename = 'sensor_data'`).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Failed to check for table existence: %v\n", err)
		return
	}

	if exists == 0 {
		_, err = c.DB.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_data (
			recorded BIGINT NOT NULL,
			location VARCHAR(255) NOT NULL,
			sensor VARCHAR(255) NOT NULL,
			measurement VARCHAR(255) NOT NULL,
			units VARCHAR(255) NOT NULL,
			value NUMERIC(10, 2) NOT NULL
		);`)
		if err != nil {
			fmt.Printf("failed to create sensor_data table: %v\n", err)
			return
		}
		fmt.Println("Table created successfully!")
	}

	// Unmarshal JSON string into SensorDataPostgres struct
	var sensorDataPostgres SensorDataPostgres
	if err := json.Unmarshal([]byte(getRandomSensorDataJsonString()), &sensorDataPostgres); err != nil {
		fmt.Printf("Failed to parse json: %v\n", err)
		return
	}
	fmt.Println("JSON string converted to SensorDataPostgres struct")

	// Prepare the SQL statement for inserting sensor data
	var query = `INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
		VALUES ($1, $2, $3, $4, $5, $6)`

	// Execute the insert statement
	_, err = c.DB.Exec(query,
		sensorDataPostgres.Recorded,
		sensorDataPostgres.Location,
		sensorDataPostgres.Sensor,
		sensorDataPostgres.Measurement,
		sensorDataPostgres.Units,
		sensorDataPostgres.Value)
	if err != nil {
		fmt.Printf("Failed to insert data: %v\n", err)
		return
	}
	fmt.Println("SensorDataPostgres stored")

	// Query all rows from the sensor_data table
	var results []string
	rows, err := c.DB.Query(`SELECT recorded, location, sensor, measurement, units, value FROM sensor_data`)
	if err != nil {
		fmt.Printf("failed to query data: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sensorDataPostgres SensorDataPostgres
		err := rows.Scan(
			&sensorDataPostgres.Recorded,
			&sensorDataPostgres.Location,
			&sensorDataPostgres.Sensor,
			&sensorDataPostgres.Measurement,
			&sensorDataPostgres.Units,
			&sensorDataPostgres.Value,
		)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)
			continue
		}

		jsonBytes, err := json.Marshal(sensorDataPostgres)
		if err != nil {
			fmt.Printf("failed to marshal row to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if len(results) > 0 {
		data := []string{"[" + strings.Join(results, ",") + "]"}
		csvString, err := jsonToCSV(strings.Join(data, "\n"))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s", "SensorDataPostgres retrieved", csvString)
		}
	}

	// Delete all records
	_, err = c.DB.Exec("DELETE FROM sensor_data")
	if err != nil {
		fmt.Printf("Failed to delete data: %v\n", err)
		return
	}
	fmt.Println("SensorDataPostgres purged")
}
