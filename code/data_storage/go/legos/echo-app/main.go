// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
	"os"
	"net/http"
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
	"echo-app/data"
)

// Global data access instance - created once, reused for all requests
var dataAccess data.SensorDataAccess

func initDataAccess() error {
	dataAccessType := os.Getenv("DATA_ACCESS")
	if dataAccessType == "" {
		dataAccessType = "redis"
	}

	var err error
	switch dataAccessType {
	case "cassandra":
		dataAccess, err = data.NewCassandraDataAccess()
	case "mongo":
		dataAccess, err = data.NewMongoDataAccess()
	case "mysql":
		dataAccess, err = data.NewMySQLDataAccess()
	case "postgres":
		dataAccess, err = data.NewPostgresDataAccess()
	case "redis":
		dataAccess, err = data.NewRedisDataAccess()
	default:
        data.UnsupportedDataAccess(dataAccessType)
        return nil
	}

	if err != nil {
		return err
	}

	log.Printf("Data access initialized: %s", dataAccessType)
	return nil
}

func main() {
	// Initialize database connection once at startup
	if err := initDataAccess(); err != nil {
		log.Fatalf("Failed to initialize data access: %v", err)
	}

	// Clean up connection on exit
	if closer, ok := dataAccess.(interface{ Close() }); ok {
		defer closer.Close()
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Echo API Server is running!"})
	})

	e.POST("/echo", func(c echo.Context) error {
		var requestBody map[string]interface{}
		if err := c.Bind(&requestBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		return c.JSON(http.StatusOK, requestBody)
	})

	e.POST("/log", func(c echo.Context) error {
		var requestBody map[string]interface{}
		if err := c.Bind(&requestBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Marshal to JSON string
		bytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("Failed to marshal request body: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		}

		// Log to database using the global connection
		if err := dataAccess.LogSensorData(string(bytes)); err != nil {
			log.Printf("Failed to log sensor data: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to log sensor data"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Data logged successfully"})
	})

	e.GET("/report", func(c echo.Context) error {
		// Fetch data using the global connection
		sensorData, err := dataAccess.FetchSensorData()
		if err != nil {
			log.Printf("Failed to fetch sensor data: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sensor data"})
		}

		// Convert JSON array to CSV
		csvData, err := data.JsonToCSV(sensorData)
		if err != nil {
			log.Printf("Failed to convert to CSV: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate CSV"})
		}

		c.Response().Header().Set(echo.HeaderContentType, "text/csv")
		c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="report.csv"`)
		return c.String(http.StatusOK, csvData)
	})

	e.Any("/purge", func(c echo.Context) error {
		// Purge using the global connection
		if err := dataAccess.PurgeSensorData(); err != nil {
			log.Printf("Failed to purge sensor data: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to purge sensor data"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Purge executed successfully"})
	})

	log.Println("Starting server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
