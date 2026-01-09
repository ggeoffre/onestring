// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

public class MongoDataAccess: SensorDataAccess {
    public init() {}

    public func logSensorData(jsonData: String) throws {
        print("Logging sensor data to MongoDB: \(jsonData)")
    }

    public func fetchSensorData() throws -> [String] {
        print("Fetching sensor data from MongoDB")
        return ["{\"sensor\": \"humidity\", \"value\": 45.6}"]
    }

    public func purgeSensorData() throws {
        print("Purging sensor data from MongoDB")
    }
}
