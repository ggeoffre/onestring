# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json
import os

import pymongo

# Configurations
MONGO_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
MONGO_PORT = 27017
MONGO_DB = "sensor_data_db"
MONGO_COLLECTION = "sensor_data"

# OneString
SENSOR_DATA_DICT = {
    "recorded": 1768570200,
    "location": "den",
    "sensor": "bmp280",
    "measurement": "temperature",
    "units": "C",
    "value": 22.3,
}


def main():
    # 1. Connect
    try:
        client = pymongo.MongoClient(
            f"mongodb://{MONGO_HOST}:{MONGO_PORT}/", serverSelectionTimeoutMS=5000
        )
        db = client[MONGO_DB]
        collection = db[MONGO_COLLECTION]
        client.admin.command("ping")  # Verify connection
        print("Connected to Mongo")
    except Exception as e:
        print(f"Connection failed: {e}")
        return

    # 2. Store
    try:
        new_sensor_data = SENSOR_DATA_DICT
        collection.insert_one(new_sensor_data)
        print("Record stored successfully")
    except Exception as e:
        print(f"Storage error: {e}")

    # 3. Retrieve
    try:
        cursor = collection.find({}, {"_id": 0})
        results = list(cursor)
        if results:
            print(f"Retrieved {len(results)} records")
            print(results)
        else:
            print("No matching records found")
    except Exception as e:
        print(f"Retrieval error: {e}")

    # 4. Delete
    try:
        collection.delete_many({})
        print("Removed all records")
    except Exception as e:
        print(f"Deletion error: {e}")
    finally:
        client.close()


if __name__ == "__main__":
    main()
