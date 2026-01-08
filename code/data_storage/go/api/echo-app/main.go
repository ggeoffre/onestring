// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
	"net/http"
	"fmt"
	"strings"
	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusOK, requestBody)
	})

	e.GET("/report", func(c echo.Context) error {
		csvData := jsonToCSV(jsonData)
		c.Response().Header().Set(echo.HeaderContentType, "text/csv")
		c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="report.csv"`)
		return c.String(http.StatusOK, csvData)
	})

	e.Any("/purge", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Purge executed"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
