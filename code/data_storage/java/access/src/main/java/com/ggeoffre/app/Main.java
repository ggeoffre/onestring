// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.app;

import com.geoffrey.data.access.*;
import java.util.Scanner;

public class Main {

    public static SensorDataAccess getDataAccess() {
        String dataAccessType = System.getenv("DATA_ACCESS");
        if (dataAccessType == null || dataAccessType.isEmpty()) {
            dataAccessType = "redis";
        }
        switch (dataAccessType.toLowerCase()) {
            case "redis":
                return new RedisDataAccess();
            case "mongo":
                return new MongoDataAccess();
            case "cassandra":
                return new CassandraDataAccess();
            case "mysql":
                return new MySQLDataAccess();
            case "postgres":
                return new PostgresDataAccess();
            default:
                throw new IllegalArgumentException(
                    "Unsupported DATA_ACCESS type: " + dataAccessType
                );
        }
    }

    public static void main(String[] args) {
        SensorDataAccess dataAccess = getDataAccess();

        dataAccess.logSensorData(
            "{\"sensor\": \"temperature\", \"value\": 22.3}"
        );
        System.out.println(
            "Fetched sensor data: " + dataAccess.fetchSensorData()
        );
        dataAccess.purgeSensorData();
    }
}
