// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.data;

import java.sql.*;
import java.util.*;
import org.json.*;

public class MySQLData implements SensorDataAccess {

    // Configuration constants
    private static final String HOST = "localhost";
    private static final int PORT = 3306;
    private static final String USER = "root";
    private static final String PASSWORD = "";
    private static final String DATABASE = "sensor_data_db";
    private static final String TABLE = "sensor_data";

    private Connection createConnection() throws SQLException {
        String url = String.format(
            "jdbc:mysql://%s:%d/%s",
            HOST,
            PORT,
            DATABASE
        );
        return DriverManager.getConnection(url, USER, PASSWORD);
    }

    private Connection createSystemDbConnection() throws SQLException {
        String url = String.format("jdbc:mysql://%s:%d/mysql", HOST, PORT);
        return DriverManager.getConnection(url, USER, PASSWORD);
    }

    // Initialize the MySQL connection
    private void init() throws SQLException {
        // Connect to the system database
        try (Connection connection = createSystemDbConnection()) {
            System.out.println("Connected to MySQL (mysql database)");

            // Create database if not exists
            try (Statement stmt = connection.createStatement()) {
                String createDbSql = String.format(
                    "CREATE DATABASE IF NOT EXISTS %s;",
                    DATABASE
                );
                stmt.executeUpdate(createDbSql);
                System.out.println("Database ensured: " + DATABASE);
            } catch (SQLException e) {
                System.err.println(
                    "Error creating database: " + e.getMessage()
                );
                e.printStackTrace();
                return;
            }
        } catch (SQLException e) {
            System.err.println(
                "Failed to connect to the MySQL server (mysql database)."
            );
            e.printStackTrace();
            return;
        }

        //connect to target database
        try (
            Connection connection = createConnection();
            Statement stmt = connection.createStatement()
        ) {
            System.out.println("Connected to MySQL (" + DATABASE + ")");

            // Create table if not exists
            String createTableSql = String.format(
                "CREATE TABLE IF NOT EXISTS %s (" +
                    "recorded BIGINT NOT NULL, " +
                    "location VARCHAR(255) NOT NULL, " +
                    "sensor VARCHAR(255) NOT NULL, " +
                    "measurement VARCHAR(255) NOT NULL, " +
                    "units VARCHAR(10) NOT NULL, " +
                    "value DECIMAL(5,2) NOT NULL" +
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
    }

    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to MySQL: " + jsonData);

        String insertSql =
            "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)";
        try (Connection connection = createConnection()) {
            JSONObject jsonObj = new JSONObject(jsonData);
            System.out.println("SensorData marshalled to JSON");

            // Validate or generate 'recorded' field
            long recorded = jsonObj.has("recorded")
                ? jsonObj.getLong("recorded")
                : System.currentTimeMillis() / 1000;

            try (
                PreparedStatement pstmt = connection.prepareStatement(insertSql)
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
    }

    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from MySQL");
        String selectSql =
            "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data";
        List<String> jsonStrings = new ArrayList<>();
        try (
            Connection connection = createConnection();
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
                jsonStrings.add(jsonObj.toString());
            }
            System.out.println("SensorData retrieved");
        } catch (SQLException | JSONException e) {
            System.err.println("Failed to retrieve sensor data.");
            e.printStackTrace();
        }
        return jsonStrings;
    }

    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from MySQL");
        String deleteSql = "DELETE FROM sensor_data";
        try (
            Connection connection = createConnection();
            Statement stmt = connection.createStatement()
        ) {
            stmt.executeUpdate(deleteSql);
            System.out.println("Sensor Data purged\n");
        } catch (SQLException e) {
            System.err.println("Failed to purge sensor data.");
            e.printStackTrace();
        }
    }
}
