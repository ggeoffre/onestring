// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

// SensorDataJSON is for unmarshaling JSON where recorded might be a string
type SensorDataJSON struct {
	Recorded    interface{} `json:"recorded"` // Accept both string and number
	Location    string      `json:"location"`
	Sensor      string      `json:"sensor"`
	Measurement string      `json:"measurement"`
	Units       string      `json:"units"`
	Value       float64     `json:"value"`
}

type SensorDataCassandra struct {
	Recorded    int64   `json:"recorded" cassandra:"recorded"`
	Location    string  `json:"location" cassandra:"location"`
	Sensor      string  `json:"sensor" cassandra:"sensor"`
	Measurement string  `json:"measurement" cassandra:"measurement"`
	Units       string  `json:"units" cassandra:"units"`
	Value       float64 `json:"value" cassandra:"value"`
}

// CassandraDataAccess implements the SensorDataAccess interface for Cassandra
type CassandraDataAccess struct {
	Session *gocql.Session
}

// NewCassandraDataAccess creates a new CassandraDataAccess instance and establishes connection
func NewCassandraDataAccess() (*CassandraDataAccess, error) {
	// Suppress gocql internal logging to reduce noise
	gocql.Logger = log.New(io.Discard, "", 0)

	// Connect to Cassandra
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "system"
	cluster.Consistency = gocql.Quorum
	cluster.DisableInitialHostLookup = true
	cluster.ProtoVersion = 4
	cluster.NumConns = 2
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	fmt.Println("Connected to Cassandra")

	query := `CREATE KEYSPACE IF NOT EXISTS sensor_data_db WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`
	if err := session.Query(query).Exec(); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}
	fmt.Println("Keyspace created or already exists")
	session.Close()

	cluster.Keyspace = "sensor_data_db"
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sensor_data_db keyspace: %w", err)
	}

	query = `CREATE TABLE IF NOT EXISTS sensor_data (
		recorded bigint,
		location text,
		sensor text,
		measurement text,
		units text,
		value double,
		PRIMARY KEY (recorded, location, sensor, measurement)
	)`
	if err := session.Query(query).Exec(); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	fmt.Println("Table created or already exists")

	return &CassandraDataAccess{Session: session}, nil
}

func (c *CassandraDataAccess) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}

// parseRecorded converts the recorded field to int64 whether it's a string or number
func parseRecorded(recorded interface{}) (int64, error) {
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

func (c *CassandraDataAccess) LogSensorData(jsonData string) error {
	var jsonStruct SensorDataJSON
	if err := json.Unmarshal([]byte(jsonData), &jsonStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}

	recorded, err := parseRecorded(jsonStruct.Recorded)
	if err != nil {
		return fmt.Errorf("failed to parse recorded field: %w", err)
	}

	sensorData := SensorDataCassandra{
		Recorded:    recorded,
		Location:    jsonStruct.Location,
		Sensor:      jsonStruct.Sensor,
		Measurement: jsonStruct.Measurement,
		Units:       jsonStruct.Units,
		Value:       jsonStruct.Value,
	}

	query := "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
	if err := c.Session.Query(query,
		sensorData.Recorded,
		sensorData.Location,
		sensorData.Sensor,
		sensorData.Measurement,
		sensorData.Units,
		sensorData.Value,
	).Exec(); err != nil {
		return fmt.Errorf("failed to insert into sensor_data: %w", err)
	}

	return nil
}

func (c *CassandraDataAccess) FetchSensorData() ([]string, error) {
	query := "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data"
	iter := c.Session.Query(query).Iter()
	var recorded int64
	var location, sensor, measurement, units string
	var value float64
	var results []string

	for iter.Scan(&recorded, &location, &sensor, &measurement, &units, &value) {
		sensorData := SensorDataCassandra{
			Recorded:    recorded,
			Location:    location,
			Sensor:      sensor,
			Measurement: measurement,
			Units:       units,
			Value:       value,
		}
		jsonBytes, err := json.Marshal(sensorData)
		if err != nil {
			fmt.Printf("failed to marshal SensorDataCassandra to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close iterator: %w", err)
	}

	return results, nil
}

func (c *CassandraDataAccess) PurgeSensorData() error {
	query := "TRUNCATE sensor_data"
	if err := c.Session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to truncate sensor_data table: %w", err)
	}
	return nil
}
