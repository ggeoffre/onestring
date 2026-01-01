# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

import json

import sensor_data_helper
from cassandra_data_access import CassandraDataAccess
from flask import Flask, Response, jsonify, request
from sensor_data_access_protocol import SensorDataAccess


def get_data_access() -> SensorDataAccess:
    return CassandraDataAccess()


app = Flask(__name__)


@app.route("/", methods=["GET"])
def root():
    return jsonify({"message": "Flask API Server is running!"})


@app.route("/echo", methods=["POST"])
def echo():
    data = request.get_json(force=True)
    if data is None:
        return jsonify({"error": "No valid JSON provided"}), 400
    return jsonify(data)


@app.route("/log", methods=["POST"])
def log():
    data_access = get_data_access()
    data = request.get_json(force=True)
    print(data)
    data_access.log_sensor_data(data)
    if data is None:
        return jsonify({"error": "No valid JSON provided"}), 400
    return jsonify({"message": "Data logged successfully"})


@app.route("/report", methods=["GET"])
def report():
    try:
        data_access = get_data_access()
        data = data_access.fetch_sensor_data()
        csv_data = sensor_data_helper.json_list_to_csv(data)
        if not csv_data:
            return jsonify({"error": "No data available"}), 404
        return Response(
            csv_data,
            mimetype="text/csv",
            headers={"Content-Disposition": "attachment; filename=sensor_report.csv"},
        )
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/purge", methods=["GET", "POST"])
def purge():
    data_access = get_data_access()
    data_access.purge_sensor_data()
    return jsonify({"message": "Data purge sequence complete"})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080, debug=False)
