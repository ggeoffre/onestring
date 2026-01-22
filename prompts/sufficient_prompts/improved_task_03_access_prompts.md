# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 3: ACCESS - Improved Prompts
## Python Abstract Base Classes + Database Implementations

---

## Core Prompts (Minimal Set)

### Prompt 1: Define the Interface
```
Create a file called sensor_data_access_protocol.py with:
- An abstract base class called SensorDataAccess using Python's ABC module
- Three abstract methods that all database implementations must provide:
  1. log_sensor_data(json_data: str) -> None - stores one sensor reading
  2. fetch_sensor_data() -> list[str] - retrieves all sensor readings
  3. purge_sensor_data() -> None - deletes all sensor readings
- Include proper type hints
- Don't implement the methods, just define them as abstract
```

### Prompt 2: MongoDB Implementation Structure
```
Create a file called mongo_data_access.py that:
- Imports the SensorDataAccess interface
- Creates a class called MongoDataAccess that implements SensorDataAccess
- Add a method called get_connection() that:
  * Reads DATA_HOSTNAME from environment variables (default "localhost")
  * Connects to MongoDB on port 27017
  * Returns the client, database ("sensor_data_db"), and collection ("sensor_data")
  * Includes a 5-second timeout
  * Prints connection success
- Don't implement the three abstract methods yet
```

### Prompt 3: MongoDB Log Method
```
In MongoDataAccess, implement log_sensor_data() that:
- Calls get_connection() to connect
- Takes json_data which is a dictionary (not a string despite the type hint)
- Inserts the dictionary directly into MongoDB
- Prints "Data stored successfully"
- Catches and prints any errors
- Returns None
```

### Prompt 4: MongoDB Fetch Method
```
In MongoDataAccess, implement fetch_sensor_data() that:
- Calls get_connection() to connect
- Retrieves all documents from the collection
- Removes the "_id" field from each document (MongoDB adds this automatically)
- Returns a list of dictionaries
- If no data found, returns empty list
- Prints how many records were retrieved
- Catches and prints any errors
```

### Prompt 5: MongoDB Purge Method
```
In MongoDataAccess, implement purge_sensor_data() that:
- Calls get_connection() to connect
- Deletes all documents from the collection using delete_many({})
- Prints "Sensor data purged from MongoDB"
- Catches and prints any errors
- Returns None
```

### Prompt 6: Main Test File
```
Create main.py that:
- Imports MongoDataAccess
- Creates an instance of MongoDataAccess
- Calls each of the three methods to test them:
  1. log_sensor_data() with sample data: {"recorded": 1768570200, "location": "den", "sensor": "bmp280", "measurement": "temperature", "units": "C", "value": 22.3}
  2. fetch_sensor_data() and print the results
  3. purge_sensor_data() to clean up
- Include clear print statements showing what's being tested
```

### Prompt 7: Requirements File
```
Create requirements.txt with:
pymongo
```

---

## Optional Enhancement Prompts (Additional Databases)

### Prompt 8: Redis Implementation
```
Create redis_data_access.py that implements SensorDataAccess:
- Use the same pattern as MongoDataAccess
- Connection method connects to Redis on port 6379 using the redis library
- Store data using RPUSH to a key called "sensor_data:list"
- Store each reading as a JSON string
- Fetch retrieves the list using LRANGE and parses each JSON string back to a dictionary
- Purge deletes the key using DELETE
- Include 5-second timeout and error handling
```

### Prompt 9: PostgreSQL Implementation
```
Create postgres_data_access.py that implements SensorDataAccess:
- Use psycopg2 library
- Connection method:
  * First connects to "postgres" database to create "sensor_data_db" if needed
  * Then connects to "sensor_data_db"
  * Creates table "sensor_data" if needed with columns: recorded (BIGINT), location (VARCHAR), sensor (VARCHAR), measurement (VARCHAR), units (VARCHAR), value (NUMERIC)
- Store method inserts using parameterized query: INSERT INTO sensor_data VALUES (%s, %s, %s, %s, %s, %s)
- Fetch retrieves all rows and returns as list of dictionaries
- Purge uses TRUNCATE TABLE
- Include error handling
```

### Prompt 10: MySQL Implementation
```
Create mysql_data_access.py that implements SensorDataAccess:
- Use pymysql library
- Connection method:
  * Connects to MySQL on port 3306 (user: root, no password)
  * Creates database "sensor_data_db" if needed
  * Creates table "sensor_data" if needed with same columns as PostgreSQL
- Store method inserts one row using parameterized query
- Fetch retrieves all rows using DictCursor to get dictionaries
- Purge uses TRUNCATE TABLE
- Include error handling
```

### Prompt 11: Cassandra Implementation
```
Create cassandra_data_access.py that implements SensorDataAccess:
- Use cassandra-driver library
- Connection method:
  * Connects to Cassandra cluster on port 9042
  * Creates keyspace "sensor_data_db" if needed with SimpleStrategy replication
  * Creates table "sensor_data" with PRIMARY KEY ((location), recorded, sensor)
  * Returns a session object
- Store method inserts using CQL with parameterized query
- Fetch retrieves all rows and converts to list of dictionaries
- Purge uses TRUNCATE
- Include error handling
```

---

## Integration Prompts

### Prompt 12: Factory Pattern
```
Update main.py to add a function called get_data_access() that:
- Reads an environment variable called DATA_ACCESS (default to "mongo")
- Returns the appropriate implementation:
  * "mongo" -> MongoDataAccess()
  * "redis" -> RedisDataAccess()
  * "postgres" -> PostgresDataAccess()
  * "mysql" -> MySQLDataAccess()
  * "cassandra" -> CassandraDataAccess()
- Raises an error for unsupported types
- Returns the abstract type SensorDataAccess

Then modify main.py to:
- Call get_data_access() instead of directly creating MongoDataAccess()
- Test that the returned object works the same way
```

### Prompt 13: Schema Creation Helper
```
For databases that need schema creation (Postgres, MySQL, Cassandra):
- Add a global variable called TABLE_CREATION_CHECKED = False
- In get_connection(), only create the database/table/keyspace on the first call
- Set TABLE_CREATION_CHECKED = True after creation
- On subsequent calls, skip the creation logic
- This avoids recreating the schema on every operation
```

### Prompt 14: Data Helper Integration
```
Create sensor_data_helper.py with a function called create_insert_data_tuple() that:
- Takes a dictionary
- Returns a tuple in this order: (recorded, location, sensor, measurement, units, value)
- Converts "recorded" to int and "value" to float
- This helps prepare data for SQL/CQL INSERT statements

Update all SQL database implementations (Postgres, MySQL, Cassandra) to:
- Import and use create_insert_data_tuple()
- Pass the tuple to the parameterized INSERT query
```

---

## Expected Sufficiency Improvement

**Original Rating: 30%**
**Improved Rating: 70-75%**

### What These Prompts Provide:
✅ Clear interface definition with ABC pattern
✅ Step-by-step implementation guidance for each database
✅ Database-specific connection patterns explained
✅ Schema creation logic specified
✅ Data transformation approach clarified
✅ Factory pattern for database switching
✅ Error handling expectations set

### What Still Requires Developer Judgment:
- Exact error message wording
- Import statement organization
- Connection pooling strategies (not needed for this simple case)
- Advanced query optimization
- Detailed docstrings and comments

### Key Improvements from Original:
- Separated connection logic from CRUD operations
- Explicitly stated database-specific quirks (schema creation, primary keys, etc.)
- Provided clear guidance on data type handling
- Specified the factory pattern implementation
- Clarified the global state management for schema creation

### Why This Balance Works:
- Each database's unique requirements are clearly stated
- Implementation pattern is consistent across databases
- Technical details (ports, credentials, timeouts) are specified
- Developers understand the "why" behind architectural decisions
- No line-by-line code dictation, just clear requirements
