# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os

import psycopg2
import sensor_data_helper
from psycopg2.extras import RealDictCursor

POSTGRES_HOST = os.environ.get("DATA_HOSTNAME", "localhost").lower()
POSTGRES_PORT = 5432
POSTGRES_DB = "sensor_data_db"
POSTGRES_TABLE = "sensor_data"
POSTGRES_USER = "postgres"
POSTGRES_PASS = ""


def raw_postgres_data():
    connection = None
    try:
        # 1.0 Administrative Connection (to 'postgres' default DB)
        connection = psycopg2.connect(
            host=POSTGRES_HOST,
            port=POSTGRES_PORT,
            database="postgres",
            user=POSTGRES_USER,
            password=POSTGRES_PASS,
        )
        connection.autocommit = True
        cursor = connection.cursor()

        # 1.1 Setup Database
        cursor.execute(
            "SELECT 1 FROM pg_catalog.pg_database WHERE datname = %s", (POSTGRES_DB,)
        )
        if not cursor.fetchone():
            cursor.execute(f"CREATE DATABASE {POSTGRES_DB}")
            print(f"Database {POSTGRES_DB} created")

        # 1. Connect
        connection.close()  # Administrative connection must close to switch DBs
        connection = psycopg2.connect(
            host=POSTGRES_HOST,
            port=POSTGRES_PORT,
            database=POSTGRES_DB,
            user=POSTGRES_USER,
            password=POSTGRES_PASS,
        )
        connection.autocommit = True
        cursor = connection.cursor()
        print(f"Connected to PostgreSQL at {POSTGRES_HOST}:{POSTGRES_PORT}")

        # 1.3 Create Table
        cursor.execute(f"""
            CREATE TABLE IF NOT EXISTS {POSTGRES_TABLE} (
                recorded BIGINT NOT NULL,
                location VARCHAR(255) NOT NULL,
                sensor VARCHAR(255) NOT NULL,
                measurement VARCHAR(255) NOT NULL,
                units VARCHAR(255) NOT NULL,
                value NUMERIC(10, 2) NOT NULL
            );
        """)

        # 2. Store
        data_tuple = sensor_data_helper.create_insert_data_tuple(
            sensor_data_helper.generate_random_sensor_data()
        )
        cursor.execute(
            f"""
            INSERT INTO {POSTGRES_TABLE} (recorded, location, sensor, measurement, units, value)
            VALUES (%s, %s, %s, %s, %s, %s)
        """,
            data_tuple,
        )
        print("Data stored successfully")

        # 3. Retrieve
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

        # 4. Delete
        cursor.execute(f"TRUNCATE TABLE {POSTGRES_TABLE}")
        print("Table truncated")

    except psycopg2.Error as e:
        print(f"PostgreSQL Error: {e}")
    finally:
        if connection:
            connection.close()
            print("Connection closed")
