// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
    "net/http"
	"fmt"
	"strings"
    "github.com/gin-gonic/gin"
)

var jsonData = map[string]interface{}{
	"recorded":    1756655999,
	"location":    "den",
	"sensor":      "bmp280",
	"measurement": "temperature",
	"units":       "C",
	"value":       22.3,
}

func jsonToCSV(jsonData map[string]interface{}) string {
	var csvData string
	var keys []string
	var values []string

	for key, value := range jsonData {
		keys = append(keys, key)
		values = append(values, fmt.Sprintf("%v", value))
	}

	csvData += strings.Join(keys, ",") + "\n"
	csvData += strings.Join(values, ",")
	return csvData
}

func main() {
    // Create a new Gin router
    r := gin.Default()

    // Define a GET route for the root path
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Gin API Server is running!",
        })
    })
    // Define a POST route for /echo
    r.POST("/echo", func(c *gin.Context) {
        var requestBody map[string]interface{}
        if err := c.ShouldBindJSON(&requestBody); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, requestBody)
    })

    // Define a POST route for /log
    r.POST("/log", func(c *gin.Context) {
        var requestBody map[string]interface{}
        if err := c.ShouldBindJSON(&requestBody); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, requestBody)
    })

    // Define a GET route for /report
    r.GET("/report", func(c *gin.Context) {
        csvData := jsonToCSV(jsonData)
        c.Header("Content-Disposition", "attachment; filename=report.csv")
        c.Header("Content-Type", "text/csv")
        c.String(http.StatusOK, csvData)
    })

    // Define a route for /purge that supports both GET and POST
    r.Any("/purge", func(c *gin.Context) {
        if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodPost {
            c.JSON(http.StatusOK, gin.H{"status": "purged"})
        } else {
            c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
        }
    })
    // Start the server on port 8080
    r.Run(":8080")
}
