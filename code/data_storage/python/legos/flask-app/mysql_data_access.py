# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

import os

import pymysql
import pymysql.cursors
import sensor_data_helper
from sensor_data_access_protocol import SensorDataAccess

# Configurations
MYSQL_HOST = os.environ.get("DATA_HOSTNAME", "localhost").lower()
MYSQL_PORT = 3306
MYSQL_DB = "sensor_data_db"
MYSQL_TABLE = "sensor_data"
MYSQL_USER = "root"
MYSQL_PASS = ""

TABLE_CREATION_CHECKED = False


class MySQLDataAccess(SensorDataAccess):
    """MySQL implementation of SensorDataAccess."""

    connection = None

    def _ensure_schema(self, connection):
        global TABLE_CREATION_CHECKED
        if TABLE_CREATION_CHECKED:
            return
        with connection.cursor() as cursor:
            cursor.execute(f"CREATE DATABASE IF NOT EXISTS {MYSQL_DB};")
            connection.select_db(MYSQL_DB)
            cursor.execute(f"""
                CREATE TABLE IF NOT EXISTS {MYSQL_TABLE} (
                    recorded BIGINT NOT NULL,
                    location VARCHAR(255) NOT NULL,
                    sensor VARCHAR(255) NOT NULL,
                    measurement VARCHAR(255) NOT NULL,
                    units VARCHAR(10) NOT NULL,
                    value DECIMAL(5,2) NOT NULL
                );
            """)
        TABLE_CREATION_CHECKED = True

    def get_cursor(self) -> tuple:
        try:
            connection = pymysql.connect(
                host=MYSQL_HOST,
                port=MYSQL_PORT,
                user=MYSQL_USER,
                password=MYSQL_PASS,
                autocommit=True,
            )
            self._ensure_schema(connection)
            connection.select_db(MYSQL_DB)
            return connection, connection.cursor()
        except Exception as ex:
            print(f"Database error: {ex}")
            raise

    def log_sensor_data(self, json_data: str) -> None:
        try:
            connection, cursor = self.get_cursor()
            data_tuple = sensor_data_helper.create_insert_data_tuple(json_data)
            insert_sql = f"INSERT INTO {MYSQL_TABLE} VALUES (%s, %s, %s, %s, %s, %s)"
            cursor.execute(insert_sql, data_tuple)
            print("Data stored successfully")
        except Exception as ex:
            print(f"Database error: {ex}")
        return None

    def fetch_sensor_data(self) -> list[str]:
        try:
            connection, cursor = self.get_cursor()
            with connection.cursor(pymysql.cursors.DictCursor) as dict_cursor:
                dict_cursor.execute(f"SELECT * FROM {MYSQL_TABLE}")
                results = dict_cursor.fetchall()
                for row in results:
                    row["value"] = float(row["value"])
                if results:
                    print(f"Retrieved {len(results)} records")
                    print(results)
                    return results
        except Exception as ex:
            print(f"Database error: {ex}")
        return []

    def purge_sensor_data(self) -> None:
        try:
            connection, cursor = self.get_cursor()
            cursor.execute(f"TRUNCATE TABLE {MYSQL_TABLE}")
            print("Table truncated")
        except pymysql.MySQLError as e:
            print(f"Error purging data: {e}")
