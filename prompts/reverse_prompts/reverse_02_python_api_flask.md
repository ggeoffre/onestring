# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 2: python/api/flask-app
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Create a minimal Flask API with basic routes and JSON-to-CSV conversion.

---

## File Structure
```
python/api/flask-app/
├── app.py
└── requirements.txt
```

---

## Prompt 1: Flask API with Routes (app.py)

```
Create a Flask application in app.py with:

1. Imports:
   - import json
   - from flask import Flask, Response, jsonify, request

2. A constant SENSOR_DATA_JSON:
   - Create a dictionary with: recorded (1768570200), location ("den"), sensor ("bmp280"), 
     measurement ("temperature"), units ("C"), value (22.3)
   - Convert to JSON string using json.dumps()

3. Helper function json_to_csv_safe(json_string):
   - Takes a JSON string
   - Parses it with json.loads()
   - If it's a dict, convert to list containing that dict
   - Return empty string if data is empty
   - Get headers from first dict's keys
   - Create CSV with headers in first line
   - For each item, create a row:
     * Convert each value to string
     * Escape quotes by replacing " with ""
     * Wrap in quotes if value contains comma or quote
   - Join rows with newlines
   - Return empty string if any exception occurs

4. Create Flask app: app = Flask(__name__)

5. Route GET /:
   - Return JSON: {"message": "Flask API Server is running!"}

6. Routes POST /echo and POST /log (same function):
   - Get JSON from request using request.get_json(force=True)
   - Return error 400 if data is None: {"error": "No valid JSON provided"}
   - Return the data as JSON using jsonify

7. Route GET /report:
   - Call json_to_csv_safe with SENSOR_DATA_JSON constant
   - Return as Response with:
     * CSV data as content
     * mimetype="text/csv"
     * headers with Content-Disposition: "attachment; filename=sensor_report.csv"

8. Route GET or POST /purge:
   - Return JSON: {"message": "Data purge sequence complete"}

9. Main block:
   - Run app on host="0.0.0.0", port=8080, debug=False
```

---

## Prompt 2: Requirements File

```
Create requirements.txt with:
Flask
```

---

## Expected Behavior

When run:
- GET / returns welcome message
- POST /echo or /log returns the JSON you send
- GET /report downloads a CSV file with the sample sensor data
- GET or POST /purge returns success message
- Server runs on all interfaces (0.0.0.0) on port 8080

---

## Key Patterns Used

- Flask route decorators with HTTP methods
- request.get_json(force=True) to parse JSON bodies
- jsonify() for JSON responses
- Response() for custom response types (CSV)
- Manual CSV creation without csv library
- Proper CSV escaping (quotes and commas)
- Error handling with status codes (400 for bad request)
- Content-Disposition header for file downloads
