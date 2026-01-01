# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import os

import pymongo
import sensor_data_helper
from pymongo.typings import ClusterTime
from sensor_data_access_protocol import SensorDataAccess

# Configurations
MONGO_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
MONGO_PORT = 27017
MONGO_DB = "sensor_data_db"
MONGO_COLLECTION = "sensor_data"


class MongoDataAccess(SensorDataAccess):
    """MongoDB implementation of SensorDataAccess."""

    def get_connection(self) -> tuple:
        try:
            client = pymongo.MongoClient(
                f"mongodb://{MONGO_HOST}:{MONGO_PORT}/", serverSelectionTimeoutMS=5000
            )
            db = client[MONGO_DB]
            collection = db[MONGO_COLLECTION]
            client.admin.command("ping")  # Verify connection
            print(f"Connected to Mongo at {MONGO_HOST}:{MONGO_PORT}")
            return client, db, collection
        except Exception as e:
            print(f"Connection failed: {e}")
            return None, None, None

    def close_connection(self, client, db, collection):
        try:
            client.close()
            print("MongoDB connection closed")
        except Exception as e:
            print(f"Failed to close MongoDB connection: {e}")

    def log_sensor_data(self, json_data: str) -> None:
        try:
            client, db, collection = self.get_connection()
            collection.insert_one(json_data)
            print("Record stored successfully")
        except Exception as e:
            print(f"Storage error: {e}")
        return None

    def fetch_sensor_data(self) -> list[str]:
        try:
            client, db, collection = self.get_connection()
            cursor = collection.find({}, {"_id": 0})
            results = list(cursor)
            if results:
                print(f"Retrieved {len(results)} records")
            else:
                print("No matching records found")
            return results
        except Exception as e:
            print(f"Fetch error: {e}")
            return []

    def purge_sensor_data(self) -> None:
        try:
            client, db, collection = self.get_connection()
            collection.delete_many({})
            print("Sensor data purged from MongoDB")
        except Exception as e:
            print(f"Purge error: {e}")
        return None
