// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package dataaccess

import "fmt"

// MySQLDataAccess implements the SensorDataAccess interface for MySQL
type MySQLDataAccess struct{}

// NewMySQLDataAccess creates a new MySQLDataAccess instance
func NewMySQLDataAccess() *MySQLDataAccess {
    return &MySQLDataAccess{}
}

func (m *MySQLDataAccess) LogSensorData(jsonData string) error {
    fmt.Printf("Logging sensor data to MySQL: %s\n", jsonData)
    return nil
}

func (m *MySQLDataAccess) FetchSensorData() ([]string, error) {
    fmt.Println("Fetching sensor data from MySQL")
    return []string{`{"sensor":"light","value":300}`}, nil
}

func (m *MySQLDataAccess) PurgeSensorData() error {
    fmt.Println("Purging sensor data from MySQL")
    return nil
}
