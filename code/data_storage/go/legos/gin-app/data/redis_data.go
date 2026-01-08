// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// SensorDataRedisJSON is for unmarshaling JSON where recorded might be a string
type SensorDataRedisJSON struct {
	Recorded    interface{} `json:"recorded"` // Accept both string and number
	Location    string      `json:"location"`
	Sensor      string      `json:"sensor"`
	Measurement string      `json:"measurement"`
	Units       string      `json:"units"`
	Value       float64     `json:"value"`
}

type SensorDataRedis struct {
	Recorded    int64   `json:"recorded"`
	Location    string  `json:"location"`
	Sensor      string  `json:"sensor"`
	Measurement string  `json:"measurement"`
	Units       string  `json:"units"`
	Value       float64 `json:"value"`
}

// RedisDataAccess implements the SensorDataAccess interface for Redis
type RedisDataAccess struct {
	Client *redis.Client
	Key    string
}

// NewRedisDataAccess creates a new RedisDataAccess instance
func NewRedisDataAccess() (*RedisDataAccess, error) {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.60:6379",
		Password: "",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis")

	return &RedisDataAccess{
		Client: client,
		Key:    "sensor_data:readings",
	}, nil
}

// Close closes the Redis client
func (r *RedisDataAccess) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}

// parseRecorded converts the recorded field to int64 whether it's a string or number
func parseRecordedRedis(recorded interface{}) (int64, error) {
	switch v := recorded.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("recorded field has unexpected type: %T", v)
	}
}

func (r *RedisDataAccess) LogSensorData(jsonData string) error {
	var jsonStruct SensorDataRedisJSON
	if err := json.Unmarshal([]byte(jsonData), &jsonStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}

	recorded, err := parseRecordedRedis(jsonStruct.Recorded)
	if err != nil {
		return fmt.Errorf("failed to parse recorded field: %w", err)
	}

	sensorData := SensorDataRedis{
		Recorded:    recorded,
		Location:    jsonStruct.Location,
		Sensor:      jsonStruct.Sensor,
		Measurement: jsonStruct.Measurement,
		Units:       jsonStruct.Units,
		Value:       jsonStruct.Value,
	}

	// Marshal back to JSON for storage
	jsonBytes, err := json.Marshal(sensorData)
	if err != nil {
		return fmt.Errorf("failed to marshal sensor data: %w", err)
	}

	ctx := context.Background()
	if err := r.Client.RPush(ctx, r.Key, string(jsonBytes)).Err(); err != nil {
		return fmt.Errorf("failed to insert into Redis: %w", err)
	}

	return nil
}

func (r *RedisDataAccess) FetchSensorData() ([]string, error) {
	ctx := context.Background()

	values, err := r.Client.LRange(ctx, r.Key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Redis: %w", err)
	}

	return values, nil
}

func (r *RedisDataAccess) PurgeSensorData() error {
	ctx := context.Background()

	if err := r.Client.Del(ctx, r.Key).Err(); err != nil {
		return fmt.Errorf("failed to delete data from Redis: %w", err)
	}

	return nil
}
