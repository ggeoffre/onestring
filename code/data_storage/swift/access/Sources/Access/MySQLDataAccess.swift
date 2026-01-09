// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

public class MySQLDataAccess: SensorDataAccess {
    public init() {}

    public func logSensorData(jsonData: String) throws {
        print("Logging sensor data to MySQL: \(jsonData)")
    }

    public func fetchSensorData() throws -> [String] {
        print("Fetching sensor data from MySQL")
        return ["{\"sensor\": \"light\", \"value\": 300}"]
    }

    public func purgeSensorData() throws {
        print("Purging sensor data from MySQL")
    }
}
