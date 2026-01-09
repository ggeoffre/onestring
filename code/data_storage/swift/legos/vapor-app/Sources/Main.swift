// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import Vapor

func getDataAccess() -> SensorDataAccess {
    let dataAccessType = ProcessInfo.processInfo.environment["DATA_ACCESS"] ?? "redis"

    switch dataAccessType {
    case "redis":
        return RedisDataAccess()
    case "mongo":
        return MongoDataAccess()
    case "cassandra":
        return CassandraDataAccess()
    case "mysql":
        return MySQLDataAccess()
    case "postgres":
        return PostgresDataAccess()
    default:
        fatalError("Unsupported DATA_ACCESS type: \(dataAccessType)")
    }
}

@main
struct App {
    static func main() async throws {
        var env = try Environment.detect()
        try LoggingSystem.bootstrap(from: &env)
        let app = try await Application.make(env, .singleton)
        try configure(app)
        try await app.execute()
    }
}

func configure(_ app: Application) throws {

    app.http.server.configuration.hostname = "0.0.0.0"
    app.http.server.configuration.port = 8080

    app.routes.get { req -> String in
        return "{\"message\": \"Vapor API Server is running!\"}"
    }

    app.routes.post("echo") { req async throws -> Response in
        // Get the raw request body as ByteBuffer
        guard let body = req.body.data else {
            throw Abort(.badRequest, reason: "Invalid JSON data")
        }

        // Echo it back as JSON response
        return Response(
            status: .ok,
            headers: ["Content-Type": "application/json"],
            body: .init(buffer: body)
        )
    }

    app.routes.post("log") { req async throws -> Response in
        // Get the raw request body as ByteBuffer
        guard let body = req.body.data,
            let jsonString = body.getString(at: 0, length: body.readableBytes)
        else {
            throw Abort(.badRequest, reason: "Invalid JSON data")
        }

        // Create an instance of SensorDataAccess
        let dataAccess = getDataAccess()

        // Log the sensor data
        do {
            try await dataAccess.logSensorData(jsonData: jsonString)
        } catch {
            throw Abort(
                .internalServerError,
                reason: "Failed to log sensor data: \(error.localizedDescription)")
        }

        // Return a success response
        return Response(
            status: .ok,
            headers: ["Content-Type": "application/json"],
            body: .init(string: "{\"message\": \"Sensor data logged successfully\"}")
        )
    }

    app.routes.get("report") { req async throws -> Response in
        // Create an instance of SensorDataAccess
        let dataAccess = getDataAccess()

        // Fetch sensor data
        let sensorData: [String]
        do {
            sensorData = try await dataAccess.fetchSensorData()
        } catch {
            throw Abort(
                .internalServerError,
                reason: "Failed to fetch sensor data: \(error.localizedDescription)")
        }

        // Convert the sensor data to CSV
        let csvString: String
        do {
            csvString = try jsonStringsToCSV(jsonStrings: sensorData)
        } catch {
            throw Abort(
                .internalServerError,
                reason: "Failed to convert sensor data to CSV: \(error.localizedDescription)")
        }

        // Return CSV response
        return Response(
            status: .ok,
            headers: ["Content-Type": "text/csv"],
            body: .init(string: csvString)
        )
    }

    app.routes.get("purge") { req async throws -> String in
        // Create an instance of SensorDataAccess
        let dataAccess = getDataAccess()

        // Purge sensor data
        do {
            try await dataAccess.purgeSensorData()
        } catch {
            throw Abort(
                .internalServerError,
                reason: "Failed to purge sensor data: \(error.localizedDescription)")
        }

        return "{\"message\": \"Purge executed\"}"
    }

    app.routes.post("purge") { req async throws -> String in
        // Create an instance of SensorDataAccess
        let dataAccess = getDataAccess()

        // Purge sensor data
        do {
            try await dataAccess.purgeSensorData()
        } catch {
            throw Abort(
                .internalServerError,
                reason: "Failed to purge sensor data: \(error.localizedDescription)")
        }

        return "{\"message\": \"Purge executed\"}"
    }

}
