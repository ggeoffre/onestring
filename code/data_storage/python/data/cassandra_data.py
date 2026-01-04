# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os

import sensor_data_helper
from cassandra import DriverException
from cassandra.cluster import Cluster

# Configurations
CASS_HOST = os.environ.get("DATA_HOSTNAME", "localhost").lower()
CASS_PORT = 9042
CASS_KEYSPACE = "sensor_data_db"
CASS_TABLE = "sensor_data"


def raw_cassandra_data():
    cluster = Cluster([CASS_HOST], port=CASS_PORT)

    try:
        # 1. Connect
        session = cluster.connect()
        print(f"Connected to Cassandra at {CASS_HOST}:{CASS_PORT}")

        # 1.1 Setup Keyspace and Table
        session.execute(f"""
            CREATE KEYSPACE IF NOT EXISTS {CASS_KEYSPACE}
            WITH REPLICATION = {{ 'class' : 'SimpleStrategy', 'replication_factor' : 1 }};
        """)
        session.set_keyspace(CASS_KEYSPACE)

        session.execute(f"""
            CREATE TABLE IF NOT EXISTS {CASS_TABLE} (
                recorded bigint,
                location text,
                sensor text,
                measurement text,
                units text,
                value double,
                PRIMARY KEY ((location), recorded, sensor)
            );
        """)

        # 2. Store
        data_dict = sensor_data_helper.generate_random_sensor_data()
        insert_query = f"""
            INSERT INTO {CASS_TABLE} (recorded, location, sensor, measurement, units, value)
            VALUES (%s, %s, %s, %s, %s, %s)
        """
        session.execute(
            insert_query, sensor_data_helper.create_insert_data_tuple(data_dict)
        )
        print("Data stored successfully")

        # 3. Retrieve
        rows = session.execute(f"SELECT * FROM {CASS_TABLE}")
        results = [row._asdict() for row in rows]
        if results:
            print(f"Retrieved {len(results)} records")
            print(results)

        # 4. Delete
        session.execute(f"TRUNCATE {CASS_TABLE}")
        print("Table truncated")

    except DriverException as e:
        print(f"Cassandra Error: {e}")
    finally:
        cluster.shutdown()
        print("Connection closed")
