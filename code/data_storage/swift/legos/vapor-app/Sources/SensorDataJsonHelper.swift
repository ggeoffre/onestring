// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

// Define the structure for sensor data
public struct SensorData: Codable {
    var recorded: Int
    var location: String
    var sensor: String
    var measurement: String
    var units: String
    var value: Double
}

// Constants for testing
struct Constants {
    let jsonString = """
        {"recorded":1737861909,"location":"den","sensor":"bmp280","measurement":"temperature","units":"C","value":33.3}
        """
}

// Generate random sensor data for testing
public func get_random_sensor_data() -> SensorData {
    SensorData(
        recorded: Int(Date().timeIntervalSince1970),
        location: "den",
        sensor: "bmp280",
        measurement: "temperature",
        units: "C",
        value: Double(round(10 * Double.random(in: 22.4...32.1)) / 10)
    )
}

// Convert an array of JSON strings to CSV format
public func jsonStringsToCSV(jsonStrings: [String]) throws -> String {
    var csvRows = [String]()
    var header = [String]()

    for jsonStr in jsonStrings {
        let data = jsonStr.data(using: .utf8)!
        let json = try JSONSerialization.jsonObject(with: data, options: [])
        guard let jsonDict = json as? [String: Any] else {
            throw NSError(domain: "Invalid JSON", code: 0, userInfo: nil)
        }

        if header.isEmpty {
            header = Array(jsonDict.keys).sorted()
        }

        let row = header.map { key -> String in
            guard let value = jsonDict[key] else { return "" }
            switch value {
            case let str as String:
                if str.contains(",") || str.contains("\"") || str.contains("\n") {
                    let escaped = str.replacingOccurrences(of: "\"", with: "\"\"")
                    return "\"\(escaped)\""
                } else {
                    return str
                }
            case let number as NSNumber:
                return number.stringValue
            case _ as NSNull:
                return ""
            default:
                return "\"\(value)\""
            }
        }.joined(separator: ",")

        csvRows.append(row)
    }

    let csv = [header.joined(separator: ","), csvRows.joined(separator: "\n")]
    return csv.joined(separator: "\n")
}
