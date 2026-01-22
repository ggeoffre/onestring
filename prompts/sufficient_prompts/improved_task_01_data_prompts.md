# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 1: DATA - Improved Prompts
## Python + MongoDB CRUD Operations

---

## Core Prompts (Minimal Set)

### Prompt 1: Database Connection Setup
```
Create a Python function called get_mongo_connection() that:
- Connects to MongoDB using pymongo
- Gets the hostname from an environment variable called DATA_HOSTNAME (use "localhost" if not set)
- Uses port 27017
- Connects to a database called "sensor_data_db" and collection called "sensor_data"
- Returns three things: the client, database, and collection
- Has a 5-second timeout if MongoDB isn't available
- Prints "Connected to Mongo at [hostname]:27017" when successful
- Catches any errors and prints them
```

### Prompt 2: Store Data Function
```
Create a Python function called store_sensor_data() that:
- Takes a dictionary as input with these keys: recorded, location, sensor, measurement, units, value
- Calls get_mongo_connection() to connect
- Inserts the dictionary into the MongoDB collection
- Prints "Data stored successfully"
- Returns nothing (None)
- Handles errors by printing them and returning None
```

### Prompt 3: Retrieve Data Function
```
Create a Python function called retrieve_sensor_data() that:
- Calls get_mongo_connection() to connect
- Gets all documents from the collection
- Removes the MongoDB "_id" field from each document before returning
- Returns a list of dictionaries
- If no data exists, returns an empty list
- Prints how many records were found
- Handles errors by printing them and returning an empty list
```

### Prompt 4: Delete Data Function
```
Create a Python function called delete_all_sensor_data() that:
- Calls get_mongo_connection() to connect
- Deletes all documents from the collection
- Prints "All data deleted"
- Returns nothing (None)
- Handles errors by printing them and returning None
```

### Prompt 5: Main Function
```
Create a main.py file that:
- Imports all four functions above
- Generates random sensor data with current Unix timestamp and random temperature between 22.4 and 32.1
- Calls store_sensor_data() to save it
- Calls retrieve_sensor_data() to fetch all data
- Calls delete_all_sensor_data() to clean up
- Uses this data structure:
  {
    "recorded": [current Unix timestamp as integer],
    "location": "den",
    "sensor": "bmp280",
    "measurement": "temperature",
    "units": "C",
    "value": [random decimal with 1 decimal place]
  }
```

### Prompt 6: Helper Functions
```
Create a sensor_data_helper.py file with three simple functions:

1. generate_random_sensor_data() - creates the data structure from Prompt 5 with random values
2. get_current_timestamp() - returns current Unix timestamp as an integer
3. get_random_temperature() - returns a random float between 22.4 and 32.1 with 1 decimal place

Also include a constant called SAMPLE_DATA that shows an example of the data structure.
```

### Prompt 7: Requirements File
```
Create a requirements.txt file with just one line:
pymongo
```

---

## Optional Enhancement Prompts (For Extra Features)

### Prompt 8: Update Function
```
Create a Python function called update_sensor_value() that:
- Takes two parameters: a filter dictionary (to find the document) and a new value
- Calls get_mongo_connection() to connect
- Updates the "value" field of the first matching document
- Prints "Data updated successfully"
- Returns nothing (None)
- Handles errors by printing them
```

### Prompt 9: CSV Conversion
```
Create a function called sensor_data_to_csv() that:
- Takes a list of sensor data dictionaries
- Returns a CSV string with headers in the first line
- Headers should be: recorded,location,sensor,measurement,units,value
- Each data row should be on a new line
- Don't use any CSV libraries - just string manipulation
```

### Prompt 10: Data Validation
```
Create a function called validate_sensor_data() that:
- Takes a dictionary
- Checks that it has all required keys: recorded, location, sensor, measurement, units, value
- Checks that "recorded" is an integer
- Checks that "value" is a number (int or float)
- Returns True if valid, False if not
- Doesn't need to print anything
```

---

## Integration Prompt (Combines Everything)

### Prompt 11: Complete Application
```
Combine all the functions into a working application where:
1. mongo_data.py contains the four main functions (get_mongo_connection, store, retrieve, delete)
2. sensor_data_helper.py contains the helper functions
3. main.py does this sequence:
   - Generate 3 random sensor readings
   - Store each one
   - Retrieve all data and print it
   - Convert the data to CSV and print it
   - Delete all data
   - Verify deletion by retrieving again (should be empty)

Make sure to handle the case where MongoDB might not be running by catching connection errors gracefully.
```

---

## Expected Sufficiency Improvement

**Original Rating: 40%**
**Improved Rating: 75-80%**

### What These Prompts Provide:
✅ Clear function signatures and purposes
✅ Specific return types and error handling patterns
✅ Environment variable usage explained
✅ Data structure explicitly defined
✅ Connection pattern clearly specified
✅ Simple, plain English descriptions

### What Still Requires Developer Judgment:
- Exact import statement formatting
- Specific variable naming preferences
- Code organization within functions
- Comments and docstring style

### Why This Balance Works:
- Developers understand WHAT to build without being told HOW to write every line
- Technical details (ports, timeout values) are specified where they matter
- Error handling approach is clear but not prescriptive about exact implementation
- Focuses on behavior and requirements rather than syntax
