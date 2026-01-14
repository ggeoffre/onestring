// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package com.ggeoffre.data;

import com.mongodb.client.*;
import com.mongodb.client.model.Projections;
import java.util.*;
import org.bson.Document;

public class MongoData implements SensorDataAccess {

    private static final String MONGO_URI = "mongodb://localhost:27017";
    private static final String DATABASE_NAME = "sensor_data_db";
    private static final String COLLECTION_NAME = "sensor_data";

    private final MongoClient mongoClient;
    private final MongoCollection<Document> collection;

    public MongoData() {
        this.mongoClient = MongoClients.create(MONGO_URI);
        MongoDatabase db = mongoClient.getDatabase(DATABASE_NAME);
        this.collection = db.getCollection(COLLECTION_NAME);
        System.out.println("MongoClient initialized and connected to MongoDB");
    }

    public MongoCollection<Document> getCollection() {
        return collection;
    }

    @Override
    public void logSensorData(String jsonData) {
        System.out.println("Logging sensor data to MongoDB: " + jsonData);

        Document sensorDoc = SensorDataJsonHelper.getSendorDataDocument(
            jsonData
        );
        System.out.println("SensorData Document created from JSON");

        try {
            collection.insertOne(sensorDoc);
            System.out.println("SensorData stored");
        } catch (Exception e) {
            System.err.println("Failed to insert sensor data");
            e.printStackTrace();
        }
    }

    @Override
    public List<String> fetchSensorData() {
        System.out.println("Fetching sensor data from MongoDB");

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
        return jsonStrings;
    }

    @Override
    public void purgeSensorData() {
        System.out.println("Purging sensor data from MongoDB");

        try {
            collection.deleteMany(new Document());
            System.out.println("SensorData purged");
        } catch (Exception e) {
            System.err.println("Failed to purge sensor data");
            e.printStackTrace();
        }
    }

    public void close() {
        mongoClient.close();
        System.out.println("MongoClient closed");
    }
}
