# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 3: python/api/django-app
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Create a standalone Django API without separate settings file.

---

## File Structure
```
python/api/django-app/
├── app.py
└── requirements.txt
```

---

## Prompt 1: Django API Application (app.py)

```
Create a Django application in a single file app.py with:

1. Imports:
   - import json, os, sys
   - from django.conf import settings
   - from django.core.management import execute_from_command_line
   - from django.http import HttpResponse, JsonResponse
   - from django.urls import path, re_path

2. Constant SENSOR_DATA_JSON:
   - Dictionary with: recorded (1768570200), location ("den"), sensor ("bmp280"), 
     measurement ("temperature"), units ("C"), value (22.3)
   - Convert to JSON string using json.dumps()

3. Helper function json_to_csv(json_string):
   - Parse JSON string
   - If dict, convert to list with one dict
   - Return empty string if empty
   - Get headers from first dict keys
   - Build CSV with headers in first line
   - For each item:
     * Convert values to strings
     * Escape quotes: replace " with ""
     * Wrap in quotes if contains comma or quote
   - Return newline-joined CSV
   - Return empty string on exception

4. Configure Django settings with settings.configure():
   - DEBUG=True
   - APPEND_SLASH=True
   - SECRET_KEY="your-secret-key"
   - ROOT_URLCONF=__name__
   - ALLOWED_HOSTS=["*"]
   - TEMPLATES with Django backend and empty DIRS

5. View function home(request):
   - Return JsonResponse: {"message": "Django API Server is running!"}

6. View function echo_view(request):
   - If GET: return error 405: {"error": "GET requests are not allowed."}
   - If POST: parse request.body with json.loads, return as JsonResponse

7. View function log_view(request):
   - If GET: return error 405: {"error": "GET requests are not allowed."}
   - If POST: parse request.body with json.loads, return as JsonResponse

8. View function report_view(request):
   - Convert SENSOR_DATA_JSON to CSV using helper
   - Create HttpResponse with CSV data, content_type="text/csv"
   - Add header Content-Disposition: 'attachment; filename="report.csv"'
   - Return response

9. View function purge_view(request):
   - If GET or POST: return JsonResponse: {"message": "Data purged."}

10. URL patterns (urlpatterns list):
    - path("", home, name="home")
    - re_path(r"^echo/?$", echo_view, name="echo")
    - re_path(r"^log/?$", log_view, name="log")
    - re_path(r"^report/?$", report_view, name="report")
    - re_path(r"^purge/?$", purge_view, name="purge")

11. Main block:
    - os.environ.setdefault("DJANGO_SETTINGS_MODULE", __name__)
    - execute_from_command_line(sys.argv)
```

---

## Prompt 2: Requirements File

```
Create requirements.txt with:
Django
```

---

## Running Instructions

```
To run:
python app.py runserver 0.0.0.0:8080
```

---

## Expected Behavior

- GET / returns welcome message
- POST /echo and /log return the JSON sent (GET returns 405 error)
- GET /report downloads CSV file with sensor data
- GET or POST /purge returns success message
- Routes work with or without trailing slash (APPEND_SLASH=True)

---

## Key Patterns Used

- Django settings.configure() for single-file app
- ROOT_URLCONF=__name__ to use current module for routing
- JsonResponse for JSON responses
- HttpResponse for custom response types (CSV)
- path() for exact matches, re_path() for regex patterns
- Request method checking (request.method)
- execute_from_command_line() to handle Django commands
- Manual CSV generation without csv library
- Proper CSV escaping for quotes and commas
