// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import MongoKitten

class MongoDatabase: SensorDataDatabase {
    private var database: MongoKitten.MongoDatabase?
    private var collection: MongoKitten.MongoCollection?

    // Connect to MongoDB server and access the database and collection
    func connect() async throws {
        let settings = try ConnectionSettings("mongodb://localhost:27017")
        let connection = try await MongoConnection.connect(settings: settings)
        database = connection["sensor_data_db"]
        collection = database?["sensor_data"]
    }

    // Insert sensor data into the MongoDB collection
    func insert(sensorData: SensorData) async throws {
        guard let collection = collection else { return }
        try await collection.insertEncoded(sensorData)
    }

    // Select all sensor data from the MongoDB collection
    func selectAll() async throws -> [SensorData] {
        guard let collection = collection else { return [] }
        var results: [SensorData] = []
        for try await doc in collection.find().decode(SensorData.self) {
            results.append(doc)
        }
        return results
    }

    // Delete all sensor data from the MongoDB collection
    func deleteAll() async throws {
        guard let collection = collection else { return }
        try await collection.deleteAll(where: [:])
    }

    // Close the connection to the MongoDB server
    func close() async throws {
        database = nil
        collection = nil
    }
}
