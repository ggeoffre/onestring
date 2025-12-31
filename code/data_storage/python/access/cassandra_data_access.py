# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

from sensor_data_access_protocol import SensorDataAccess


class CassandraDataAccess(SensorDataAccess):
    """Cassandra implementation of SensorDataAccess."""

    def log_sensor_data(self, json_data: str) -> None:
        print(f"Logging sensor data to Cassandra: {json_data}")

    def fetch_sensor_data(self) -> list[str]:
        print("Fetching sensor data from Cassandra")
        return ['{"sensor": "pressure", "value": 1013}']

    def purge_sensor_data(self) -> None:
        print("Purging sensor data from Cassandra")
