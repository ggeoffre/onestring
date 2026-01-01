# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import os

import sensor_data_helper
from cassandra import DriverException
from cassandra.cluster import Cluster, Session
from sensor_data_access_protocol import SensorDataAccess

# Configurations
CASS_HOST = os.environ.get("DATA_HOSTNAME", "localhost").lower()
CASS_PORT = 9042
CASS_KEYSPACE = "sensor_data_db"
CASS_TABLE = "sensor_data"

TABLE_CREATION_CHECKED = False


class CassandraDataAccess(SensorDataAccess):
    """Cassandra implementation of SensorDataAccess."""

    cluster = Cluster([CASS_HOST], port=CASS_PORT)

    def get_connection(self) -> Session:
        global TABLE_CREATION_CHECKED
        try:
            session = self.cluster.connect()
            print(f"Connected to Cassandra at {CASS_HOST}:{CASS_PORT}")

            if not TABLE_CREATION_CHECKED:
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
                TABLE_CREATION_CHECKED = True
            else:
                session.set_keyspace(CASS_KEYSPACE)
            return session
        except Exception as e:
            print(f"Error creating table: {e}")
            return None

    def close_connection(self):
        try:
            self.cluster.shutdown()
        except Exception as e:
            print(f"Error closing connection: {e}")

    def log_sensor_data(self, json_data: str) -> None:
        try:
            session = self.get_connection()
            insert_query = f"""
                INSERT INTO {CASS_TABLE} (recorded, location, sensor, measurement, units, value)
                VALUES (%s, %s, %s, %s, %s, %s)
            """
            session.execute(
                insert_query, sensor_data_helper.create_insert_data_tuple(json_data)
            )
            print("Data stored successfully")
        except Exception as e:
            print(f"Error storing data: {e}")
        return None

    def fetch_sensor_data(self) -> list[str]:
        try:
            session = self.get_connection()
            rows = session.execute(f"SELECT * FROM {CASS_TABLE}")
            results = [row._asdict() for row in rows]
            if results:
                print(f"Retrieved {len(results)} records")
                print(results)
                return results
            else:
                print("No data found")
                return []
        except Exception as e:
            print(f"Error fetching data: {e}")
            return []

    def purge_sensor_data(self) -> None:
        try:
            session = self.get_connection()
            session.execute(f"TRUNCATE {CASS_TABLE}")
            print("Table truncated")
            return None
        except Exception as e:
            print(f"Error purging data: {e}")
            return None
