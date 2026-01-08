// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import "fmt"

// SensorDataAccess defines the interface for data access operations
type SensorDataAccess interface {
    LogSensorData(jsonData string) error
    FetchSensorData() ([]string, error)
    PurgeSensorData() error
}

// Helper function to handle unsupported data access types
func UnsupportedDataAccess(dataAccessType string) {
    panic(fmt.Sprintf("Unsupported DATA_ACCESS type: %s", dataAccessType))
}
