// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import Logging
import PostgresNIO

class PostgresDatabase: SensorDataDatabase {
    private var eventLoopGroup: EventLoopGroup?
    private var eventLoop: EventLoop?
    private var connection: PostgresConnection?

    // Establish a connection to the PostgreSQL database
    func connect() async throws {
        eventLoopGroup = MultiThreadedEventLoopGroup(numberOfThreads: 1)
        eventLoop = eventLoopGroup?.next()
        let logger = Logger(label: "PostgresNIO_Connection")
        let configuration = PostgresConnection.Configuration(
            host: "localhost",
            port: 5432,
            username: "postgres",
            password: "postgres",
            database: "sensor_data_db",
            tls: .disable
        )
        connection = try await awaitPostgresFuture(
            PostgresConnection.connect(
                on: eventLoop!,
                configuration: configuration,
                id: 1,
                logger: logger
            )
        )
    }

    // Insert a new sensor data record into the database
    func insert(sensorData: SensorData) async throws {
        let insertSQL = """
                INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value)
                VALUES ($1, $2, $3, $4, $5, $6)
            """
        let insertParams: [PostgresData] = [
            PostgresData(int64: Int64(sensorData.recorded)),
            PostgresData(string: sensorData.location),
            PostgresData(string: sensorData.sensor),
            PostgresData(string: sensorData.measurement),
            PostgresData(string: sensorData.units),
            PostgresData(double: sensorData.value),
        ]
        _ = try await awaitPostgresFuture(connection!.query(insertSQL, insertParams))
    }

    // Select all sensor data records from the database
    func selectAll() async throws -> [SensorData] {
        var results: [SensorData] = []
        let selectSQL = """
            SELECT recorded, location, sensor, measurement, units, value::DOUBLE PRECISION AS value
            FROM sensor_data
            """
        do {
            let queryResult = try await awaitPostgresFuture(connection!.query(selectSQL))
            for rawRow in queryResult.rows {
                let row = PostgresRandomAccessRow(rawRow)
                let recorded = try row["recorded"].decode(Int64?.self) ?? 0
                let location = try row["location"].decode(String?.self) ?? ""
                let sensor = try row["sensor"].decode(String?.self) ?? ""
                let measurement = try row["measurement"].decode(String?.self) ?? ""
                let units = try row["units"].decode(String?.self) ?? ""
                let value = try row["value"].decode(Double?.self) ?? 0.0
                let sensorData = SensorData(
                    recorded: Int(recorded),
                    location: location,
                    sensor: sensor,
                    measurement: measurement,
                    units: units,
                    value: value
                )
                results.append(sensorData)
            }
        } catch {
            print("Error decoding data: \(String(reflecting: error))")
            throw error
        }
        return results
    }

    // Delete all sensor data records from the database
    func deleteAll() async throws {
        let deleteSQL = "TRUNCATE TABLE sensor_data"
        _ = try await awaitPostgresFuture(connection!.query(deleteSQL))
    }

    // Close the database connection and clean up resources
    func close() async throws {
        if let connection = connection {
            try await awaitPostgresFuture(connection.close())
            self.connection = nil
        }
        if let eventLoopGroup = eventLoopGroup {
            try await eventLoopGroup.shutdownGracefully()
            self.eventLoopGroup = nil
            self.eventLoop = nil
        }
    }
}

// Helper for bridging EventLoopFuture to async/await
func awaitPostgresFuture<T>(_ future: EventLoopFuture<T>) async throws -> T {
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
