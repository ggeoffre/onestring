# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 5: python/legos/flask-app
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Integrate Flask API with multiple database implementations and helper utilities.

---

## File Structure
```
python/legos/flask-app/
├── app.py
├── sensor_data_access_protocol.py
├── mongo_data_access.py
├── cassandra_data_access.py
├── mysql_data_access.py
├── postgres_data_access.py
├── redis_data_access.py
├── sensor_data_helper.py
└── requirements.txt
```

---

## Prompt 1: Abstract Protocol (sensor_data_access_protocol.py)

```
Create sensor_data_access_protocol.py - same as Project 4:

Use ABC with three abstract methods:
- log_sensor_data(self, json_data: str) -> None
- fetch_sensor_data(self) -> List[str]
- purge_sensor_data(self) -> None
```

---

## Prompt 2: Helper Module (sensor_data_helper.py)

```
Create sensor_data_helper.py - same core functions as Project 1:

Include:
- generate_random_sensor_data() - returns dict with current timestamp and random temperature
- create_insert_data_tuple(data) - converts dict to tuple for SQL inserts
- json_default(obj) - handles Decimal, datetime, ObjectId for JSON serialization
- json_list_to_csv(json_string_list) - converts list of dicts to CSV string with escaping

Plus constants for sample sensor data.
```

---

## Prompt 3: Flask App with Data Access Integration (app.py)

```
Create app.py with Flask routes integrated with data access layer:

1. Imports:
   - import os
   - import sensor_data_helper
   - from flask import Flask, Response, jsonify, request
   - from cassandra_data_access import CassandraDataAccess
   - from mongo_data_access import MongoDataAccess
   - from mysql_data_access import MySQLDataAccess
   - from postgres_data_access import PostgresDataAccess
   - from redis_data_access import RedisDataAccess
   - from sensor_data_access_protocol import SensorDataAccess

2. Factory function get_data_access() -> SensorDataAccess:
   - Read DATA_ACCESS from environment (default "mongo")
   - Return appropriate implementation based on value
   - Raise ValueError for unsupported types

3. Create Flask app: app = Flask(__name__)

4. Route GET /:
   - Return JSON: {"message": "Flask API Server is running!"}

5. Route POST /echo:
   - Get JSON from request (force=True)
   - Return 400 error if None: {"error": "No valid JSON provided"}
   - Return data as JSON

6. Route POST /log:
   - Get data access instance
   - Get JSON from request
   - Print the data
   - Call data_access.log_sensor_data(data)
   - Return 400 if None: {"error": "No valid JSON provided"}
   - Return success: {"message": "Data logged successfully"}

7. Route GET /report:
   - Wrap in try-except
   - Get data access instance
   - Call fetch_sensor_data()
   - Convert to CSV using sensor_data_helper.json_list_to_csv()
   - Return 404 if no data: {"error": "No data available"}
   - Return CSV as Response with:
     * mimetype="text/csv"
     * Content-Disposition header for download
   - Return 500 error on exception

8. Route GET or POST /purge:
   - Get data access instance
   - Call purge_sensor_data()
   - Return: {"message": "Data purge sequence complete"}

9. Main block:
   - Run on host="0.0.0.0", port=8080, debug=False
```

---

## Prompt 4: MongoDB Implementation (mongo_data_access.py)

```
Create mongo_data_access.py with full database operations:

Configuration:
- MONGO_HOST from DATA_HOSTNAME env (default "localhost", lowercase)
- MONGO_PORT = 27017
- MONGO_DB = "sensor_data_db"
- MONGO_COLLECTION = "sensor_data"

Class MongoDataAccess(SensorDataAccess):

Method get_connection():
- Connect with pymongo.MongoClient, 5-second timeout
- Verify with ping command
- Print connected message
- Return (client, db, collection) tuple

Method close_connection(client, db, collection):
- Close client
- Print "MongoDB connection closed"

Method log_sensor_data(json_data):
- Get connection
- Insert data with collection.insert_one()
- Print "Record stored successfully"
- Catch and print exceptions

Method fetch_sensor_data():
- Get connection
- Find all, exclude "_id" field
- Convert cursor to list
- Print count or "No matching records found"
- Return list
- Catch and print exceptions, return empty list

Method purge_sensor_data():
- Get connection
- Delete all with delete_many({})
- Print "Sensor data purged from MongoDB"
- Catch and print exceptions
```

---

## Prompt 5: Cassandra Implementation (cassandra_data_access.py)

```
Create cassandra_data_access.py:

Configuration:
- CASS_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- CASS_PORT = 9042
- CASS_KEYSPACE = "sensor_data_db"
- CASS_TABLE = "sensor_data"
- Global TABLE_CREATION_CHECKED = False

Class CassandraDataAccess:
- Class variable: cluster = Cluster([CASS_HOST], port=CASS_PORT)

Method get_connection():
- Use global TABLE_CREATION_CHECKED
- Connect and get session
- Print connected message
- If not TABLE_CREATION_CHECKED:
  * Create keyspace with SimpleStrategy
  * Set keyspace
  * Create table with PRIMARY KEY ((location), recorded, sensor)
  * Set TABLE_CREATION_CHECKED = True
- Else: just set keyspace
- Return session
- Catch exceptions and return None

Method log_sensor_data(json_data):
- Get connection
- INSERT with parameterized query
- Use sensor_data_helper.create_insert_data_tuple()
- Print "Data stored successfully"

Method fetch_sensor_data():
- Get connection
- SELECT * from table
- Convert rows to dicts with row._asdict()
- Print count and results
- Return list or empty list

Method purge_sensor_data():
- Get connection
- TRUNCATE table
- Print "Table truncated"
```

---

## Prompt 6: MySQL Implementation (mysql_data_access.py)

```
Create mysql_data_access.py:

Configuration:
- MYSQL_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- MYSQL_PORT = 3306
- MYSQL_DB = "sensor_data_db"
- MYSQL_TABLE = "sensor_data"
- MYSQL_USER = "root"
- MYSQL_PASS = ""
- Global TABLE_CREATION_CHECKED = False

Class MySQLDataAccess:

Method _ensure_schema(connection):
- Use global TABLE_CREATION_CHECKED
- Return early if already checked
- CREATE DATABASE IF NOT EXISTS
- Select database
- CREATE TABLE with BIGINT, VARCHAR, DECIMAL columns
- Set TABLE_CREATION_CHECKED = True

Method get_cursor():
- Connect with pymysql (autocommit=True)
- Call _ensure_schema()
- Select database
- Return (connection, cursor) tuple

Method log_sensor_data(json_data):
- Get cursor
- Convert data to tuple
- INSERT with parameterized query
- Print "Data stored successfully"

Method fetch_sensor_data():
- Get cursor
- Use DictCursor
- SELECT * from table
- Convert Decimal values to float
- Print count and results
- Return list or empty list

Method purge_sensor_data():
- Get cursor
- TRUNCATE table
- Print "Table truncated"
```

---

## Prompt 7: PostgreSQL Implementation (postgres_data_access.py)

```
Create postgres_data_access.py:

Configuration:
- POSTGRES_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- POSTGRES_PORT = 5432
- POSTGRES_DB = "sensor_data_db"
- POSTGRES_TABLE = "sensor_data"
- POSTGRES_USER = "postgres"
- POSTGRES_PASS = ""
- Global TABLE_CREATION_CHECKED = False

Class PostgresDataAccess:

Method _ensure_schema():
- Use global TABLE_CREATION_CHECKED
- Return early if already checked
- Connect to "postgres" database first
- Check if target database exists
- CREATE DATABASE if needed
- Close admin connection
- Connect to target database
- CREATE TABLE with BIGINT, VARCHAR, NUMERIC columns
- Set TABLE_CREATION_CHECKED = True

Method get_cursor():
- Call _ensure_schema()
- Connect to target database with autocommit=True
- Return (connection, cursor) tuple

Method log_sensor_data(json_data):
- Get cursor
- Convert to tuple
- INSERT with parameterized query
- Print "Data stored successfully"

Method fetch_sensor_data():
- Get cursor
- Use RealDictCursor
- SELECT * from table
- Convert Decimal to float
- Print count and results
- Return list or empty list

Method purge_sensor_data():
- Get cursor
- TRUNCATE table
- Print "Table truncated"
```

---

## Prompt 8: Redis Implementation (redis_data_access.py)

```
Create redis_data_access.py:

Configuration:
- REDIS_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- REDIS_PORT = 6379
- REDIS_KEY_BASE = "location:den:list"

Class RedisDataAccess:

Method get_connection():
- Create StrictRedis with decode_responses=True, socket_timeout=5
- Call ping() to verify
- Print connected message
- Return client
- Raise on ConnectionError

Method log_sensor_data(json_data):
- Get connection
- Convert data to JSON string with json_default handler
- RPUSH to list
- Print "Data stored successfully"

Method fetch_sensor_data():
- Get connection
- LRANGE to get all items
- Parse each JSON string to dict
- Print count and results or "No data found"
- Return list or empty list

Method purge_sensor_data():
- Get connection
- DELETE the key
- Print "Key deleted"
```

---

## Prompt 9: Requirements File

```
Create requirements.txt with:
Flask
pymongo
cassandra-driver
pymysql
psycopg2-binary
redis
```

---

## Expected Behavior

This integrated system:
- Runs Flask API on port 8080
- Routes use factory pattern to get database implementation
- POST /log stores data in selected database
- GET /report retrieves and exports as CSV
- GET or POST /purge clears all data
- Database selection via DATA_ACCESS environment variable
- Each database creates schema on first use
- Helper functions handle data generation and CSV conversion

---

## Key Integration Patterns

- Factory pattern for database selection
- Abstract protocol for consistent interface
- Dependency injection (routes call factory)
- CSV conversion utility
- Error handling with HTTP status codes
- Environment-based configuration
- Global state for schema creation (one-time setup)
- Proper HTTP response types (JSON vs CSV)
