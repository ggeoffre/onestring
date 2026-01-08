// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package main

import (
	"fmt"
	"log"
	"strings"

	// Import the package where the other files reside.
	// If you named your module something else in 'go mod init', update this path.
	"go-ai-data/database"
)

func main() {
	fmt.Println("Starting Database Operations...")

	// 1. Cassandra
	printSeparator("Cassandra")
	if cass, err := database.NewCassandraData(); err != nil {
		log.Printf("Skipping Cassandra due to initialization error: %v\n", err)
	} else {
		// Ensure connection is closed when function exits (or explicitly here)
		// We use a func literal or explicit block to control scope if we want immediate cleanup
		func() {
			defer cass.Close()
			cass.RawCassandraData()
		}()
	}

	// 2. MongoDB
	printSeparator("MongoDB")
	if mongo, err := database.NewMongoData(); err != nil {
		log.Printf("Skipping MongoDB due to initialization error: %v\n", err)
	} else {
		func() {
			defer mongo.Close()
			mongo.RawMongoData()
		}()
	}

	// 3. Redis
	printSeparator("Redis")
	// Note: NewRedisData (in our refactor) does not return an error
	redis := database.NewRedisData()
	func() {
		defer redis.Close()
		redis.RawRedisData()
	}()

	// 4. Postgres
	printSeparator("Postgres")
	if pg, err := database.NewPostgresData(); err != nil {
		log.Printf("Skipping Postgres due to initialization error: %v\n", err)
	} else {
		func() {
			defer pg.Close()
			pg.RawPostgresData()
		}()
	}

	// 5. MySQL
	printSeparator("MySQL")
	if mysql, err := database.NewMySQLData(); err != nil {
		log.Printf("Skipping MySQL due to initialization error: %v\n", err)
	} else {
		func() {
			defer mysql.Close()
			mysql.RawMySQLData()
		}()
	}

	printSeparator("Done")
	fmt.Println("All database operations completed.")
}

func printSeparator(name string) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf(">>> EXECUTING %s <<<\n", strings.ToUpper(name))
	fmt.Println(strings.Repeat("-", 60))
}
