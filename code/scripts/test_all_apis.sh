#!/bin/bash

# API Smoke Test Script
#
# This script is designed to [insert a brief summary of what the script does,
# including its main functionality and purpose]. It aims to [describe the
# intended outcome or effect of running the script]. The script utilizes
# [mention any key technologies, libraries, or methods used] to achieve its
# objectives, ensuring [mention any important aspects such as efficiency,
# reliability, or user-friendliness].
#
# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC
#

usage() {
  cat <<EOF
Usage: $0 [DIRECTORY]
Run a simple API smoke test against http://localhost:8080 using either wget or curl.

ARGUMENT:
    DIRECTORY   Name of the HTTP client to use (case-insensitive): "wget" or "curl".

BEHAVIOR:
    - Generates a JSON measurement with a random temperature value (22.1–32.4 C) and the current Unix timestamp.
    - Prints the JSON and the chosen tester name.
    - Executes the following sequence of requests against localhost:8080:
            GET  /           (log)
            POST /echo       (Content-Type: application/json) with the generated JSON
            POST /log        (Content-Type: application/json) with the generated JSON
            GET  /report
            GET  /purge
            POST /purge      (empty POST)

EXAMPLES:
    $0 wget
    $0 curl
    $0 --help

SAMPLE OUTPUT:
    {"recorded":"1768570200","location":"den","sensor":"bmp280","measurement":"temperature","units":"C","value":25.3}
    API Tester: WGET
    LOG
    ... (output from wget) ...

If no argument or an unrecognized value is given the script will report "Unknown API tester".
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

JSON_VALUE=$(awk -v min=22.1 -v max=32.4 'BEGIN{srand(); printf "%.1f\n", min+rand()*(max-min)}');
JSON_RECORDED=$(date +%s);
printf -v JSON_STRING '{"recorded":"%s","location":"den","sensor":"bmp280","measurement":"temperature","units":"C","value":%s}' "$JSON_RECORDED" "$JSON_VALUE";
echo $JSON_STRING; echo

API_TESTER=$(printf '%s' "${1:-}" | tr '[:lower:]' '[:upper:]')

echo "API Tester: $API_TESTER"

HOST="${HOST:-localhost}"
PORT="${PORT:-8080}"
BASE_URL="http://$HOST:$PORT"

case "$API_TESTER" in
    WGET)
        echo "ROOT"
        wget -q -O - "$BASE_URL/"; echo
        echo "ECHO"
        wget -q -O - --header='Content-Type: application/json' --post-data="$JSON_STRING" "$BASE_URL/echo"; echo
        echo "LOG"
        wget -q -O - --header='Content-Type: application/json' --post-data="$JSON_STRING" "$BASE_URL/log"; echo
        echo "REPORT"
        wget -q -O - "$BASE_URL/report"; echo
        echo "PURGE"
        wget -q -O - "$BASE_URL/purge"; echo
        wget -q -O - --post-data='' "$BASE_URL/purge"; echo
        ;;
    CURL)
        echo "ROOT"
        curl "$BASE_URL/";echo
        echo "ECHO"
        curl -X POST -H "Content-Type: application/json" -d "$JSON_STRING" "$BASE_URL/echo";echo
        echo "LOG"
        curl -X POST -H "Content-Type: application/json" -d "$JSON_STRING" "$BASE_URL/log";echo
        echo "REPORT"
        curl "$BASE_URL/report" -w "\n"
        echo "PURGE"
        curl "$BASE_URL/purge"; echo
        curl -X POST -d '' "$BASE_URL/purge"; echo
        ;;
    *)
        echo "Unknown API tester: $API_TESTER"
        ;;
esac
