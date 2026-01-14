// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.database;

import com.datastax.oss.driver.api.core.CqlSession;
import com.datastax.oss.driver.api.core.cql.*;
import java.net.InetSocketAddress;
import java.util.*;
import org.json.*;

public class CassandraData {

    // Configuration constants
    private static final String CONTACT_POINT = "localhost";
    private static final int PORT = 9042;
    private static final String LOCAL_DATACENTER = "datacenter1";
    private static final String KEYSPACE = "sensor_data_db";
    private static final String TABLE = "sensor_data";

    public static void processSensorData() {
        System.out.println("CASSANDRA");
        System.out.println("#########");

        // Step 1: Connect to the "system" keyspace to create the keyspace if needed
        try (
            CqlSession session = CqlSession.builder()
                .addContactPoint(new InetSocketAddress(CONTACT_POINT, PORT))
                .withLocalDatacenter(LOCAL_DATACENTER)
                .withKeyspace("system")
                .build()
        ) {
            System.out.println("Connected to Cassandra (system keyspace)");

            // Create keyspace if not exists
            String createKeyspace = String.format(
                "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };",
                KEYSPACE
            );
            session.execute(createKeyspace);
            System.out.println("Keyspace ensured: " + KEYSPACE);
        } catch (Exception e) {
            System.err.println(
                "Failed to connect to Cassandra (system keyspace) or create keyspace."
            );
            e.printStackTrace();
            return;
        }

        // Step 2: Connect to the target keyspace for all further operations
        try (
            CqlSession session = CqlSession.builder()
                .addContactPoint(new InetSocketAddress(CONTACT_POINT, PORT))
                .withLocalDatacenter(LOCAL_DATACENTER)
                .withKeyspace(KEYSPACE)
                .build()
        ) {
            System.out.println(
                "Connected to Cassandra (keyspace: " + KEYSPACE + ")"
            );

            // Create table if not exists
            try {
                String createTable = String.format(
                    "CREATE TABLE IF NOT EXISTS %s (" +
                        "recorded bigint, " +
                        "location text, " +
                        "sensor text, " +
                        "measurement text, " +
                        "units text, " +
                        "value double, " +
                        "PRIMARY KEY (recorded, location, sensor)" +
                        ");",
                    TABLE
                );
                session.execute(createTable);
                System.out.println("Table ensured: " + TABLE);
            } catch (Exception e) {
                System.err.println("Error creating table: " + e.getMessage());
                e.printStackTrace();
                return;
            }

            // Store sensor data
            try {
                JSONObject jsonObject =
                    SensorDataJsonHelper.getSensorDataJSONObject();
                System.out.println("JSON string marshalled to JSON Object");
                PreparedStatement insertStatement = session.prepare(
                    "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
                );
                BoundStatement boundStatement = insertStatement.bind(
                    (System.currentTimeMillis() / 1000),
                    jsonObject.getString("location"),
                    jsonObject.getString("sensor"),
                    jsonObject.getString("measurement"),
                    jsonObject.getString("units"),
                    jsonObject.getDouble("value")
                );
                session.execute(boundStatement);
                System.out.println("SensorData stored");
            } catch (Exception e) {
                System.err.println("Failed to save sensor data from JSON");
                e.printStackTrace();
            }

            // Retrieve sensor data
            List<String> jsonStrings = new ArrayList<>();
            try {
                ResultSet resultSet = session.execute(
                    "SELECT JSON * FROM sensor_data"
                );
                for (Row row : resultSet) {
                    jsonStrings.add(row.getString("[json]"));
                }
                System.out.println("SensorData retrieved");
            } catch (Exception e) {
                System.err.println("Failed to retrieve sensor data");
                e.printStackTrace();
            }

            // Print as CSV
            try {
                System.out.print(
                    SensorDataJsonHelper.jsonArrayToCsv(
                        jsonStrings.toArray(new String[0])
                    )
                );
            } catch (Exception e) {
                System.err.println("Failed to convert JSON to CSV");
                e.printStackTrace();
            }

            // Purge sensor data
            try {
                session.execute("TRUNCATE sensor_data");
                System.out.println("Sensor Data purged\n");
            } catch (Exception e) {
                System.err.println("Failed to purge sensor data");
                e.printStackTrace();
            }
        } catch (Exception e) {
            System.err.println(
                "Failed to connect to Cassandra (with keyspace) or perform operations."
            );
            e.printStackTrace();
        }
    }
}
