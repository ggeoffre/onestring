# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 1: python/data
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Create standalone database CRUD operations for multiple databases with a helper module.

---

## File Structure
```
python/data/
├── main.py
├── sensor_data_helper.py
├── mongo_data.py
├── cassandra_data.py
├── mysql_data.py
├── postgres_data.py
├── redis_data.py
└── requirements.txt
```

---

## Prompt 1: Helper Module (sensor_data_helper.py)

```
Create a Python module called sensor_data_helper.py with:

1. Constants at the top:
   - sample_timestamp = 1768570200
   - sample_location = "den"
   - sample_sensor = "bmp280"
   - sample_measurement = "temperature"
   - sample_units = "C"
   - sample_value = 22.3
   - min_temperature = 22.4
   - max_temperature = 32.1
   - temperature_precision = 1

2. A dictionary called sensor_data with the above values as a template

3. Function generate_random_sensor_data():
   - Makes a copy of sensor_data
   - Updates "recorded" with current Unix timestamp (int)
   - Updates "value" with random float between min and max temperature, rounded to 1 decimal
   - Returns the dictionary

4. Function create_insert_data_tuple(data):
   - Takes a dictionary
   - Returns a tuple: (recorded, location, sensor, measurement, units, value)
   - Converts recorded to int and value to float

5. Function json_default(obj):
   - Handles special JSON types for serialization
   - Converts Decimal to float
   - Converts datetime to isoformat
   - Converts ObjectId to string (if it exists)
   - Raises TypeError for unknown types

6. Function json_list_to_csv(json_string_list):
   - Takes a list of dictionaries
   - Returns a CSV string with headers in first line
   - Uses keys from first dict as headers
   - Each row on new line
   - Handles commas in values by wrapping in quotes
```

---

## Prompt 2: MongoDB CRUD (mongo_data.py)

```
Create a Python file called mongo_data.py with:

Configuration at top:
- Import os, pymongo, sensor_data_helper
- MONGO_HOST from environment variable DATA_HOSTNAME (default "localhost", lowercase)
- MONGO_PORT = 27017
- MONGO_DB = "sensor_data_db"
- MONGO_COLLECTION = "sensor_data"

Function raw_mongo_data() that does all 4 operations in sequence:

1. Connect:
   - Use pymongo.MongoClient with f-string: "mongodb://{MONGO_HOST}:{MONGO_PORT}/"
   - Set serverSelectionTimeoutMS=5000
   - Get db and collection
   - Verify with client.admin.command("ping")
   - Print "Connected to Mongo at {MONGO_HOST}:{MONGO_PORT}"
   - Return early if connection fails

2. Store:
   - Generate random sensor data using sensor_data_helper
   - Insert using collection.insert_one()
   - Print "Record stored successfully"
   - Catch and print exceptions

3. Retrieve:
   - Find all documents, exclude "_id" field
   - Convert cursor to list
   - Print count and results
   - Print "No matching records found" if empty
   - Catch and print exceptions

4. Delete:
   - Use collection.delete_many({}) to remove all
   - Print "Removed all records"
   - Catch and print exceptions
   - Close client in finally block
   - Print "Connection closed"
```

---

## Prompt 3: Cassandra CRUD (cassandra_data.py)

```
Create cassandra_data.py with:

Configuration:
- Import os, sensor_data_helper, Cluster from cassandra.cluster, DriverException
- CASS_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- CASS_PORT = 9042
- CASS_KEYSPACE = "sensor_data_db"
- CASS_TABLE = "sensor_data"

Function raw_cassandra_data():
- Create cluster = Cluster([CASS_HOST], port=CASS_PORT)
- Put all operations in try-except-finally

1. Connect:
   - session = cluster.connect()
   - Print connected message

2. Setup keyspace and table:
   - Create keyspace with SimpleStrategy, replication_factor 1
   - Set keyspace
   - Create table with columns: recorded (bigint), location, sensor, measurement, units (text), value (double)
   - PRIMARY KEY: ((location), recorded, sensor)

3. Store:
   - Generate random data
   - INSERT using parameterized query with %s placeholders
   - Use create_insert_data_tuple() for values
   - Print "Data stored successfully"

4. Retrieve:
   - SELECT * from table
   - Convert rows to list of dicts using row._asdict()
   - Print count and results

5. Delete:
   - TRUNCATE table
   - Print "Table truncated"

Catch DriverException and print error
Finally: cluster.shutdown() and print "Connection closed"
```

---

## Prompt 4: MySQL CRUD (mysql_data.py)

```
Create mysql_data.py with:

Configuration:
- Import os, pymysql, pymysql.cursors, sensor_data_helper
- MYSQL_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- MYSQL_PORT = 3306
- MYSQL_DB = "sensor_data_db"
- MYSQL_TABLE = "sensor_data"
- MYSQL_USER = "root"
- MYSQL_PASS = ""

Function raw_mysql_data():
- Initialize connection = None
- Put all in try-except-finally

1. Connect:
   - pymysql.connect with host, port, user, password, autocommit=True
   - Print connected message

2. Setup database and table:
   - Get cursor
   - CREATE DATABASE IF NOT EXISTS
   - connection.select_db(database)
   - CREATE TABLE with: recorded (BIGINT), location, sensor, measurement (VARCHAR(255)), units (VARCHAR(10)), value (DECIMAL(5,2))

3. Store:
   - Generate random data and convert to tuple
   - INSERT INTO with %s placeholders
   - cursor.execute with tuple
   - Print "Data stored successfully"

4. Retrieve:
   - Use DictCursor context manager
   - SELECT * from table
   - Fetch all results
   - Convert Decimal value to float in each row
   - Print count and results

5. Delete:
   - TRUNCATE TABLE
   - Print "Table truncated"

Catch pymysql.MySQLError and print
Finally: close connection if not None
```

---

## Prompt 5: PostgreSQL CRUD (postgres_data.py)

```
Create postgres_data.py with:

Configuration:
- Import os, psycopg2, psycopg2.extras.RealDictCursor, sensor_data_helper
- POSTGRES_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- POSTGRES_PORT = 5432
- POSTGRES_DB = "sensor_data_db"
- POSTGRES_TABLE = "sensor_data"
- POSTGRES_USER = "postgres"
- POSTGRES_PASS = ""

Function raw_postgres_data():
- Initialize connection = None
- Use try-except-finally

1. Connect to "postgres" database first:
   - Connect with host, port, database="postgres", user, password
   - Set autocommit = True
   - Get cursor
   - Check if target database exists with: SELECT 1 FROM pg_catalog.pg_database WHERE datname = %s
   - If not exists, CREATE DATABASE
   - Print "Database created" if new

2. Close and reconnect to target database:
   - Close admin connection
   - Connect to POSTGRES_DB
   - Set autocommit = True
   - Print connected message

3. Create table:
   - CREATE TABLE IF NOT EXISTS with: recorded (BIGINT), location, sensor, measurement, units (VARCHAR(255)), value (NUMERIC(10,2))

4. Store:
   - Generate data and convert to tuple
   - INSERT INTO with %s placeholders
   - Print "Data stored successfully"

5. Retrieve:
   - Use RealDictCursor
   - SELECT * from table
   - Fetchall
   - Convert Decimal value to float in each row
   - Print count and results

6. Delete:
   - TRUNCATE TABLE
   - Print "Table truncated"

Catch psycopg2.Error and print
Finally: close connection if exists
```

---

## Prompt 6: Redis CRUD (redis_data.py)

```
Create redis_data.py with:

Configuration:
- Import json, os, redis, sensor_data_helper
- REDIS_HOST from DATA_HOSTNAME (default "localhost", lowercase)
- REDIS_PORT = 6379
- REDIS_KEY_BASE = "location:den:list"

Function raw_redis_data():

1. Connect:
   - redis.StrictRedis with host, port, decode_responses=True, socket_timeout=5
   - Call redis_client.ping() to verify
   - Print connected message
   - Catch ConnectionError and return early

2. Store:
   - Generate random data
   - Convert to JSON string using json.dumps with json_default handler
   - Use RPUSH to add to list
   - Print "Data stored successfully"
   - Catch and print exceptions

3. Retrieve:
   - Use LRANGE to get all items (0, -1)
   - Parse each JSON string to dict
   - Print count and results
   - Print "No data found" if empty
   - Catch and print exceptions

4. Delete:
   - DELETE the key
   - Print "Key deleted"
   - Catch and print exceptions

5. Close:
   - Check if ping works
   - Disconnect connection pool
   - Print "Connection closed"
   - Catch ConnectionError silently
```

---

## Prompt 7: Main Script (main.py)

```
Create main.py that:
- Imports all the raw functions:
  * from cassandra_data import raw_cassandra_data
  * from mongo_data import raw_mongo_data
  * from mysql_data import raw_mysql_data
  * from postgres_data import raw_postgres_data
  * from redis_data import raw_redis_data

- Calls each function with a header:
  print("\nCASSANDRA\n#########")
  raw_cassandra_data()
  
  (Repeat for MONGO, MYSQL, POSTGRES, REDIS)
```

---

## Prompt 8: Requirements File

```
Create requirements.txt with these libraries:
pymongo
cassandra-driver
pymysql
psycopg2-binary
redis
```

---

## Expected Results

When run, this code will:
1. Connect to each database type
2. Create schema/keyspace/database if needed
3. Insert one random sensor reading
4. Retrieve and print all data
5. Delete all data
6. Close connections gracefully

Each database operation is self-contained in its own function with full CRUD cycle.

---

## Key Patterns Used

- Configuration from environment variables with defaults
- Try-except-finally for resource management
- Print statements for operation visibility
- Helper functions for data generation and transformation
- Consistent pattern across all databases (connect, setup, store, retrieve, delete, close)
