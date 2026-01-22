# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# PROJECT 4: python/access
## Reverse-Engineered Prompts from Working Code

**Project Goal**: Create abstract base class with multiple stub implementations demonstrating polymorphism.

---

## File Structure
```
python/access/
├── sensor_data_access_protocol.py
├── mongo_data_access.py
├── cassandra_data_access.py
├── mysql_data_access.py
├── postgres_data_access.py
├── redis_data_access.py
└── main.py
```

---

## Prompt 1: Abstract Base Class (sensor_data_access_protocol.py)

```
Create sensor_data_access_protocol.py with:

1. Imports:
   - from abc import ABC, abstractmethod
   - from typing import List

2. Class SensorDataAccess(ABC):
   - Docstring: "Abstract base class for sensor data access."
   
   - Abstract method log_sensor_data(self, json_data: str) -> None:
     * Docstring: "Log sensor data."
     * Just pass
   
   - Abstract method fetch_sensor_data(self) -> List[str]:
     * Docstring: "Fetch sensor data."
     * Just pass
   
   - Abstract method purge_sensor_data(self) -> None:
     * Docstring: "Purge sensor data."
     * Just pass
```

---

## Prompt 2: MongoDB Implementation (mongo_data_access.py)

```
Create mongo_data_access.py with:

1. Import: from sensor_data_access_protocol import SensorDataAccess

2. Class MongoDataAccess(SensorDataAccess):
   - Docstring: "MongoDB implementation of SensorDataAccess."
   
   - Method log_sensor_data(self, json_data: str) -> None:
     * Print: f"Logging sensor data to MongoDB: {json_data}"
   
   - Method fetch_sensor_data(self) -> list[str]:
     * Print: "Fetching sensor data from MongoDB"
     * Return: ['{"sensor": "humidity", "value": 45.6}']
   
   - Method purge_sensor_data(self) -> None:
     * Print: "Purging sensor data from MongoDB"
```

---

## Prompt 3: Cassandra Implementation (cassandra_data_access.py)

```
Create cassandra_data_access.py with:

1. Import: from sensor_data_access_protocol import SensorDataAccess

2. Class CassandraDataAccess(SensorDataAccess):
   - Docstring: "Cassandra implementation of SensorDataAccess."
   
   - Method log_sensor_data(self, json_data: str) -> None:
     * Print: f"Logging sensor data to Cassandra: {json_data}"
   
   - Method fetch_sensor_data(self) -> list[str]:
     * Print: "Fetching sensor data from Cassandra"
     * Return: ['{"sensor": "pressure", "value": 1013}']
   
   - Method purge_sensor_data(self) -> None:
     * Print: "Purging sensor data from Cassandra"
```

---

## Prompt 4: MySQL Implementation (mysql_data_access.py)

```
Create mysql_data_access.py with:

1. Import: from sensor_data_access_protocol import SensorDataAccess

2. Class MySQLDataAccess(SensorDataAccess):
   - Docstring: "MySQL implementation of SensorDataAccess."
   
   - Method log_sensor_data(self, json_data: str) -> None:
     * Print: f"Logging sensor data to MySQL: {json_data}"
   
   - Method fetch_sensor_data(self) -> list[str]:
     * Print: "Fetching sensor data from MySQL"
     * Return: ['{"sensor": "light", "value": 300}']
   
   - Method purge_sensor_data(self) -> None:
     * Print: "Purging sensor data from MySQL"
```

---

## Prompt 5: PostgreSQL Implementation (postgres_data_access.py)

```
Create postgres_data_access.py with:

1. Import: from sensor_data_access_protocol import SensorDataAccess

2. Class PostgresDataAccess(SensorDataAccess):
   - Docstring: "PostgreSQL implementation of SensorDataAccess."
   
   - Method log_sensor_data(self, json_data: str) -> None:
     * Print: f"Logging sensor data to PostgreSQL: {json_data}"
   
   - Method fetch_sensor_data(self) -> list[str]:
     * Print: "Fetching sensor data from PostgreSQL"
     * Return: ['{"sensor": "sound", "value": 75}']
   
   - Method purge_sensor_data(self) -> None:
     * Print: "Purging sensor data from PostgreSQL"
```

---

## Prompt 6: Redis Implementation (redis_data_access.py)

```
Create redis_data_access.py with:

1. Import: from sensor_data_access_protocol import SensorDataAccess

2. Class RedisDataAccess(SensorDataAccess):
   - Docstring: "Redis implementation of SensorDataAccess."
   
   - Method log_sensor_data(self, json_data: str) -> None:
     * Print: f"Logging sensor data to Redis: {json_data}"
   
   - Method fetch_sensor_data(self) -> list[str]:
     * Print: "Fetching sensor data from Redis"
     * Return: ['{"sensor": "temperature", "value": 22.3}']
   
   - Method purge_sensor_data(self) -> None:
     * Print: "Purging sensor data from Redis"
```

---

## Prompt 7: Factory and Main (main.py)

```
Create main.py with:

1. Imports:
   - import os
   - from cassandra_data_access import CassandraDataAccess
   - from mongo_data_access import MongoDataAccess
   - from mysql_data_access import MySQLDataAccess
   - from postgres_data_access import PostgresDataAccess
   - from redis_data_access import RedisDataAccess
   - from sensor_data_access_protocol import SensorDataAccess

2. Function get_data_access() -> SensorDataAccess:
   - Get DATA_ACCESS from environment (default "mongo")
   - If "redis": return RedisDataAccess()
   - If "mongo": return MongoDataAccess()
   - If "cassandra": return CassandraDataAccess()
   - If "mysql": return MySQLDataAccess()
   - If "postgres": return PostgresDataAccess()
   - Else: raise ValueError(f"Unsupported DATA_ACCESS type: {data_access_type}")

3. Function main():
   - Call get_data_access() to get instance
   - Call log_sensor_data with: '{"sensor": "temperature", "value": 22.3}'
   - Call fetch_sensor_data and print result: "Fetched sensor data:", data
   - Call purge_sensor_data()

4. Main guard:
   - if __name__ == "__main__": call main()
```

---

## Expected Behavior

When run:
- Reads DATA_ACCESS environment variable (defaults to "mongo")
- Creates appropriate data access implementation
- Calls the three methods (log, fetch, purge)
- Each method prints what it's doing and returns stub data

Example output (with DATA_ACCESS=mongo):
```
Logging sensor data to MongoDB: {"sensor": "temperature", "value": 22.3}
Fetching sensor data from MongoDB
Fetched sensor data: ['{"sensor": "humidity", "value": 45.6}']
Purging sensor data from MongoDB
```

---

## Key Patterns Used

- Abstract Base Class (ABC) pattern
- @abstractmethod decorator
- Type hints (List[str], -> None)
- Factory pattern (get_data_access function)
- Environment-based configuration
- Polymorphism (all implementations conform to same interface)
- Strategy pattern (swap implementations via environment variable)
- Stub implementations (print instead of real database operations)
