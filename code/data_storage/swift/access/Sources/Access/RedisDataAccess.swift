// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

public class RedisDataAccess: SensorDataAccess {
    public init() {}

    public func logSensorData(jsonData: String) throws {
        print("Logging sensor data to Redis: \(jsonData)")
    }

    public func fetchSensorData() throws -> [String] {
        print("Fetching sensor data from Redis")
        return ["{\"sensor\": \"temperature\", \"value\": 22.3}"]
    }

    public func purgeSensorData() throws {
        print("Purging sensor data from Redis")
    }
}
