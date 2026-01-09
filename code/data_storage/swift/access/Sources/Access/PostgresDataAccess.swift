// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

public class PostgresDataAccess: SensorDataAccess {
    public init() {}

    public func logSensorData(jsonData: String) throws {
        print("Logging sensor data to PostgreSQL: \(jsonData)")
    }

    public func fetchSensorData() throws -> [String] {
        print("Fetching sensor data from PostgreSQL")
        return ["{\"sensor\": \"sound\", \"value\": 75}"]
    }

    public func purgeSensorData() throws {
        print("Purging sensor data from PostgreSQL")
    }
}
