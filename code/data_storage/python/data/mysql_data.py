# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import os

import pymysql
import pymysql.cursors
import sensor_data_helper

# Configurations
MYSQL_HOST = os.environ.get("DATA_HOSTNAME", "192.168.1.60").lower()
MYSQL_PORT = 3306
MYSQL_DB = "sensor_data_db"
MYSQL_TABLE = "sensor_data"
MYSQL_USER = "root"
MYSQL_PASS = ""


def raw_mysql_data():
    connection = None
    try:
        # 1. Connect
        connection = pymysql.connect(
            host=MYSQL_HOST,
            port=MYSQL_PORT,
            user=MYSQL_USER,
            password=MYSQL_PASS,
            autocommit=True,
        )
        print(f"Connected to MySQL at {MYSQL_HOST}:{MYSQL_PORT}")

        # 1.1 Setup Database and Table
        cursor = connection.cursor()
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

        # 2. Store
        data_tuple = sensor_data_helper.create_insert_data_tuple(
            sensor_data_helper.generate_random_sensor_data()
        )
        insert_sql = f"INSERT INTO {MYSQL_TABLE} VALUES (%s, %s, %s, %s, %s, %s)"
        cursor.execute(insert_sql, data_tuple)
        print("Data stored in MySQL")

        # 3. Retrieve
        with connection.cursor(pymysql.cursors.DictCursor) as dict_cursor:
            dict_cursor.execute(f"SELECT * FROM {MYSQL_TABLE}")
            results = dict_cursor.fetchall()
            if results:
                print(f"Retrieved {len(results)} records")
                print(results)

        # 4. Delete
        try:
            cursor.execute(f"TRUNCATE TABLE {MYSQL_TABLE}")
            print("Table truncated")
        except pymysql.MySQLError as e:
            print(f"Error purging data: {e}")

    except pymysql.MySQLError as e:
        print(f"MySQL Connection/Execution Error: {e}")
    finally:
        if connection:
            connection.close()
            print("Connection closed")
