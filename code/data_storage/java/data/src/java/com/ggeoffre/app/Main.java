// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.app;

import com.geoffrey.database.*;

public class Main {

    public static void main(String[] args) {
        try {
            CassandraData.processSensorData();
        } catch (Exception e) {
            System.out.println(
                "Error connecting to Cassandra: " + e.getMessage()
            );
        }
        try {
            MongoData.processSensorData();
        } catch (Exception e) {
            System.out.println(
                "Error connecting to MongoDB: " + e.getMessage()
            );
        }
        try {
            MySQLData.processSensorData();
        } catch (Exception e) {
            System.out.println("Error connecting to MySQL: " + e.getMessage());
        }
        try {
            PostgresData.processSensorData();
        } catch (Exception e) {
            System.out.println(
                "Error connecting to PostgreSQL: " + e.getMessage()
            );
        }
        try {
            RedisData.processSensorData();
        } catch (Exception e) {
            System.out.println("Error connecting to Redis: " + e.getMessage());
        }
    }
}
