// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.data.access;

import java.util.Collections;
import java.util.List;

public class PostgresDataAccess implements SensorDataAccess {

    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to PostgreSQL: " + jsonData);
    }

    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from PostgreSQL");
        return Collections.singletonList(
            "{\"sensor\": \"sound\", \"value\": 75}"
        );
    }

    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from PostgreSQL");
    }
}
