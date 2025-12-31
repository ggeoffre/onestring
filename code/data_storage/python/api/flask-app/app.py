# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

from flask import Flask

app = Flask(__name__)


@app.route("/", methods=["GET"])
def root():
    return "Flask API Server is running!"


@app.route("/echo", methods=["POST"])
def echo():
    return "Flask API Server is responding to /echo"


@app.route("/log", methods=["POST"])
def log():
    return "Flask API Server is responding to /log"


@app.route("/report", methods=["GET"])
def report():
    return "Flask API Server is responding to /report"


@app.route("/purge", methods=["GET", "POST"])
def purge():
    return "Flask API Server is responding to /purge"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
