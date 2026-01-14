// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import NIO
@preconcurrency import RediStack

// import SensorDataJsonHelper

class RedisDatabase: SensorDataDatabase {
    private var connection: RedisConnection?

    // Connect to the Redis server
    func connect() async throws {
        let group = MultiThreadedEventLoopGroup(numberOfThreads: 1)
        let future = RedisConnection.make(
            configuration: try RedisConnection.Configuration(hostname: "localhost", port: 6379),
            boundEventLoop: group.next()
        )
        connection = try await future.get()
    }

    // Insert sensor data as a JSON string into a Redis list
    func insert(sensorData: SensorData) async throws {
        let encoder = JSONEncoder()
        let jsonString = String(data: try encoder.encode(sensorData), encoding: .utf8)!
        _ = try await connection?.rpush(jsonString, into: RedisKey("sensor_data")).get()
    }

    // Retrieve all sensor data from the Redis list and decode from JSON strings
    func selectAll() async throws -> [SensorData] {
        let respValues =
            try await connection?.lrange(
                from: RedisKey("sensor_data"), firstIndex: 0, lastIndex: -1
            ).get() ?? []
        let decoder = JSONDecoder()
        var results: [SensorData] = []
        for value in respValues {
            if let rawString = value.string, let data = rawString.data(using: .utf8) {
                if let sensorData = try? decoder.decode(SensorData.self, from: data) {
                    results.append(sensorData)
                }
            }
        }
        return results
    }

    // Delete all sensor data from the Redis list
    func deleteAll() async throws {
        _ = try await connection?.delete(RedisKey("sensor_data")).get()
    }

    // Close the connection to the Redis server
    func close() async throws {
        connection = nil
    }
}
