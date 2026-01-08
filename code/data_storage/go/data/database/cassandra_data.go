// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
)

type SensorDataCassandra struct {
	Recorded    int64   `cassandra:"recorded"`
	Location    string  `cassandra:"location"`
	Sensor      string  `cassandra:"sensor"`
	Measurement string  `cassandra:"measurement"`
	Units       string  `cassandra:"units"`
	Value       float64 `cassandra:"value"`
}

// CassandraData struct for Cassandra
type CassandraData struct {
	Session *gocql.Session
}

// NewCassandraData creates a new CassandraData instance and establishes connection
func NewCassandraData() (*CassandraData, error) {
	// Connect to Cassandra
	cluster := gocql.NewCluster("192.168.1.60")
	cluster.Port = 9042
	cluster.Keyspace = "system" // Connect to system to create keyspace
	cluster.Consistency = gocql.Quorum
	cluster.DisableInitialHostLookup = true

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %w", err)
	}

	fmt.Println("Connected to Cassandra")

	// Create the keyspace if it doesn't exist
	var query = `SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name = ?;`
	var kspName string
	if err := session.Query(query, "sensor_data_db").Scan(&kspName); err != nil {
		if err != gocql.ErrNotFound {
			session.Close()
			return nil, fmt.Errorf("database error checking keyspace: %w", err)
		}
	}

	if kspName == "" {
		query = `CREATE KEYSPACE IF NOT EXISTS sensor_data_db WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`
		if err := session.Query(query).Exec(); err != nil {
			session.Close()
			return nil, fmt.Errorf("failed to create keyspace: %w", err)
		}
		fmt.Println("Keyspace created")
	}
	session.Close()

	// Reconnect to the specific Keyspace
	cluster.Keyspace = "sensor_data_db"
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sensor_data_db keyspace: %w", err)
	}

	return &CassandraData{Session: session}, nil
}

// Close cleans up the session
func (c *CassandraData) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}

func (c *CassandraData) RawCassandraData() {
	// Create the table if it doesn't exist
	query := `CREATE TABLE IF NOT EXISTS sensor_data (
			 recorded bigint,
			 location text,
			 sensor text,
			 measurement text,
			 units text,
			 value double,
			 PRIMARY KEY (recorded, location, sensor, measurement)
			);`
	if err := c.Session.Query(query).Exec(); err != nil {
		fmt.Printf("Failed to create table: %v\n", err)
		return
	}
	fmt.Println("Table created or already exists")

	// Convert JSON string to SensorDataCassandra struct
	var sensorDataCassandra SensorDataCassandra
	if err := json.Unmarshal([]byte(getRandomSensorDataJsonString()), &sensorDataCassandra); err != nil {
		fmt.Printf("Failed to unmarshal JSON string: %v\n", err)
		return
	}
	fmt.Println("JSON string converted to SensorDataCassandra struct")

	// Insert the SensorData into the "sensor_data" table
	query = "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
	if err := c.Session.Query(query,
		sensorDataCassandra.Recorded,
		sensorDataCassandra.Location,
		sensorDataCassandra.Sensor,
		sensorDataCassandra.Measurement,
		sensorDataCassandra.Units,
		sensorDataCassandra.Value,
	).Exec(); err != nil {
		fmt.Printf("Failed to insert into sensor_data: %v\n", err)
		return
	}
	fmt.Println("SensorDataCassandra stored")

	// Retrieve all SensorData from the "sensor_data" table
	query = "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data"
	iter := c.Session.Query(query).Iter()
	var recorded int64
	var location, sensor, measurement, units string
	var value float64
	var results []string

	// Iterate over the results
	for iter.Scan(&recorded, &location, &sensor, &measurement, &units, &value) {
		sensorData := SensorDataCassandra{
			Recorded:    recorded,
			Location:    location,
			Sensor:      sensor,
			Measurement: measurement,
			Units:       units,
			Value:       value,
		}
		// Marshal the SensorDataCassandra struct to JSON
		jsonBytes, err := json.Marshal(sensorData)
		if err != nil {
			fmt.Printf("failed to marshal SensorDataCassandra to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if err := iter.Close(); err != nil {
		fmt.Printf("failed to close iterator: %v\n", err)
		return
	}

	// Join results with commas to form a valid JSON array
	if len(results) > 0 {
		data := []string{"[" + strings.Join(results, ",") + "]"}

		// Convert JSON array to CSV string
		csvString, err := jsonToCSV(strings.Join(data, "\n"))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s", "SensorDataCassandra retrieved", csvString)
		}
	}

	// Truncate the sensor_data table
	query = "TRUNCATE sensor_data"
	if err := c.Session.Query(query).Exec(); err != nil {
		fmt.Printf("Failed to truncate sensor_data table: %v\n", err)
		return
	}
	fmt.Println("SensorDataCassandra purged")
}
