// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation
import Vapor

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

    app.routes.get("report") { req -> Response in
        let csvString = try jsonStringsToCSV(jsonStrings: [Constants().jsonString])

        // Return CSV response
        return Response(
            status: .ok,
            headers: ["Content-Type": "text/csv"],
            body: .init(string: csvString)
        )
    }

    app.routes.get("purge") { req -> String in
        return "{\"message\": \"Purge executed\"}"
    }

    app.routes.post("purge") { req -> String in
        return "{\"message\": \"Purge executed\"}"
    }

}

struct Constants {
    let jsonString = """
        {"recorded":1737861909,"location":"den","sensor":"bmp280","measurement":"temperature","units":"C","value":33.3}
        """
}

func jsonStringsToCSV(jsonStrings: [String]) throws -> String {
    var csvRows = [String]()
    var header = [String]()

    for jsonStr in jsonStrings {
        let data = jsonStr.data(using: .utf8)!
        let json = try JSONSerialization.jsonObject(with: data, options: [])
        guard let jsonDict = json as? [String: Any] else {
            throw NSError(domain: "Invalid JSON", code: 0, userInfo: nil)
        }

        if header.isEmpty {
            header = Array(jsonDict.keys).sorted()  // or preserve insertion order if preferred
        }

        let row = header.map { key -> String in
            guard let value = jsonDict[key] else {
                return ""
            }

            switch value {
            case let str as String:
                // Quote only if string contains a comma, newline, or quote
                if str.contains(",") || str.contains("\"") || str.contains("\n") {
                    let escaped = str.replacingOccurrences(of: "\"", with: "\"\"")
                    return "\"\(escaped)\""
                } else {
                    return str
                }
            case let number as NSNumber:
                // NSNumber can represent Int, Float, Bool, etc.
                return number.stringValue
            case _ as NSNull:
                return ""
            default:
                // Fallback: try to stringify (e.g., for arrays or nested objects)
                return "\"\(value)\""
            }
        }.joined(separator: ",")

        csvRows.append(row)
    }

    let csv = [header.joined(separator: ","), csvRows.joined(separator: "\n")]
    return csv.joined(separator: "\n")
}
