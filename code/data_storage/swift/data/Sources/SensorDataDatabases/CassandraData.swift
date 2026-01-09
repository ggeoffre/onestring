// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import CassandraClient
import Foundation

// import SensorDataJsonHelper

class CassandraDatabase: SensorDataDatabase {
    private var client: CassandraClient?

    // Connect to the Cassandra cluster
    func connect() async throws {
        let configuration = CassandraClient.Configuration(
            contactPointsProvider: { completion in
                completion(.success(["192.168.1.60"]))
            },
            port: 9042,
            protocolVersion: .v4
        )
        client = CassandraClient(configuration: configuration)
    }

    // Insert sensor data into the Cassandra table
    func insert(sensorData: SensorData) async throws {
        let query =
            "INSERT INTO sensor_data_db.sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
        let params: [CassandraClient.Statement.Value] = [
            .int64(Int64(sensorData.recorded)),
            .string(sensorData.location),
            .string(sensorData.sensor),
            .string(sensorData.measurement),
            .string(sensorData.units),
            .double(sensorData.value),
        ]
        try await client?.run(query, parameters: params)
    }

    // Retrieve all sensor data from the Cassandra table
    func selectAll() async throws -> [SensorData] {
        let query =
            "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data_db.sensor_data"
        guard let rows = try await client?.query(query) else { return [] }
        return rows.map { row in
            SensorData(
                recorded: Int(row.column("recorded")?.int64 ?? 0),
                location: row.column("location")?.string ?? "",
                sensor: row.column("sensor")?.string ?? "",
                measurement: row.column("measurement")?.string ?? "",
                units: row.column("units")?.string ?? "",
                value: row.column("value")?.double ?? 0.0
            )
        }
    }

    // Delete all sensor data from the Cassandra table
    func deleteAll() async throws {
        let query = "TRUNCATE sensor_data_db.sensor_data"
        _ = try await client?.query(query)
    }

    // Close the connection to the Cassandra cluster
    func close() async throws {
        try? client?.shutdown()
        client = nil
    }
}
