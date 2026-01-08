// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
	"os"
	"net/http"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gin-app/data"
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

	// Create a new Gin router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Gin API Server is running!",
		})
	})

	r.POST("/echo", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		c.JSON(http.StatusOK, requestBody)
	})

	r.POST("/log", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Marshal to JSON string
		bytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("Failed to marshal request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Log to database using the global connection
		if err := dataAccess.LogSensorData(string(bytes)); err != nil {
			log.Printf("Failed to log sensor data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log sensor data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data logged successfully"})
	})

	r.GET("/report", func(c *gin.Context) {
		// Fetch data using the global connection
		sensorData, err := dataAccess.FetchSensorData()
		if err != nil {
			log.Printf("Failed to fetch sensor data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sensor data"})
			return
		}

		// Convert JSON array to CSV
		csvData, err := data.JsonToCSV(sensorData)
		if err != nil {
			log.Printf("Failed to convert to CSV: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSV"})
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", `attachment; filename="report.csv"`)
		c.String(http.StatusOK, csvData)
	})

	r.Any("/purge", func(c *gin.Context) {
		// Purge using the global connection
		if err := dataAccess.PurgeSensorData(); err != nil {
			log.Printf("Failed to purge sensor data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purge sensor data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Purge executed successfully"})
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
