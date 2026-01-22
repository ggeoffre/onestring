# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 2: API - Improved Prompts
## Python + Flask Web Service

---

## Core Prompts (Minimal Set)

### Prompt 1: Basic Flask Setup
```
Create a Flask application in app.py that:
- Imports Flask
- Creates an app that listens on all IP addresses (0.0.0.0) and port 8080
- Has a global variable to store sensor data as a list of dictionaries
- Runs with debug=True when the script is executed directly
- Don't worry about routes yet - just the basic app structure
```

### Prompt 2: Root Route
```
Add a route for GET / that:
- Returns a JSON response
- The response should have one field: "message" with value "Flask API is running"
- Use jsonify to create the JSON response
```

### Prompt 3: Echo Route
```
Add a route for POST /echo that:
- Accepts JSON in the request body
- Returns the exact same JSON back to the caller
- If the request isn't JSON, return an error message "Invalid JSON" with status code 400
- The route should only accept POST (not GET)
```

### Prompt 4: Log Route
```
Add a route for POST /log that:
- Accepts JSON in the request body
- Appends the JSON to the global list of sensor data
- Returns a JSON response: {"message": "Data logged successfully"}
- If the request isn't JSON, return an error message "Invalid JSON" with status code 400
- The route should only accept POST (not GET)
```

### Prompt 5: Report Route
```
Add a route for GET /report that:
- Gets all sensor data from the global list
- Converts it to CSV format (headers in first line, one data row per line)
- Returns the CSV as a downloadable file named "sensor_report.csv"
- Sets the Content-Type header to "text/csv"
- Sets the Content-Disposition header to make it downloadable
- If there's no data, return an empty CSV with just the headers
```

### Prompt 6: Purge Route
```
Add a route for /purge that:
- Accepts both GET and POST methods
- Clears all data from the global list
- Returns a JSON response: {"message": "All data purged"}
```

### Prompt 7: CSV Conversion Helper
```
Create a helper function called json_to_csv() that:
- Takes a list of dictionaries
- Returns a CSV string
- Uses the keys from the first dictionary as the CSV headers
- Handles the case where the list is empty (just return headers: "recorded,location,sensor,measurement,units,value")
- Don't use the csv library - just string operations
- Put this function in the app.py file before the routes
```

### Prompt 8: Requirements File
```
Create a requirements.txt with:
Flask
```

---

## Optional Enhancement Prompts

### Prompt 9: Static Sample Data
```
Add to the top of app.py:
- A constant dictionary called SAMPLE_SENSOR_DATA with this structure:
  {
    "recorded": 1768570200,
    "location": "den",
    "sensor": "bmp280",
    "measurement": "temperature",
    "units": "C",
    "value": 22.3
  }
- When the app starts, automatically add 3 copies of this sample data (with slightly different values) to the global list
```

### Prompt 10: Error Handling
```
Update all routes to:
- Wrap their logic in try-except blocks
- If an error occurs, return a JSON response: {"error": "An error occurred", "details": [error message]}
- Return status code 500 for server errors
- Still return 400 for invalid JSON as specified before
```

### Prompt 11: Test Data Generation
```
Add a route for POST /generate that:
- Takes an optional "count" parameter from the JSON body (default to 1 if not provided)
- Generates that many random sensor readings (use random timestamp and temperature)
- Adds them all to the global list
- Returns: {"message": "Generated X readings", "count": X}
```

---

## Integration Prompt

### Prompt 12: Complete Flask API
```
Create a complete Flask API with all routes working together:

Structure:
- app.py has all routes and the csv conversion helper
- requirements.txt has Flask

Routes should work in this sequence:
1. GET / - shows API is running
2. POST /generate with {"count": 5} - creates 5 random readings
3. POST /log with a single reading - adds one more
4. POST /echo with any JSON - returns it unchanged
5. GET /report - downloads all data as CSV
6. GET /purge - clears everything
7. GET /report again - returns empty CSV (just headers)

Make sure:
- The app binds to 0.0.0.0:8080
- All JSON responses use jsonify()
- CSV download has proper headers
- Global data list persists between requests (but resets when app restarts)
```

---

## Testing Prompt

### Prompt 13: Create Test Script
```
Create a test.sh bash script that uses curl to test all endpoints:

1. Test GET / - expect success message
2. Test POST /echo with sample JSON - expect same JSON back
3. Test POST /log with sample sensor data - expect success message
4. Test POST /log again with different data - expect success message
5. Test GET /report - save the CSV to report.csv file
6. Test GET /purge - expect purge confirmation
7. Test GET /report again - expect empty result

Print clear messages before each test like "Testing GET /..."
Print the response from each test
Assume the API is running on localhost:8080
```

---

## Expected Sufficiency Improvement

**Original Rating: 35%**
**Improved Rating: 75-80%**

### What These Prompts Provide:
✅ Explicit route behavior and responses
✅ HTTP method constraints specified
✅ Error handling approach clearly defined
✅ CSV conversion logic explained
✅ Response format standardization
✅ Global state management approach
✅ Clear request/response examples

### What Still Requires Developer Judgment:
- Exact Flask import style (from flask import Flask vs import flask)
- Detailed exception handling patterns
- Logging format preferences
- Code organization and comments

### Why This Balance Works:
- Routes are functionally specified without prescribing implementation
- Error cases are identified but not over-engineered
- CSV logic is explained conceptually, not line-by-line
- Testing approach guides verification without rigid test requirements
- Developers retain control over code style and structure
