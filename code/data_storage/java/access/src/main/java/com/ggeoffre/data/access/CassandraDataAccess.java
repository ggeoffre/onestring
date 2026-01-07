// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.data.access;

import java.util.Collections;
import java.util.List;

public class CassandraDataAccess implements SensorDataAccess {

    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to Cassandra: " + jsonData);
    }

    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from Cassandra");
        return Collections.singletonList(
            "{\"sensor\": \"pressure\", \"value\": 1013}"
        );
    }

    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from Cassandra");
    }
}
