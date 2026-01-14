// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import Logging
import NIOCore
import NIOPosix
import PostgresNIO

public class PostgresDataAccess: SensorDataAccess {
    public init() {}

    // Helper for bridging EventLoopFuture to async/await
    private func awaitPostgresFuture<T>(_ future: EventLoopFuture<T>) async throws -> T {
        return try await withCheckedThrowingContinuation { continuation in
            future.whenComplete { result in
                switch result {
                case .success(let value):
                    continuation.resume(returning: value)
                case .failure(let error):
                    continuation.resume(throwing: error)
                }
            }
        }
    }

    // Establish a connection to the PostgreSQL database
    private func connect() async throws -> PostgresConnection {
        let eventLoopGroup = MultiThreadedEventLoopGroup.singleton
        let eventLoop = eventLoopGroup.next()
        let logger = Logger(label: "PostgresNIO_Connection")
        let configuration = PostgresConnection.Configuration(
            host: "localhost",
            port: 5432,
            username: "postgres",
            password: "postgres",
            database: "sensor_data_db",
            tls: .disable
        )
        return try await awaitPostgresFuture(
            PostgresConnection.connect(
                on: eventLoop,
                configuration: configuration,
                id: 1,
                logger: logger
            )
        )
    }

    public func logSensorData(jsonData: String) async throws {
        print("Logging sensor data to PostgreSQL: \(jsonData)")

        let connection = try await connect()

        // Parse the JSON string to extract sensor data fields
        guard let data = jsonData.data(using: .utf8),
            let json = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any]
        else {
            try await awaitPostgresFuture(connection.close())
            throw NSError(
                domain: "PostgresDataAccess", code: 1,
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

        let insertSQL = """
                INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
                VALUES ($1, $2, $3, $4, $5, $6)
            """
        let insertParams: [PostgresData] = [
            PostgresData(int64: recorded),
            PostgresData(string: location),
            PostgresData(string: sensor),
            PostgresData(string: measurement),
            PostgresData(string: units),
            PostgresData(double: value),
        ]

        do {
            _ = try await awaitPostgresFuture(connection.query(insertSQL, insertParams))
            print("Insert completed successfully")
        } catch {
            try await awaitPostgresFuture(connection.close())
            throw error
        }

        try await awaitPostgresFuture(connection.close())
    }

    public func fetchSensorData() async throws -> [String] {
        print("Fetching sensor data from PostgreSQL")

        let connection = try await connect()

        let selectSQL = """
                SELECT recorded, location, sensor, measurement, units, value::DOUBLE PRECISION AS value
                FROM sensor_data
            """

        var results: [String] = []

        do {
            let queryResult = try await awaitPostgresFuture(connection.query(selectSQL))
            for rawRow in queryResult.rows {
                let row = PostgresRandomAccessRow(rawRow)
                let recorded = try row["recorded"].decode(Int64?.self) ?? 0
                let location = try row["location"].decode(String?.self) ?? ""
                let sensor = try row["sensor"].decode(String?.self) ?? ""
                let measurement = try row["measurement"].decode(String?.self) ?? ""
                let units = try row["units"].decode(String?.self) ?? ""
                let value = try row["value"].decode(Double?.self) ?? 0.0

                let jsonDict: [String: Any] = [
                    "recorded": recorded,
                    "location": location,
                    "sensor": sensor,
                    "measurement": measurement,
                    "units": units,
                    "value": value,
                ]

                if let jsonData = try? JSONSerialization.data(
                    withJSONObject: jsonDict, options: []),
                    let jsonString = String(data: jsonData, encoding: .utf8)
                {
                    results.append(jsonString)
                }
            }
        } catch {
            try await awaitPostgresFuture(connection.close())
            throw error
        }

        try await awaitPostgresFuture(connection.close())
        return results
    }

    public func purgeSensorData() async throws {
        print("Purging sensor data from PostgreSQL")

        let connection = try await connect()

        let deleteSQL = "TRUNCATE TABLE sensor_data"

        do {
            _ = try await awaitPostgresFuture(connection.query(deleteSQL))
        } catch {
            try await awaitPostgresFuture(connection.close())
            throw error
        }

        try await awaitPostgresFuture(connection.close())
    }
}
