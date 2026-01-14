// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import MySQLNIO

class MySQLDatabase: SensorDataDatabase {
    private var connection: MySQLConnection?

    // Establish a connection to the MySQL database
    func connect() async throws {
        let eventLoopGroup = MultiThreadedEventLoopGroup.singleton
        let eventLoop = eventLoopGroup.next()
        let address = try SocketAddress.makeAddressResolvingHost("localhost", port: 3306)
        let logger = Logger(label: "MySQLNIO_Connection")
        var tlsConfiguration = TLSConfiguration.makeClientConfiguration()
        tlsConfiguration.certificateVerification = .none
        let connectFuture = MySQLConnection.connect(
            to: address,
            username: "root",
            database: "sensor_data_db",
            password: "",
            tlsConfiguration: tlsConfiguration,
            serverHostname: nil,
            logger: logger,
            on: eventLoop
        )
        connection = try await awaitFuture(connectFuture)
    }

    // Insert a new sensor data record into the database
    func insert(sensorData: SensorData) async throws {
        let sql =
            "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
        let params: [MySQLData] = [
            MySQLData(int: sensorData.recorded),
            MySQLData(string: sensorData.location),
            MySQLData(string: sensorData.sensor),
            MySQLData(string: sensorData.measurement),
            MySQLData(string: sensorData.units),
            MySQLData(double: sensorData.value),
        ]
        _ = try await awaitFuture(connection!.query(sql, params))
    }

    // Select all sensor data records from the database
    func selectAll() async throws -> [SensorData] {
        let sql = "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data"
        let rows = try await awaitFuture(connection!.query(sql))
        return rows.map { row in
            SensorData(
                recorded: row.column("recorded")?.int ?? 0,
                location: row.column("location")?.string ?? "",
                sensor: row.column("sensor")?.string ?? "",
                measurement: row.column("measurement")?.string ?? "",
                units: row.column("units")?.string ?? "",
                value: row.column("value")?.double ?? 0.0
            )
        }
    }

    // Delete all sensor data records from the database
    func deleteAll() async throws {
        let sql = "DELETE FROM sensor_data"
        _ = try await awaitFuture(connection!.query(sql))
    }

    // Close the database connection and clean up resources
    func close() async throws {
        _ = try await awaitFuture(connection!.close())
        connection = nil
    }
}

// Helper for bridging EventLoopFuture to async/await
func awaitFuture<T>(_ future: EventLoopFuture<T>) async throws -> T {
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
