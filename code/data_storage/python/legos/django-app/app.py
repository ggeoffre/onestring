# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json
import os
import sys

import sensor_data_helper
from django.conf import settings
from django.core.management import execute_from_command_line
from django.http import HttpResponse, JsonResponse
from django.urls import path, re_path
from mongo_data_access import MongoDataAccess
from sensor_data_access_protocol import SensorDataAccess


def get_data_access() -> SensorDataAccess:
    return MongoDataAccess()


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
        if data is None:
            return JsonResponse({"error": "No valid JSON provided"}), 400
        data_access = get_data_access()
        data_access.log_sensor_data(data)
        return JsonResponse({"message": "Data logged successfully"})


def report_view(request):
    data_access = get_data_access()
    data = data_access.fetch_sensor_data()
    csv_data = sensor_data_helper.json_list_to_csv(data)
    response = HttpResponse(csv_data, content_type="text/csv")
    response["Content-Disposition"] = 'attachment; filename="report.csv"'
    return response


def purge_view(request):
    if request.method in ["GET", "POST"]:
        data_access = get_data_access()
        data_access.purge_sensor_data()
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
