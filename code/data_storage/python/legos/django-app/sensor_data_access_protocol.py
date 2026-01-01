# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

from abc import ABC, abstractmethod
from typing import List


class SensorDataAccess(ABC):
    """Abstract base class for sensor data access."""

    @abstractmethod
    def log_sensor_data(self, json_data: str) -> None:
        """Log sensor data."""
        pass

    @abstractmethod
    def fetch_sensor_data(self) -> List[str]:
        """Fetch sensor data."""
        pass

    @abstractmethod
    def purge_sensor_data(self) -> None:
        """Purge sensor data."""
        pass
