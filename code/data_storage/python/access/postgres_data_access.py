# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

from sensor_data_access_protocol import SensorDataAccess


class PostgresDataAccess(SensorDataAccess):
    """PostgreSQL implementation of SensorDataAccess."""

    def log_sensor_data(self, json_data: str) -> None:
        print(f"Logging sensor data to PostgreSQL: {json_data}")

    def fetch_sensor_data(self) -> list[str]:
        print("Fetching sensor data from PostgreSQL")
        return ['{"sensor": "sound", "value": 75}']

    def purge_sensor_data(self) -> None:
        print("Purging sensor data from PostgreSQL")
