// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.data;

import java.util.*;
import redis.clients.jedis.Jedis;

public class RedisData implements SensorDataAccess {

    // Configuration constants (consider externalizing for production)
    private static final String REDIS_HOST = "192.168.1.60";
    private static final int REDIS_PORT = 6379;
    private static final String SENSOR_DATA_KEY = "sensor_data_db:sensor_data";

    /**
     * Processes sensor data: inserts, retrieves, prints as CSV, and purges the list.
     */
    public Jedis getJedisClient() {
        try {
            Jedis jedisClient = new Jedis(REDIS_HOST, REDIS_PORT);
            System.out.println("Connected to Redis database.");
            return jedisClient;
        } catch (Exception e) {
            System.err.println("Failed to connect to Redis database.");
            e.printStackTrace();
            throw e;
        }
    }

    // Store sensor data
    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to Redis: " + jsonData);
        try {
            Jedis jedisClient = getJedisClient();
            jedisClient.rpush(SENSOR_DATA_KEY, jsonData);
            System.out.println("SensorData stored");
        } catch (Exception e) {
            System.err.println("Failed to store sensor data in Redis.");
            e.printStackTrace();
        }
    }

    // Retrieve sensor data
    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from Redis");
        List<String> jsonStrs = null;
        try {
            Jedis jedisClient = getJedisClient();
            jsonStrs = jedisClient.lrange(SENSOR_DATA_KEY, 0, -1);
        } catch (Exception e) {
            System.err.println("Failed to retrieve sensor data from Redis.");
            e.printStackTrace();
        }
        return jsonStrs != null ? jsonStrs : Collections.emptyList();
    }

    // Purge sensor data (use with caution in production!)
    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from Redis");
        try {
            Jedis jedisClient = getJedisClient();
            jedisClient.del(SENSOR_DATA_KEY);
            System.out.println("Sensor Data purged\n");
        } catch (Exception e) {
            System.err.println("Failed to purge sensor data in Redis.");
            e.printStackTrace();
        }
    }
}
