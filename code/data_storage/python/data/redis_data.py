# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json
import os

import redis
import sensor_data_helper

REDIS_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
REDIS_PORT = 6379
REDIS_KEY_BASE = "location:den:list"


def raw_redis_data():
    # 1. Connect
    try:
        redis_client = redis.StrictRedis(
            host=REDIS_HOST,
            port=REDIS_PORT,
            decode_responses=True,
            socket_timeout=5,  # Prevent hanging if Redis is down
        )
        redis_client.ping()
        print(f"Connected to Redis at {REDIS_HOST}:{REDIS_PORT}")
    except redis.ConnectionError as e:
        print(f"Could not connect to Redis: {e}")
        return

    # 2. Store
    try:
        data_dict = sensor_data_helper.generate_random_sensor_data()
        json_string = json.dumps(data_dict, default=sensor_data_helper.json_default)

        redis_client.rpush(REDIS_KEY_BASE, json_string)
        print("Data stored successfully")
    except Exception as e:
        print(f"Error storing data: {e}")

    # 3. Retrieve
    try:
        raw_results = redis_client.lrange(REDIS_KEY_BASE, 0, -1)

        if raw_results:
            results = [json.loads(item) for item in raw_results]
            print(f"Retrieved {len(results)} records")
            print(results)
        else:
            print("No data found in Redis")
    except Exception as e:
        print(f"Error retrieving data: {e}")

    # 4. Delete
    try:
        redis_client.delete(REDIS_KEY_BASE)
        print("Key deleted")
    except Exception as e:
        print(f"Error purging data: {e}")

    try:
        if redis_client.ping():
            redis_client.connection_pool.disconnect()
            print("Connection closed")
    except redis.ConnectionError:
        print("Connection was already closed or server is unreachable.")
    except Exception as e:
        print(f"Error closing connection: {e}")
