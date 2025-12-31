# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json
import os
import sys

from django.conf import settings
from django.core.management import execute_from_command_line
from django.http import HttpResponse, JsonResponse
from django.urls import path, re_path

SENSOR_DATA_JSON = json.dumps(
    {
        "recorded": 1768570200,
        "location": "den",
        "sensor": "bmp280",
        "measurement": "temperature",
        "units": "C",
        "value": 22.3,
    }
)


def json_to_csv(json_string):
    """Converts JSON to CSV with manual escaping for reliability."""
    try:
        data = json.loads(json_string)
        if isinstance(data, dict):
            data = [data]
        if not data:
            return ""

        headers = list(data[0].keys())
        csv_lines = [",".join(headers)]

        for item in data:
            row = []
            for key in headers:
                val = str(item.get(key, ""))
                # 2025 Safety: Escape quotes and handle commas
                val = val.replace('"', '""')
                if "," in val or '"' in val:
                    val = f'"{val}"'
                row.append(val)
            csv_lines.append(",".join(row))
        return "\n".join(csv_lines)
    except Exception:
        return ""


# Configure Django settings
settings.configure(
    DEBUG=True,
    APPEND_SLASH=True,
    SECRET_KEY="your-secret-key",
    ROOT_URLCONF=__name__,
    ALLOWED_HOSTS=["*"],
    TEMPLATES=[
        {
            "BACKEND": "django.template.backends.django.DjangoTemplates",
            "DIRS": [],
        },
    ],
)


# Views
def home(request):
    return JsonResponse({"message": "Django API Server is running!"})


def echo_view(request):
    if request.method == "GET":
        return JsonResponse({"error": "GET requests are not allowed."}, status=405)
    elif request.method == "POST":
        data = json.loads(request.body)
        return JsonResponse(data)


def log_view(request):
    if request.method == "GET":
        return JsonResponse({"error": "GET requests are not allowed."}, status=405)
    elif request.method == "POST":
        data = json.loads(request.body)
        return JsonResponse(data)


def report_view(request):
    csv_data = json_to_csv(SENSOR_DATA_JSON)
    response = HttpResponse(csv_data, content_type="text/csv")
    response["Content-Disposition"] = 'attachment; filename="report.csv"'
    return response


def purge_view(request):
    if request.method == "GET":
        return JsonResponse({"message": "Data purged."})
    elif request.method == "POST":
        return JsonResponse({"message": "Data purged."})


# URL patterns
urlpatterns = [
    path("", home, name="home"),
    re_path(r"^echo/?$", echo_view, name="echo"),
    re_path(r"^log/?$", log_view, name="log"),
    re_path(r"^report/?$", report_view, name="report"),
    re_path(r"^purge/?$", purge_view, name="purge"),
]

# Run the app
if __name__ == "__main__":
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", __name__)
    execute_from_command_line(sys.argv)
