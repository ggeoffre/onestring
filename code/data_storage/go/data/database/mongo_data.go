// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorDataMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Recorded    int32              `bson:"recorded"`
	Location    string             `bson:"location"`
	Sensor      string             `bson:"sensor"`
	Measurement string             `bson:"measurement"`
	Units       string             `bson:"units"`
	Value       float64            `bson:"value"`
}

// MongoData struct for MongoDB
type MongoData struct {
	Client *mongo.Client
}

// NewMongoData creates a new MongoData instance and connects
func NewMongoData() (*MongoData, error) {
	// Connect to MongoDB
	clientOpts := options.Client().ApplyURI("mongodb://192.168.1.60:27017/sensor_data")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB")
	return &MongoData{Client: client}, nil
}

// Close disconnects the client
func (c *MongoData) Close() {
	if c.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.Client.Disconnect(ctx)
	}
}

func (c *MongoData) RawMongoData() {
	// Convert JSON string to SensorDataMongo struct
	var sensorDataMongo SensorDataMongo
	if err := json.Unmarshal([]byte(getRandomSensorDataJsonString()), &sensorDataMongo); err != nil {
		fmt.Printf("Failed to unmarshal JSON string into SensorData: %v\n", err)
		return
	}
	fmt.Println("JSON string converted to SensorDataMongo struct")

	// Insert the SensorData into the "sensor_data" collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := c.Client.Database("sensor_data_db").Collection("sensor_data")

	_, err := collection.InsertOne(ctx, sensorDataMongo)
	if err != nil {
		fmt.Printf("Failed to insert SensorData into MongoDB: %v\n", err)
		return
	}
	fmt.Println("SensorDataMongo stored")

	// Retrieve all SensorData from the "sensor_data" collection
	// Refresh context
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Failed to retrieve Collection from MongoDB: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and print each document
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		fmt.Printf("failed to retrieve SensorData from MongoDB: %v\n", err)
		return
	}

	for _, result := range results {
		delete(result, "_id")
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		fmt.Printf("Failed to marshal results: %v\n", err)
		return
	}

	// Convert JSON array to CSV string
	csvString, err := jsonToCSV(string(jsonResults))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n%s", "SensorDataMongo retrieved", csvString)
	}

	// Purge all SensorData from the "sensor_data" collection
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		fmt.Printf("failed to delete documents: %v\n", err)
		return
	}
	fmt.Println("SensorDataMongo purged")
}
