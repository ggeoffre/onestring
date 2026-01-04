# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os

from cassandra_data_access import CassandraDataAccess
from mongo_data_access import MongoDataAccess
from mysql_data_access import MySQLDataAccess
from postgres_data_access import PostgresDataAccess
from redis_data_access import RedisDataAccess
from sensor_data_access_protocol import SensorDataAccess


def get_data_access() -> SensorDataAccess:
    data_access_type = os.getenv("DATA_ACCESS", "mongo")

    if data_access_type == "redis":
        return RedisDataAccess()
    elif data_access_type == "mongo":
        return MongoDataAccess()
    elif data_access_type == "cassandra":
        return CassandraDataAccess()
    elif data_access_type == "mysql":
        return MySQLDataAccess()
    elif data_access_type == "postgres":
        return PostgresDataAccess()
    else:
        raise ValueError(f"Unsupported DATA_ACCESS type: {data_access_type}")


def main():
    data_access = get_data_access()

    # Log sensor data
    data_access.log_sensor_data('{"sensor": "temperature", "value": 22.3}')

    # Fetch sensor data
    data = data_access.fetch_sensor_data()
    print("Fetched sensor data:", data)

    # Purge sensor data
    data_access.purge_sensor_data()


if __name__ == "__main__":
    main()
