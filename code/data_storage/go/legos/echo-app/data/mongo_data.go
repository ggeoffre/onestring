// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SensorDataJSON is for unmarshaling JSON where recorded might be a string
type SensorDataMongoJSON struct {
	Recorded    interface{} `json:"recorded"` // Accept both string and number
	Location    string      `json:"location"`
	Sensor      string      `json:"sensor"`
	Measurement string      `json:"measurement"`
	Units       string      `json:"units"`
	Value       float64     `json:"value"`
}

type SensorDataMongo struct {
	Recorded    int64   `json:"recorded" bson:"recorded"`
	Location    string  `json:"location" bson:"location"`
	Sensor      string  `json:"sensor" bson:"sensor"`
	Measurement string  `json:"measurement" bson:"measurement"`
	Units       string  `json:"units" bson:"units"`
	Value       float64 `json:"value" bson:"value"`
}

// MongoDataAccess implements the SensorDataAccess interface for MongoDB
type MongoDataAccess struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// NewMongoDataAccess creates a new MongoDataAccess instance and connects
func NewMongoDataAccess() (*MongoDataAccess, error) {
	// Suppress mongo driver logging
	log.SetOutput(io.Discard)
	defer log.SetOutput(log.Writer())

	// Connect to MongoDB
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
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

	// Get collection reference
	collection := client.Database("sensor_data_db").Collection("sensor_data")

	return &MongoDataAccess{
		Client:     client,
		Collection: collection,
	}, nil
}

// Close disconnects the client
func (m *MongoDataAccess) Close() {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		m.Client.Disconnect(ctx)
	}
}

// parseRecorded converts the recorded field to int64 whether it's a string or number
func parseRecordedMongo(recorded interface{}) (int64, error) {
	switch v := recorded.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("recorded field has unexpected type: %T", v)
	}
}

func (m *MongoDataAccess) LogSensorData(jsonData string) error {
	var jsonStruct SensorDataMongoJSON
	if err := json.Unmarshal([]byte(jsonData), &jsonStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}

	recorded, err := parseRecordedMongo(jsonStruct.Recorded)
	if err != nil {
		return fmt.Errorf("failed to parse recorded field: %w", err)
	}

	sensorData := SensorDataMongo{
		Recorded:    recorded,
		Location:    jsonStruct.Location,
		Sensor:      jsonStruct.Sensor,
		Measurement: jsonStruct.Measurement,
		Units:       jsonStruct.Units,
		Value:       jsonStruct.Value,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = m.Collection.InsertOne(ctx, sensorData)
	if err != nil {
		return fmt.Errorf("failed to insert into MongoDB: %w", err)
	}

	return nil
}

func (m *MongoDataAccess) FetchSensorData() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var results []string
	for cursor.Next(ctx) {
		var sensorData SensorDataMongo
		if err := cursor.Decode(&sensorData); err != nil {
			fmt.Printf("failed to decode document: %v\n", err)
			continue
		}

		jsonBytes, err := json.Marshal(sensorData)
		if err != nil {
			fmt.Printf("failed to marshal SensorDataMongo to JSON: %v\n", err)
			continue
		}
		results = append(results, string(jsonBytes))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}

func (m *MongoDataAccess) PurgeSensorData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.Collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete documents from MongoDB: %w", err)
	}

	return nil
}
