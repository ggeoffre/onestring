// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import NIOCore
import NIOPosix
@preconcurrency import RediStack

public class RedisDataAccess: SensorDataAccess {
    public init() {}

    // Connect to the Redis server
    private func connect() async throws -> RedisConnection {
        let group = MultiThreadedEventLoopGroup.singleton
        let future = RedisConnection.make(
            configuration: try RedisConnection.Configuration(hostname: "192.168.1.60", port: 6379),
            boundEventLoop: group.next()
        )
        return try await future.get()
    }

    public func logSensorData(jsonData: String) async throws {
        print("Logging sensor data to Redis: \(jsonData)")

        let connection = try await connect()

        // Parse the JSON string to extract sensor data fields
        guard let data = jsonData.data(using: .utf8),
            let json = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any]
        else {
            throw NSError(
                domain: "RedisDataAccess", code: 1,
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

        let sensorData = SensorData(
            recorded: recorded,
            location: location,
            sensor: sensor,
            measurement: measurement,
            units: units,
            value: value
        )

        let encoder = JSONEncoder()
        let jsonString = String(data: try encoder.encode(sensorData), encoding: .utf8)!
        _ = try await connection.rpush(jsonString, into: RedisKey("sensor_data")).get()
        print("Insert completed successfully")
    }

    public func fetchSensorData() async throws -> [String] {
        print("Fetching sensor data from Redis")

        let connection = try await connect()

        let respValues = try await connection.lrange(
            from: RedisKey("sensor_data"), firstIndex: 0, lastIndex: -1
        ).get()

        var results: [String] = []
        for value in respValues {
            if let rawString = value.string {
                results.append(rawString)
            }
        }

        return results
    }

    public func purgeSensorData() async throws {
        print("Purging sensor data from Redis")

        let connection = try await connect()
        _ = try await connection.delete(RedisKey("sensor_data")).get()
    }
}
