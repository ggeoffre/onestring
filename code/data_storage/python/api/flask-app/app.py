# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json

from flask import Flask, Response, jsonify, request

# Standard JSON constant
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


def json_to_csv_safe(json_string):
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


app = Flask(__name__)


@app.route("/", methods=["GET"])
def root():
    return jsonify({"message": "Flask API Server is running!"})


@app.route("/echo", methods=["POST"])
@app.route("/log", methods=["POST"])
def echo_log():
    data = request.get_json(force=True)
    if data is None:
        return jsonify({"error": "No valid JSON provided"}), 400
    return jsonify(data)


@app.route("/report", methods=["GET"])
def report():
    csv_data = json_to_csv_safe(SENSOR_DATA_JSON)
    return Response(
        csv_data,
        mimetype="text/csv",
        headers={"Content-Disposition": "attachment; filename=sensor_report.csv"},
    )


@app.route("/purge", methods=["GET", "POST"])
def purge():
    return jsonify({"message": "Data purge sequence complete"})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080, debug=False)
