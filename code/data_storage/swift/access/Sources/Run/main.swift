// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Access
import Foundation

func getDataAccess() -> SensorDataAccess {
    let dataAccessType = ProcessInfo.processInfo.environment["DATA_ACCESS"] ?? "redis"

    switch dataAccessType {
    case "redis":
        return RedisDataAccess()
    case "mongo":
        return MongoDataAccess()
    case "cassandra":
        return CassandraDataAccess()
    case "mysql":
        return MySQLDataAccess()
    case "postgres":
        return PostgresDataAccess()
    default:
        fatalError("Unsupported DATA_ACCESS type: \(dataAccessType)")
    }
}

func main() {
    let dataAccess = getDataAccess()

    do {
        // Log sensor data
        try dataAccess.logSensorData(jsonData: "{\"sensor\": \"temperature\", \"value\": 22.3}")

        // Fetch sensor data
        let data = try dataAccess.fetchSensorData()
        print("Fetched sensor data: \(data)")

        // Purge sensor data
        try dataAccess.purgeSensorData()
    } catch {
        print("An error occurred: \(error)")
    }
}

main()
