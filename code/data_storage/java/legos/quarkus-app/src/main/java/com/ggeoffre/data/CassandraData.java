// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.data;

import com.datastax.oss.driver.api.core.CqlSession;
import com.datastax.oss.driver.api.core.cql.*;
import java.net.InetSocketAddress;
import java.util.*;
import org.json.*;

public class CassandraData implements SensorDataAccess {

    // Configuration constants
    private static final String CONTACT_POINT = "192.168.1.60";
    private static final int PORT = 9042;
    private static final String LOCAL_DATACENTER = "datacenter1";
    private static final String KEYSPACE = "sensor_data_db";
    private static final String TABLE = "sensor_data";

    private static CqlSession session; // Singleton session

    // Initialize the session once
    public void init() {
        try {
            session = CqlSession.builder()
                .addContactPoint(new InetSocketAddress(CONTACT_POINT, PORT))
                .withLocalDatacenter(LOCAL_DATACENTER)
                .withKeyspace(KEYSPACE)
                .build();

            System.out.println(
                "Connected to Cassandra (keyspace: " + KEYSPACE + ")"
            );

            // Create keyspace and table if not exists
            String createKeyspace = String.format(
                "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };",
                KEYSPACE
            );
            session.execute(createKeyspace);

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

            System.out.println("Keyspace and table ensured.");
        } catch (Exception e) {
            System.err.println("Failed to initialize Cassandra session.");
            e.printStackTrace();
        }
    }

    // Reuse the singleton session
    public CqlSession getSession() {
        if (session == null || session.isClosed()) {
            session = CqlSession.builder()
                .addContactPoint(new InetSocketAddress(CONTACT_POINT, PORT))
                .withLocalDatacenter(LOCAL_DATACENTER)
                .withKeyspace(KEYSPACE)
                .build();
        }
        return session;
    }

    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to Cassandra: " + jsonData);
        try {
            JSONObject jsonObject =
                SensorDataJsonHelper.getSensorDataJSONObject(jsonData);
            PreparedStatement insertStatement = getSession().prepare(
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
            getSession().execute(boundStatement);
            System.out.println("SensorData stored");
        } catch (Exception e) {
            System.err.println("Failed to save sensor data from JSON");
            e.printStackTrace();
        }
    }

    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from Cassandra");
        List<String> jsonStrings = new ArrayList<>();
        try {
            ResultSet resultSet = getSession().execute(
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
        return jsonStrings;
    }

    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from Cassandra");
        try {
            getSession().execute("TRUNCATE sensor_data");
            System.out.println("Sensor Data purged\n");
        } catch (Exception e) {
            System.err.println("Failed to purge sensor data");
            e.printStackTrace();
        }
    }

    // Close the session when the application shuts down
    public void close() {
        if (session != null && !session.isClosed()) {
            session.close();
            System.out.println("Cassandra session closed.");
        }
    }
}
