# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json
import os

import redis
import sensor_data_helper
from sensor_data_access_protocol import SensorDataAccess

REDIS_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
REDIS_PORT = 6379
REDIS_KEY_BASE = "location:den:list"


class RedisDataAccess(SensorDataAccess):
    """Redis implementation of SensorDataAccess."""

    def get_connection(self) -> redis.StrictRedis:
        try:
            redis_client = redis.StrictRedis(
                host=REDIS_HOST,
                port=REDIS_PORT,
                decode_responses=True,
                socket_timeout=5,  # Prevent hanging if Redis is down
            )
            redis_client.ping()
            print(f"Connected to Redis at {REDIS_HOST}:{REDIS_PORT}")
            return redis_client
        except redis.ConnectionError as e:
            print(f"Could not connect to Redis: {e}")
            raise

    def log_sensor_data(self, json_data: str) -> None:
        try:
            redis_client = self.get_connection()
            json_string = json.dumps(json_data, default=sensor_data_helper.json_default)
            redis_client.rpush(REDIS_KEY_BASE, json_string)
            print("Data stored successfully")
        except Exception as e:
            print(f"Error storing data: {e}")

    def fetch_sensor_data(self) -> list[str]:
        try:
            redis_client = self.get_connection()
            raw_results = redis_client.lrange(REDIS_KEY_BASE, 0, -1)
            if raw_results:
                results = [json.loads(item) for item in raw_results]
                print(f"Retrieved {len(results)} records")
                print(results)
                return results
            else:
                print("No data found in Redis")
                return []
        except Exception as e:
            print(f"Error retrieving data: {e}")
            return []

    def purge_sensor_data(self) -> None:
        try:
            redis_client = self.get_connection()
            redis_client.delete(REDIS_KEY_BASE)
            print("Key deleted")
        except Exception as e:
            print(f"Error purging data: {e}")
