// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.data.access;

import java.util.List;

public interface SensorDataAccess {
    void logSensorData(String jsonData);
    List<String> fetchSensorData();
    void purgeSensorData();
}
