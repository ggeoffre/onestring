# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import os

import psycopg2
import sensor_data_helper
from psycopg2.extras import RealDictCursor
from sensor_data_access_protocol import SensorDataAccess

POSTGRES_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
POSTGRES_PORT = 5432
POSTGRES_DB = "sensor_data_db"
POSTGRES_TABLE = "sensor_data"
POSTGRES_USER = "postgres"
POSTGRES_PASS = ""

TABLE_CREATION_CHECKED = False


class PostgresDataAccess(SensorDataAccess):
    """PostgreSQL implementation of SensorDataAccess."""

    connection = None

    def _ensure_schema(self):
        global TABLE_CREATION_CHECKED
        if TABLE_CREATION_CHECKED:
            return

        try:
            # 1. Connect to default 'postgres' db to check/create the target DB
            admin_conn = psycopg2.connect(
                host=POSTGRES_HOST,
                port=POSTGRES_PORT,
                database="postgres",
                user=POSTGRES_USER,
                password=POSTGRES_PASS,
            )
            admin_conn.autocommit = True
            with admin_conn.cursor() as cur:
                cur.execute(
                    "SELECT 1 FROM pg_database WHERE datname = %s", (POSTGRES_DB,)
                )
                if not cur.fetchone():
                    cur.execute(f"CREATE DATABASE {POSTGRES_DB}")
            admin_conn.close()

            # 2. Connect to the NEWLY CREATED/EXISTING target DB to create the table
            target_conn = psycopg2.connect(
                host=POSTGRES_HOST,
                port=POSTGRES_PORT,
                database=POSTGRES_DB,
                user=POSTGRES_USER,
                password=POSTGRES_PASS,
            )
            target_conn.autocommit = True
            with target_conn.cursor() as cur:
                cur.execute(f"""
                    CREATE TABLE IF NOT EXISTS {POSTGRES_TABLE} (
                        recorded BIGINT NOT NULL,
                        location VARCHAR(255) NOT NULL,
                        sensor VARCHAR(255) NOT NULL,
                        measurement VARCHAR(255) NOT NULL,
                        units VARCHAR(255) NOT NULL,
                        value NUMERIC(10, 2) NOT NULL
                    );
                """)
            target_conn.close()

            TABLE_CREATION_CHECKED = True
        except Exception as e:
            print(f"Schema Setup Error: {e}")

    def get_cursor(self) -> tuple:
        # Call schema check BEFORE attempting regular connection
        self._ensure_schema()

        try:
            connection = psycopg2.connect(
                host=POSTGRES_HOST,
                port=POSTGRES_PORT,
                database=POSTGRES_DB,
                user=POSTGRES_USER,
                password=POSTGRES_PASS,
            )
            connection.autocommit = True
            return connection, connection.cursor()
        except Exception as e:
            print(f"Error getting cursor: {e}")
            raise

    def log_sensor_data(self, json_data: str) -> None:
        try:
            connection, cursor = self.get_cursor()
            data_tuple = sensor_data_helper.create_insert_data_tuple(json_data)
            cursor.execute(
                f"""
                INSERT INTO {POSTGRES_TABLE} (recorded, location, sensor, measurement, units, value)
                VALUES (%s, %s, %s, %s, %s, %s)
            """,
                data_tuple,
            )
            print("Data stored successfully")
        except Exception as e:
            print(f"Error storing data: {e}")
        return None

    def fetch_sensor_data(self) -> list[str]:
        try:
            connection, cursor = self.get_cursor()
            with connection.cursor(cursor_factory=RealDictCursor) as dict_cursor:
                dict_cursor.execute(f"SELECT * FROM {POSTGRES_TABLE}")
                raw_results = dict_cursor.fetchall()
                results = []
                for row in raw_results:
                    clean_row = dict(row)
                    clean_row["value"] = float(clean_row["value"])
                    results.append(clean_row)
                if results:
                    print(f"Retrieved {len(results)} records")
                    print(results)
                    return results
                else:
                    print("No records found")
        except Exception as e:
            print(f"Error fetching data: {e}")
        return []

    def purge_sensor_data(self) -> None:
        try:
            connection, cursor = self.get_cursor()
            cursor.execute(f"TRUNCATE TABLE {POSTGRES_TABLE}")
            print("Table truncated")
        except Exception as e:
            print(f"Error purging data: {e}")
        return None
