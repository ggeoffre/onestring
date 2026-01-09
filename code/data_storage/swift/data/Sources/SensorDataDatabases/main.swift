// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

import Foundation

let databases: [SensorDataDatabase] = [
    CassandraDatabase(),
    MongoDatabase(),
    MySQLDatabase(),
    PostgresDatabase(),
    RedisDatabase(),
]
let dbNames = ["Cassandra", "MongoDB", "MySQL", "PostgreSQL", "Redis"]

Task {
    for (index, db) in databases.enumerated() {
        print("\n--- Running workflow for \(dbNames[index]) ---")
        do {
            try await runDatabaseWorkflow(db: db)
        } catch {
            print("Error running workflow for \(dbNames[index]): \(error)")
        }
    }
    print("\nAll workflows complete.")
    exit(0)
}

// Keep the process alive until the async work is done
dispatchMain()
