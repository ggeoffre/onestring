// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package dataaccess

import "fmt"

// PostgresDataAccess implements the SensorDataAccess interface for PostgreSQL
type PostgresDataAccess struct{}

// NewPostgresDataAccess creates a new PostgresDataAccess instance
func NewPostgresDataAccess() *PostgresDataAccess {
    return &PostgresDataAccess{}
}

func (p *PostgresDataAccess) LogSensorData(jsonData string) error {
    fmt.Printf("Logging sensor data to PostgreSQL: %s\n", jsonData)
    return nil
}

func (p *PostgresDataAccess) FetchSensorData() ([]string, error) {
    fmt.Println("Fetching sensor data from PostgreSQL")
    return []string{`{"sensor":"sound","value":75}`}, nil
}

func (p *PostgresDataAccess) PurgeSensorData() error {
    fmt.Println("Purging sensor data from PostgreSQL")
    return nil
}
