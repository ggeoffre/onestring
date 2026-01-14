// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.database;

import com.geoffrey.database.SensorDataJsonHelper;
import java.util.List;
import redis.clients.jedis.Jedis;

public class RedisData {

    // Configuration constants (consider externalizing for production)
    private static final String REDIS_HOST = "localhost";
    private static final int REDIS_PORT = 6379;
    private static final String SENSOR_DATA_KEY = "sensor_data_db:sensor_data";

    /**
     * Processes sensor data: inserts, retrieves, prints as CSV, and purges the list.
     */
    public static void processSensorData() {
        System.out.println("REDIS");
        System.out.println("#####");

        // Try-with-resources ensures Jedis client is always closed
        try (Jedis jedisClient = new Jedis(REDIS_HOST, REDIS_PORT)) {
            System.out.println("Connected to Redis database.");

            // Store sensor data
            try {
                jedisClient.rpush(
                    SENSOR_DATA_KEY,
                    SensorDataJsonHelper.getSensorDataJSONString()
                );
                System.out.println("SensorData stored");
            } catch (Exception e) {
                System.err.println("Failed to store sensor data in Redis.");
                e.printStackTrace();
            }

            // Retrieve sensor data
            try {
                List<String> jsonStrs = jedisClient.lrange(
                    SENSOR_DATA_KEY,
                    0,
                    -1
                );
                System.out.println("SensorData retrieved");
                System.out.println(
                    SensorDataJsonHelper.jsonArrayToCsv(
                        jsonStrs.toArray(new String[0])
                    )
                );
            } catch (Exception e) {
                System.err.println(
                    "Failed to retrieve sensor data from Redis."
                );
                e.printStackTrace();
            }

            // Purge sensor data (use with caution in production!)
            try {
                jedisClient.del(SENSOR_DATA_KEY);
                System.out.println("Sensor Data purged\n");
            } catch (Exception e) {
                System.err.println("Failed to purge sensor data in Redis.");
                e.printStackTrace();
            }
        } catch (Exception e) {
            System.err.println("Failed to connect to Redis.");
            e.printStackTrace();
        }
    }
}
