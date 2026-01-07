// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.database;

import com.geoffrey.database.SensorDataJsonHelper;
import java.sql.*;
import java.util.*;
import org.json.*;

public class PostgresData {

    // Configuration constants
    private static final String HOST = "192.168.1.60";
    private static final int PORT = 5432;
    private static final String USER = "postgres";
    private static final String PASSWORD = "";
    private static final String DATABASE = "sensor_data_db";
    private static final String TABLE = "sensor_data";

    public static void processSensorData() {
        System.out.println("POSTGRES");
        System.out.println("########");

        // Step 1: Connect to the default "postgres" database to create the target database if needed
        String urlPostgres = String.format(
            "jdbc:postgresql://%s:%d/postgres",
            HOST,
            PORT
        );
        try (
            Connection connection = DriverManager.getConnection(
                urlPostgres,
                USER,
                PASSWORD
            )
        ) {
            System.out.println("Connected to Postgres (postgres database)");

            // Check if database exists, create if not
            try (Statement stmt = connection.createStatement()) {
                String checkDbSql = String.format(
                    "SELECT 1 FROM pg_database WHERE datname = '%s'",
                    DATABASE
                );
                ResultSet rs = stmt.executeQuery(checkDbSql);
                boolean dbExists = rs.next();
                rs.close();

                if (!dbExists) {
                    String createDbSql = String.format(
                        "CREATE DATABASE %s;",
                        DATABASE
                    );
                    stmt.executeUpdate(createDbSql);
                    System.out.println("Database created: " + DATABASE);
                } else {
                    System.out.println("Database already exists: " + DATABASE);
                }
            } catch (SQLException e) {
                System.err.println(
                    "Error creating database: " + e.getMessage()
                );
                e.printStackTrace();
                return;
            }
        } catch (SQLException e) {
            System.err.println(
                "Failed to connect to the Postgres server (postgres database)."
            );
            e.printStackTrace();
            return;
        }

        // Step 2: Connect to the target database for all further operations
        String urlTargetDb = String.format(
            "jdbc:postgresql://%s:%d/%s",
            HOST,
            PORT,
            DATABASE
        );
        try (
            Connection connection = DriverManager.getConnection(
                urlTargetDb,
                USER,
                PASSWORD
            )
        ) {
            System.out.println("Connected to Postgres (" + DATABASE + ")");

            // Create table if not exists
            try (Statement stmt = connection.createStatement()) {
                String createTableSql = String.format(
                    "CREATE TABLE IF NOT EXISTS %s (" +
                        "recorded BIGINT NOT NULL, " +
                        "location VARCHAR(255) NOT NULL, " +
                        "sensor VARCHAR(255) NOT NULL, " +
                        "measurement VARCHAR(255) NOT NULL, " +
                        "units VARCHAR(255) NOT NULL, " +
                        "value NUMERIC(10, 2) NOT NULL" +
                        ");",
                    TABLE
                );
                stmt.executeUpdate(createTableSql);
                System.out.println("Table ensured: " + TABLE);
            } catch (SQLException e) {
                System.err.println("Error creating table: " + e.getMessage());
                e.printStackTrace();
                return;
            }

            // Store sensor data
            String insertSql =
                "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)";
            try {
                JSONObject jsonObj =
                    SensorDataJsonHelper.getSensorDataJSONObject();
                System.out.println("SensorData marshalled to JSON");

                // Validate or generate 'recorded' field
                long recorded = jsonObj.has("recorded")
                    ? jsonObj.getLong("recorded")
                    : System.currentTimeMillis() / 1000;

                try (
                    PreparedStatement pstmt = connection.prepareStatement(
                        insertSql
                    )
                ) {
                    pstmt.setLong(1, recorded);
                    pstmt.setString(2, jsonObj.getString("location"));
                    pstmt.setString(3, jsonObj.getString("sensor"));
                    pstmt.setString(4, jsonObj.getString("measurement"));
                    pstmt.setString(5, jsonObj.getString("units"));
                    pstmt.setBigDecimal(6, jsonObj.getBigDecimal("value"));
                    pstmt.executeUpdate();
                }
                System.out.println("SensorData stored");
            } catch (SQLException | JSONException e) {
                System.err.println("Failed to store sensor data.");
                e.printStackTrace();
            }

            // Retrieve sensor data
            String selectSql =
                "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data";
            List<String> jsonList = new ArrayList<>();
            try (
                Statement stmt = connection.createStatement();
                ResultSet rs = stmt.executeQuery(selectSql)
            ) {
                while (rs.next()) {
                    JSONObject jsonObj = new JSONObject();
                    jsonObj.put("recorded", rs.getLong("recorded"));
                    jsonObj.put("location", rs.getString("location"));
                    jsonObj.put("sensor", rs.getString("sensor"));
                    jsonObj.put("measurement", rs.getString("measurement"));
                    jsonObj.put("units", rs.getString("units"));
                    jsonObj.put("value", rs.getBigDecimal("value"));
                    jsonList.add(jsonObj.toString());
                }
                System.out.println("SensorData retrieved");
            } catch (SQLException | JSONException e) {
                System.err.println("Failed to retrieve sensor data.");
                e.printStackTrace();
            }

            // Print as CSV
            try {
                System.out.print(
                    SensorDataJsonHelper.jsonArrayToCsv(
                        jsonList.toArray(new String[0])
                    )
                );
            } catch (Exception e) {
                System.err.println("Failed to convert JSON to CSV.");
                e.printStackTrace();
            }

            // Purge sensor data (use with caution in production!)
            String deleteSql = "DELETE FROM sensor_data";
            try (Statement stmt = connection.createStatement()) {
                stmt.executeUpdate(deleteSql);
                System.out.println("Sensor Data purged\n");
            } catch (SQLException e) {
                System.err.println("Failed to purge sensor data.");
                e.printStackTrace();
            }
        } catch (SQLException e) {
            System.err.println("Failed to connect to the Postgres database.");
            e.printStackTrace();
        }
    }
}
