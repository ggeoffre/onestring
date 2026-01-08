// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// RedisData struct for Redis
type RedisData struct {
	Client *redis.Client
}

// NewRedisData creates a new RedisData instance
func NewRedisData() *RedisData {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.60:6379",
		Password: "",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	// Note: redis.NewClient doesn't actually connect until a command is run,
	// but strictly speaking the 'connection object' is created here.
	return &RedisData{Client: client}
}

// Close closes the Redis client
func (c *RedisData) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}

func (c *RedisData) RawRedisData() {
	ctx := context.Background()
	key := "sensor_data:readings"

	// Use RPush to append the JSON string to the list
	if err := c.Client.RPush(ctx, key, getRandomSensorDataJsonString()).Err(); err != nil {
		fmt.Printf("failed to rpush to Redis: %v\n", err)
		return
	}
	fmt.Println("Successfully rpushed data to Redis")

	// Use LRange to retrieve the list of JSON strings
	values, err := c.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Printf("failed to lrange from Redis: %v\n", err)
		return
	}

	// Join results with commas to form a valid JSON array
	if len(values) > 0 {
		data := []string{"[" + strings.Join(values, ",") + "]"}
		// Convert JSON array to CSV string
		csvString, err := jsonToCSV(strings.Join(data, "\n"))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s", "Successfully lranged data from Redis", csvString)
		}
	}

	// Delete the key
	if err := c.Client.Del(context.Background(), key).Err(); err != nil {
		fmt.Printf("failed to delete key from Redis: %v\n", err)
		return
	}
	fmt.Println("Successfully deleted key from Redis")
}
