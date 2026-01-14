// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import Logging
import MySQLNIO
import NIOCore
import NIOPosix
import NIOSSL

public class MySQLDataAccess: SensorDataAccess {
    public init() {}

    // Establish a connection to the MySQL database
    private func connect() async throws -> MySQLConnection {
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
        return try await connectFuture.get()
    }

    public func logSensorData(jsonData: String) async throws {
        print("Logging sensor data to MySQL: \(jsonData)")

        let connection = try await connect()

        // Parse the JSON string to extract sensor data fields
        guard let data = jsonData.data(using: .utf8),
            let json = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any]
        else {
            try await connection.close().get()
            throw NSError(
                domain: "MySQLDataAccess", code: 1,
                userInfo: [NSLocalizedDescriptionKey: "Invalid JSON data"])
        }

        let recorded: Int
        if let recordedInt = json["recorded"] as? Int {
            recorded = recordedInt
        } else if let recordedDouble = json["recorded"] as? Double {
            recorded = Int(recordedDouble)
        } else {
            recorded = Int(Date().timeIntervalSince1970)
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

        let sql =
            "INSERT INTO sensor_data (recorded, location, sensor, measurement, units, value) VALUES (?, ?, ?, ?, ?, ?)"
        let params: [MySQLData] = [
            MySQLData(int: recorded),
            MySQLData(string: location),
            MySQLData(string: sensor),
            MySQLData(string: measurement),
            MySQLData(string: units),
            MySQLData(double: value),
        ]

        do {
            _ = try await connection.query(sql, params).get()
            print("Insert completed successfully")
        } catch {
            try await connection.close().get()
            throw error
        }

        try await connection.close().get()
    }

    public func fetchSensorData() async throws -> [String] {
        print("Fetching sensor data from MySQL")

        let connection = try await connect()

        let sql = "SELECT recorded, location, sensor, measurement, units, value FROM sensor_data"
        let rows: [MySQLRow]
        do {
            rows = try await connection.query(sql).get()
        } catch {
            try await connection.close().get()
            throw error
        }

        let result = rows.compactMap { row -> String? in
            let recorded = row.column("recorded")?.int ?? 0
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

        try await connection.close().get()
        return result
    }

    public func purgeSensorData() async throws {
        print("Purging sensor data from MySQL")

        let connection = try await connect()

        let sql = "DELETE FROM sensor_data"
        do {
            _ = try await connection.query(sql).get()
        } catch {
            try await connection.close().get()
            throw error
        }

        try await connection.close().get()
    }
}
