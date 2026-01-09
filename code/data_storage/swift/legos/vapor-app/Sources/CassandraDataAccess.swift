// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import CassandraClient
import Foundation

public class CassandraDataAccess: SensorDataAccess {

    public init() {
    }

    // Create a connection to the Cassandra cluster
    private func connect() async throws -> CassandraClient {
        let configuration = CassandraClient.Configuration(
            contactPointsProvider: { completion in
                completion(.success(["192.168.1.60"]))
            },
            port: 9042,
            protocolVersion: .v4
        )
        return CassandraClient(configuration: configuration)
    }

    public func logSensorData(jsonData: String) async throws {
        print("Logging sensor data to Cassandra: \(jsonData)")

        let client = try await connect()
        defer { try? client.shutdown() }

        // Parse the JSON string to extract sensor data fields
        guard let data = jsonData.data(using: .utf8),
            let json = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any]
        else {
            throw NSError(
                domain: "CassandraDataAccess", code: 1,
                userInfo: [NSLocalizedDescriptionKey: "Invalid JSON data"])
        }

        let recorded: Int64
        if let recordedInt = json["recorded"] as? Int {
            recorded = Int64(recordedInt)
        } else if let recordedDouble = json["recorded"] as? Double {
            recorded = Int64(recordedDouble)
        } else {
            recorded = Int64(Date().timeIntervalSince1970)
        }

        let location = json["location"] as? String ?? ""
        let sensor = json["sensor"] as? String ?? ""
        let measurement = json["measurement"] as? String ?? ""
        let units = json["units"] as? String ?? ""

        let value: Double
        if let valueDouble = json["value"] as? Double {
            value = valueDouble
        } else if let valueInt = json["value"] as? Int {
            value = Double(valueInt)
        } else {
            value = 0.0
        }

        let query =
            "INSERT INTO sensor_data_db.sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
        let params: [CassandraClient.Statement.Value] = [
            .int64(recorded),
            .string(location),
            .string(sensor),
            .string(measurement),
            .string(units),
            .double(value),
        ]

        try await client.run(query, parameters: params)
        print("Insert completed successfully")
    }

    public func fetchSensorData() async throws -> [String] {
        print("Fetching sensor data from Cassandra")

        let client = try await connect()
        defer { try? client.shutdown() }

        let query =
            "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data_db.sensor_data"
        let rows = try await client.query(query)

        return rows.compactMap { row -> String? in
            let recorded = row.column("recorded")?.int64 ?? 0
            let location = row.column("location")?.string ?? ""
            let sensor = row.column("sensor")?.string ?? ""
            let measurement = row.column("measurement")?.string ?? ""
            let units = row.column("units")?.string ?? ""
            let value = row.column("value")?.double ?? 0.0

            let jsonDict: [String: Any] = [
                "recorded": recorded,
                "location": location,
                "sensor": sensor,
                "measurement": measurement,
                "units": units,
                "value": value,
            ]

            if let jsonData = try? JSONSerialization.data(withJSONObject: jsonDict, options: []),
                let jsonString = String(data: jsonData, encoding: .utf8)
            {
                return jsonString
            }
            return nil
        }
    }

    public func purgeSensorData() async throws {
        print("Purging sensor data from Cassandra")

        let client = try await connect()
        defer { try? client.shutdown() }

        let query = "TRUNCATE sensor_data_db.sensor_data"
        _ = try await client.query(query)
    }
}
