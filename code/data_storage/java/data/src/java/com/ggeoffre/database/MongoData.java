// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.geoffrey.database;

import com.geoffrey.database.SensorDataJsonHelper;
import com.mongodb.client.*;
import com.mongodb.client.model.Projections;
import java.util.*;
import org.bson.Document;

public class MongoData {

    // Configuration constants (consider externalizing for production)
    private static final String MONGO_URI = "mongodb://192.168.1.60:27017";
    private static final String DATABASE_NAME = "sensor_data_db";
    private static final String COLLECTION_NAME = "sensor_data";

    /**
     * Processes sensor data: inserts, retrieves, prints as CSV, and purges the collection.
     */
    public static void processSensorData() {
        System.out.println("MONGO");
        System.out.println("#####");

        // Try-with-resources ensures MongoClient is always closed
        try (MongoClient mongoClient = MongoClients.create(MONGO_URI)) {
            // Get the database and collection
            MongoDatabase db = mongoClient.getDatabase(DATABASE_NAME);
            MongoCollection<Document> collection = db.getCollection(
                COLLECTION_NAME
            );

            System.out.println("Connected to Mongo");

            // Convert JSON string to BSON Document
            Document sensorDoc = SensorDataJsonHelper.getSensorDataDocument();
            System.out.println("SensorData Document created from JSON");

            // Insert the document
            try {
                collection.insertOne(sensorDoc);
                System.out.println("SensorData stored");
            } catch (Exception e) {
                System.err.println("Failed to insert sensor data");
                e.printStackTrace();
            }

            // Retrieve all documents and convert to a JSON strings array
            List<String> jsonStrings = new ArrayList<>();
            try (
                MongoCursor<Document> cursor = collection
                    .find()
                    .projection(Projections.excludeId())
                    .iterator()
            ) {
                while (cursor.hasNext()) {
                    jsonStrings.add(cursor.next().toJson());
                }
                System.out.println("SensorData retrieved");
            } catch (Exception e) {
                System.err.println("Failed to retrieve sensor data");
                e.printStackTrace();
            }

            // Print as CSV
            try {
                System.out.print(
                    SensorDataJsonHelper.jsonArrayToCsv(
                        jsonStrings.toArray(new String[0])
                    )
                );
            } catch (Exception e) {
                System.err.println("Failed to convert JSON to CSV");
                e.printStackTrace();
            }

            // Delete all documents (use with caution in production!)
            try {
                collection.deleteMany(new Document());
                System.out.println("SensorData purged\n");
            } catch (Exception e) {
                System.err.println("Failed to purge sensor data");
                e.printStackTrace();
            }
        } catch (Exception e) {
            System.err.println("Failed to connect to MongoDB");
            e.printStackTrace();
        }
    }
}
