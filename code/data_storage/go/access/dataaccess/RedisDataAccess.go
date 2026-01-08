// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package dataaccess

import "fmt"

// RedisDataAccess implements the SensorDataAccess interface for Redis
type RedisDataAccess struct{}

// NewRedisDataAccess creates a new RedisDataAccess instance
func NewRedisDataAccess() *RedisDataAccess {
    return &RedisDataAccess{}
}

func (r *RedisDataAccess) LogSensorData(jsonData string) error {
    fmt.Printf("Logging sensor data to Redis: %s\n", jsonData)
    return nil
}

func (r *RedisDataAccess) FetchSensorData() ([]string, error) {
    fmt.Println("Fetching sensor data from Redis")
    return []string{`{"sensor":"temperature","value":22.3}`}, nil
}

func (r *RedisDataAccess) PurgeSensorData() error {
    fmt.Println("Purging sensor data from Redis")
    return nil
}
