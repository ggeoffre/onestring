# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 6: python/legos/django-app
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Integrate Django API with multiple database implementations and helper utilities.

---

## File Structure
```
python/legos/django-app/
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
Same as Project 4 and 5:

Create abstract base class with ABC and three abstract methods.
```

---

## Prompt 2: Helper Module (sensor_data_helper.py)

```
Same as Project 1 and 5:

Include helper functions for:
- Generating random sensor data
- Converting to tuple for inserts
- JSON serialization handlers
- Converting list of dicts to CSV
```

---

## Prompt 3: Database Implementations

```
Use the SAME implementations as Project 5 (legos/flask-app):
- mongo_data_access.py
- cassandra_data_access.py
- mysql_data_access.py
- postgres_data_access.py
- redis_data_access.py

These are identical to the Flask version - they're framework-independent.
```

---

## Prompt 4: Django App with Data Access Integration (app.py)

```
Create app.py with Django routes integrated with data access layer:

1. Imports:
   - import json, os, sys
   - import sensor_data_helper
   - from django.conf import settings
   - from django.core.management import execute_from_command_line
   - from django.http import HttpResponse, JsonResponse
   - from django.urls import path, re_path
   - from mongo_data_access import MongoDataAccess
   - from mysql_data_access import MySQLDataAccess
   - from postgres_data_access import PostgresDataAccess
   - from redis_data_access import RedisDataAccess
   - from sensor_data_access_protocol import SensorDataAccess

2. Factory function get_data_access() -> SensorDataAccess:
   - Read DATA_ACCESS from environment (default "mongo")
   - If "redis": return RedisDataAccess()
   - If "mongo": return MongoDataAccess()
   - If "cassandra": return CassandraDataAccess()
   - If "mysql": return MySQLDataAccess()
   - If "postgres": return PostgresDataAccess()
   - Else: raise ValueError

3. Configure Django settings with settings.configure():
   - DEBUG=True
   - APPEND_SLASH=True
   - SECRET_KEY="your-secret-key"
   - ROOT_URLCONF=__name__
   - ALLOWED_HOSTS=["*"]
   - TEMPLATES with Django backend

4. View function home(request):
   - Return JsonResponse: {"message": "Django API Server is running!"}

5. View function echo_view(request):
   - If GET: return error 405: {"error": "GET requests are not allowed."}
   - If POST: parse request.body with json.loads, return as JsonResponse

6. View function log_view(request):
   - If GET: return error 405: {"error": "GET requests are not allowed."}
   - If POST:
     * Parse request.body with json.loads
     * Return 400 if None: {"error": "No valid JSON provided"}
     * Get data access instance
     * Call log_sensor_data(data)
     * Return success: {"message": "Data logged successfully"}

7. View function report_view(request):
   - Get data access instance
   - Call fetch_sensor_data()
   - Convert to CSV using sensor_data_helper.json_list_to_csv()
   - Create HttpResponse with CSV data, content_type="text/csv"
   - Add Content-Disposition header for download as "report.csv"
   - Return response

8. View function purge_view(request):
   - If GET or POST:
     * Get data access instance
     * Call purge_sensor_data()
     * Return: {"message": "Data purged."}

9. URL patterns (urlpatterns list):
   - path("", home, name="home")
   - re_path(r"^echo/?$", echo_view, name="echo")
   - re_path(r"^log/?$", log_view, name="log")
   - re_path(r"^report/?$", report_view, name="report")
   - re_path(r"^purge/?$", purge_view, name="purge")

10. Main block:
    - os.environ.setdefault("DJANGO_SETTINGS_MODULE", __name__)
    - execute_from_command_line(sys.argv)
```

---

## Prompt 5: Requirements File

```
Create requirements.txt with:
Django
pymongo
cassandra-driver
pymysql
psycopg2-binary
redis
```

---

## Running Instructions

```
To run:
python app.py runserver 0.0.0.0:8080
```

---

## Expected Behavior

This integrated Django system:
- Runs Django API
- Routes use factory pattern to get database implementation
- POST /log stores data in selected database
- GET /report retrieves and exports as CSV
- POST or GET /purge clears all data
- Database selection via DATA_ACCESS environment variable
- Each database creates schema on first use
- Routes work with or without trailing slash

---

## Key Differences from Flask Version

**Django-Specific:**
- settings.configure() for single-file app
- JsonResponse instead of jsonify
- request.body instead of request.get_json()
- HttpResponse instead of Response
- re_path for regex routes
- execute_from_command_line for runserver
- Request method checking (request.method == "POST")

**Similarities:**
- Same factory pattern
- Same database implementations
- Same helper functions
- Same overall architecture
- Same route behavior

---

## Key Integration Patterns

- Factory pattern for database selection
- Abstract protocol ensures all databases work identically
- Single-file Django app (no separate settings.py)
- Dependency injection (views call factory)
- CSV conversion utility
- Environment-based configuration
- Global state for schema creation
- Error handling with HTTP status codes (405, 400, 500)
- Request method validation

---

## Reusability Lesson

Notice that:
- All database implementations are IDENTICAL to Flask version
- Only app.py changes (Django vs Flask specifics)
- Helper functions are framework-independent
- Abstract protocol enables this portability

This demonstrates proper separation of concerns:
- Data access layer = framework-independent
- API layer = framework-specific
- Business logic (helpers) = framework-independent
