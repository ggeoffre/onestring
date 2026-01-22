# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 4: LEGOS - Improved Prompts
## Integrated Flask + Database Access + Helper Functions

---

## Core Integration Prompts (Minimal Set)

### Prompt 1: Project Structure Setup
```
Create a Flask application that integrates three components:
1. Flask routes (from Task 2)
2. Database access layer (from Task 3)
3. Helper functions for data manipulation

File structure:
- app.py (Flask routes)
- sensor_data_access_protocol.py (abstract interface)
- mongo_data_access.py (MongoDB implementation)
- sensor_data_helper.py (helper functions)
- requirements.txt

Don't implement the integration yet, just set up the files and imports.
```

### Prompt 2: Helper Functions for Integration
```
In sensor_data_helper.py, create these functions:

1. generate_random_sensor_data() -> dict
   - Returns a dictionary with current Unix timestamp
   - Random temperature value between 22.4 and 32.1 (1 decimal place)
   - Fixed values: location="den", sensor="bmp280", measurement="temperature", units="C"

2. json_list_to_csv(data_list: list) -> str
   - Takes a list of dictionaries
   - Returns CSV string with headers in first line
   - If list is empty, return just headers: "recorded,location,sensor,measurement,units,value"

3. create_insert_data_tuple(data_dict: dict) -> tuple
   - Takes a dictionary
   - Returns tuple: (recorded, location, sensor, measurement, units, value)
   - Converts recorded to int and value to float

All functions should handle errors gracefully and print error messages if needed.
```

### Prompt 3: Update Database Access for Integration
```
Modify mongo_data_access.py:

In log_sensor_data():
- Accept either a dictionary OR a JSON string
- If it's a dictionary, use it directly
- If it's a string, parse it to a dictionary first
- Then insert into MongoDB

In fetch_sensor_data():
- Return a list of dictionaries (not strings despite the type hint)
- Make sure to remove the "_id" field MongoDB adds
- This makes it compatible with the CSV converter

Keep the same connection logic and error handling.
```

### Prompt 4: Flask Routes with Database Integration
```
In app.py, create Flask routes that use the database access layer:

Setup:
- Import the get_data_access() factory function
- Import sensor_data_helper functions
- Create the Flask app on 0.0.0.0:8080

Routes:
1. GET / - Return {"message": "Flask API with Database is running"}

2. POST /log - 
   - Get JSON from request
   - Call data_access.log_sensor_data() with the JSON
   - Return {"message": "Data logged successfully"}

3. GET /report - 
   - Call data_access.fetch_sensor_data() to get all data
   - Convert to CSV using json_list_to_csv()
   - Return as downloadable file "sensor_report.csv"

4. POST or GET /purge - 
   - Call data_access.purge_sensor_data()
   - Return {"message": "All data purged"}

5. POST /echo - 
   - Return the exact JSON received

At the top of each route, get the data access object by calling get_data_access().
```

### Prompt 5: Database Selection via Environment
```
Update app.py to support multiple databases:

At the top of the file:
- Import all database implementations (MongoDataAccess, RedisDataAccess, etc.)
- Create a get_data_access() function that:
  * Reads DATA_ACCESS environment variable (default "mongo")
  * Returns appropriate database implementation
  * Supports: mongo, redis, postgres, mysql, cassandra

Each route should call get_data_access() to get the right database connection based on the environment variable.

This allows switching databases without changing code, just by setting DATA_ACCESS="redis" or DATA_ACCESS="postgres", etc.
```

### Prompt 6: Requirements File
```
Create requirements.txt with:
Flask
pymongo
redis
psycopg2-binary
pymysql
cassandra-driver

These support Flask and all five database options.
```

---

## Optional Enhancement Prompts

### Prompt 7: Automatic Data Generation Route
```
Add a route POST /generate to app.py:
- Takes optional "count" parameter from JSON (default 1)
- Generates that many random sensor readings using generate_random_sensor_data()
- Logs each one using data_access.log_sensor_data()
- Returns {"message": "Generated X readings", "count": X}

This makes testing easier by creating sample data.
```

### Prompt 8: Health Check Route
```
Add a route GET /health to app.py:
- Try to connect to the database
- If successful, return {"status": "healthy", "database": [database type]}
- If failed, return {"status": "unhealthy", "error": [error message]} with status 503
- This helps verify the database connection is working
```

### Prompt 9: Django Alternative
```
Create a second version using Django instead of Flask:

In django_app.py:
- Configure Django settings inline (no separate settings.py needed)
- Use Django's settings.configure() at the top
- Create the same 5 routes as views
- Use JsonResponse instead of jsonify
- Support both GET and POST for /purge using request.method check
- Use the same database access layer and helpers (no changes needed there)

This demonstrates the database access layer works with any framework.
```

### Prompt 10: Logging Enhancement
```
Improve error visibility across all modules:

1. In sensor_data_access_protocol.py:
   - Add a method called get_database_type() that returns a string
   - Each implementation should return its type: "MongoDB", "Redis", etc.

2. In all database implementations:
   - Before each operation, print: "[Database Type] Starting [operation name]"
   - After each operation, print: "[Database Type] Completed [operation name]"
   - On errors, print: "[Database Type] Error in [operation name]: [error details]"

This makes debugging much easier without changing the interface.
```

---

## Testing Prompts

### Prompt 11: Integration Test Script
```
Create test_integration.sh that:
1. Starts with DATA_ACCESS=mongo (tests MongoDB)
2. POST /generate with {"count": 3}
3. POST /log with one manual reading
4. GET /report (should have 4 readings)
5. POST /purge
6. GET /report (should be empty)

Then repeat with DATA_ACCESS=redis if Redis is available.

Use curl for all requests.
Print clear test headers and results.
Save the CSV report to a file.
```

### Prompt 12: Multi-Database Test
```
Create test_all_databases.sh that:
- Loops through: mongo, redis, postgres, mysql, cassandra
- For each database:
  * Sets DATA_ACCESS environment variable
  * Restarts the Flask app (or assumes it reads env vars on each request)
  * Runs the same test sequence: generate, log, report, purge
  * Prints which database is being tested
  * Catches errors if a database isn't available (don't fail the whole script)

This validates that all database implementations work identically.
```

---

## Final Integration Prompt

### Prompt 13: Complete Working System
```
Verify the entire system works end-to-end:

Components:
1. sensor_data_helper.py - helper functions (no database dependency)
2. sensor_data_access_protocol.py - abstract interface
3. mongo_data_access.py (and other database implementations) - implement the interface
4. app.py - Flask routes that use the database access layer via factory pattern

Flow:
- User sends JSON to POST /log
- Flask route receives it
- Calls get_data_access() to get the right database implementation
- Calls log_sensor_data() on that implementation
- Database implementation stores it
- User calls GET /report
- Flask route calls fetch_sensor_data()
- Helper converts the list to CSV
- Flask returns downloadable file

Key Features:
- Switch databases by changing DATA_ACCESS environment variable
- No code changes needed to support different databases
- Flask routes don't know which database is being used
- Helper functions are pure utility (no database coupling)
- All errors are caught and reported gracefully

Test that:
- MongoDB works
- Redis works (if available)
- CSV export works with data from any database
- Purge works and actually clears the data
- Can switch databases mid-session by changing environment variable
```

---

## Expected Sufficiency Improvement

**Original Rating: 25%**
**Improved Rating: 70-75%**

### What These Prompts Provide:
✅ Clear integration architecture
✅ Factory pattern implementation guidance
✅ Explicit data flow description
✅ Database-agnostic route implementation
✅ Helper function reusability
✅ Environment-based configuration
✅ Multi-framework support (Flask + Django)
✅ Comprehensive testing approach

### What Still Requires Developer Judgment:
- Error message formatting
- Logging detail level
- Import organization
- Code comments and documentation
- Performance optimization
- Connection pooling strategies

### Key Improvements from Original:
- Separated concerns clearly (routes, data access, helpers)
- Explained the factory pattern purpose and implementation
- Showed how environment variables drive behavior
- Demonstrated framework independence of data layer
- Provided testing strategy that validates integration

### Why This Balance Works:
- Architecture is clearly explained without prescribing every line
- Integration points are explicit but implementation is flexible
- Developers understand the component relationships
- Testing validates the integration without rigid requirements
- Simple language describes complex architectural patterns
- Focus on "what connects to what" rather than "how to write each function"

---

## Progressive Complexity Note

These prompts build on Tasks 1-3:
- Task 1: Basic CRUD with MongoDB
- Task 2: Basic Flask routes
- Task 3: Abstract interface + multiple implementations
- Task 4: Integrates all three with clean architecture

Each task's prompts can stand alone, but Task 4 assumes familiarity with the patterns from Tasks 1-3. This progressive approach teaches:
1. Database operations
2. Web API design
3. Abstraction and polymorphism
4. System integration and architecture

The prompts maintain simplicity while building increasingly sophisticated systems.
