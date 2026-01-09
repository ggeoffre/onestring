// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import MongoKitten

public class MongoDataAccess: SensorDataAccess {
    public init() {}

    // Connect to MongoDB server and access the database and collection
    private func connect() async throws -> MongoKitten.MongoCollection {
        let settings = try ConnectionSettings("mongodb://192.168.1.60:27017")
        let connection = try await MongoConnection.connect(settings: settings)
        let database = connection["sensor_data_db"]
        return database["sensor_data"]
    }

    public func logSensorData(jsonData: String) async throws {
        print("Logging sensor data to MongoDB: \(jsonData)")

        let collection = try await connect()

        // Parse the JSON string to extract sensor data fields
        guard let data = jsonData.data(using: .utf8),
            let json = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any]
        else {
            throw NSError(
                domain: "MongoDataAccess", code: 1,
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

        try await collection.insertEncoded(sensorData)
        print("Insert completed successfully")
    }

    public func fetchSensorData() async throws -> [String] {
        print("Fetching sensor data from MongoDB")

        let collection = try await connect()

        var results: [String] = []
        for try await doc in collection.find().decode(SensorData.self) {
            let jsonDict: [String: Any] = [
                "recorded": doc.recorded,
                "location": doc.location,
                "sensor": doc.sensor,
                "measurement": doc.measurement,
                "units": doc.units,
                "value": doc.value,
            ]

            if let jsonData = try? JSONSerialization.data(withJSONObject: jsonDict, options: []),
                let jsonString = String(data: jsonData, encoding: .utf8)
            {
                results.append(jsonString)
            }
        }

        return results
    }

    public func purgeSensorData() async throws {
        print("Purging sensor data from MongoDB")

        let collection = try await connect()
        try await collection.deleteAll(where: [:])
    }
}
