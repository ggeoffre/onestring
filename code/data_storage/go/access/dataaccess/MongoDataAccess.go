// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package dataaccess

import "fmt"

// MongoDataAccess implements the SensorDataAccess interface for MongoDB
type MongoDataAccess struct{}

// NewMongoDataAccess creates a new MongoDataAccess instance
func NewMongoDataAccess() *MongoDataAccess {
    return &MongoDataAccess{}
}

func (m *MongoDataAccess) LogSensorData(jsonData string) error {
    fmt.Printf("Logging sensor data to MongoDB: %s\n", jsonData)
    return nil
}

func (m *MongoDataAccess) FetchSensorData() ([]string, error) {
    fmt.Println("Fetching sensor data from MongoDB")
    return []string{`{"sensor":"humidity","value":45.6}`}, nil
}

func (m *MongoDataAccess) PurgeSensorData() error {
    fmt.Println("Purging sensor data from MongoDB")
    return nil
}
