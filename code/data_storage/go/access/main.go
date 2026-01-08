// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
    "fmt"
    "os"
	"access/dataaccess"
)

func getDataAccess() dataaccess.SensorDataAccess {
    dataAccessType := os.Getenv("DATA_ACCESS")
    if dataAccessType == "" {
        dataAccessType = "redis"
    }

    switch dataAccessType {
    case "redis":
        return dataaccess.NewRedisDataAccess()
    case "mongo":
        return dataaccess.NewMongoDataAccess()
    case "cassandra":
        return dataaccess.NewCassandraDataAccess()
    case "mysql":
        return dataaccess.NewMySQLDataAccess()
    case "postgres":
        return dataaccess.NewPostgresDataAccess()
    default:
        dataaccess.UnsupportedDataAccess(dataAccessType)
        return nil
    }
}

func main() {
    dataAccess := getDataAccess()

    // Log sensor data
    err := dataAccess.LogSensorData(`{"sensor":"temperature","value":22.3}`)
    if err != nil {
        fmt.Println("Error logging sensor data:", err)
    }

    // Fetch sensor data
    data, err := dataAccess.FetchSensorData()
    if err != nil {
        fmt.Println("Error fetching sensor data:", err)
    } else {
        fmt.Println("Fetched sensor data:", data)
    }

    // Purge sensor data
    err = dataAccess.PurgeSensorData()
    if err != nil {
        fmt.Println("Error purging sensor data:", err)
    }
}
