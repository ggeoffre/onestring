// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type SensorDataMySQL struct {
	Recorded    int64   `mysql:"recorded"`
	Location    string  `mysql:"location"`
	Sensor      string  `mysql:"sensor"`
	Measurement string  `mysql:"measurement"`
	Units       string  `mysql:"units"`
	Value       float64 `mysql:"value"`
}

// MySQLData struct for MySQL
type MySQLData struct {
	DB *sql.DB
}

// NewMySQLData creates a new MySQLData instance and connects
func NewMySQLData() (*MySQLData, error) {
	// Initialize MySQL connection to default 'mysql' database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mysql?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create database if it doesn't exist
	databaseName := "sensor_data_db"
	query := fmt.Sprintf("SHOW DATABASES LIKE '%s'", databaseName)
	rows, err := db.Query(query)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("show database query failed: %w", err)
	}

	if !rows.Next() {
		rows.Close() // Close rows before executing new query
		// Note: removed single quotes around database name for better SQL compliance
		query = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", databaseName)
		if _, err := db.Exec(query); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create database: %v", err)
		}
		fmt.Println("Database created successfully")
	} else {
		rows.Close()
	}

	// Use the created database
	if _, err := db.Exec("USE sensor_data_db"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to select database: %v", err)
	}

	fmt.Println("Connected to MySQL and selected DB")
	return &MySQLData{DB: db}, nil
}

// Close closes the database connection
func (c *MySQLData) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func (c *MySQLData) RawMySQLData() {
	// Create table if it doesn't exist
	tableName := "sensor_data"
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	rows, err := c.DB.Query(query)
	if err != nil {
		fmt.Printf("show table query failed: %v\n", err)
		return
	}

	tableExists := rows.Next()
	rows.Close()

	if !tableExists {
		var tableQuery = `CREATE TABLE sensor_data (
			recorded BIGINT,
			location VARCHAR(255),
			sensor VARCHAR(255),
			measurement VARCHAR(255),
			units VARCHAR(255),
			value DECIMAL(10, 2)
		);`
		if _, err := c.DB.Exec(tableQuery); err != nil {
			fmt.Printf("failed to create table: %v\n", err)
			return
		}
		fmt.Println("Table created successfully")
	}

	// Unmarshal JSON data into SensorDataMySQL struct
	var sensorDataMySQL SensorDataMySQL
	if err := json.Unmarshal([]byte(getRandomSensorDataJsonString()), &sensorDataMySQL); err != nil {
		fmt.Printf("failed to parse JSON: %v\n", err)
		return
	}
	fmt.Println("JSON string converted to SensorDataMySQL struct")

	// Prepare the SQL statement for inserting sensor data
	query = `INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
		VALUES (?, ?, ?, ?, ?, ?)`

	// Execute the insert statement
	_, err = c.DB.Exec(query,
		sensorDataMySQL.Recorded,
		sensorDataMySQL.Location,
		sensorDataMySQL.Sensor,
		sensorDataMySQL.Measurement,
		sensorDataMySQL.Units,
		sensorDataMySQL.Value)
	if err != nil {
		fmt.Printf("failed to insert data: %v\n", err)
		return
	}
	fmt.Println("SensorDataMySQL stored")

	// Query all rows from the sensor_data table
	var results []string
	rows, err = c.DB.Query(`SELECT recorded, location, sensor, measurement, units, value FROM sensor_data`)
	if err != nil {
		fmt.Printf("failed to query data: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sensorDataMySQL SensorDataMySQL
		err := rows.Scan(
			&sensorDataMySQL.Recorded,
			&sensorDataMySQL.Location,
			&sensorDataMySQL.Sensor,
			&sensorDataMySQL.Measurement,
			&sensorDataMySQL.Units,
			&sensorDataMySQL.Value,
		)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)
			continue
		}

		jsonBytes, err := json.Marshal(sensorDataMySQL)
		if err != nil {
			fmt.Printf("failed to marshal SensorDataMySQL to JSON String: %v\n", err)
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
			fmt.Printf("%s\n%s", "SensorDataMySQL retrieved", csvString)
		}
	}

	// Prepare the SQL statement to delete all records from the sensor_data table
	_, err = c.DB.Exec("DELETE FROM sensor_data")
	if err != nil {
		fmt.Printf("Failed to delete data: %v\n", err)
		return
	}
	fmt.Println("SensorDataMySQL purged")
}
