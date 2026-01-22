# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 3: ACCESS - Better Quality Prompts
## Focus: Clean Architecture, Thread-Safe, Production-Ready Data Access

---

## Quality-Focused Core Prompts

### Prompt 1: Define Interface with Type Safety
```
Create sensor_data_access_protocol.py with an abstract base class:
- Use Python's ABC module for the interface
- Define clear method signatures with type hints for all parameters and returns
- Include docstrings explaining what each method does
- Methods should be:
  * log(data: dict) -> str  # Returns ID of inserted record
  * find_all(limit: int = None, offset: int = 0) -> list[dict]
  * find_by_id(record_id: str) -> dict | None
  * delete_all() -> int  # Returns count deleted
  * close() -> None  # Cleanup resources
- Add abstract property: database_type -> str
- Add abstract method: health_check() -> bool

Why: Strong typing catches errors early and makes the interface self-documenting.
```

### Prompt 2: Thread-Safe Connection Management
```
For all database implementations, use this pattern:
- Create a connection pool at class initialization (not on every request)
- Store the pool as an instance variable (not global or class variable)
- Use thread-local storage for connections
- Implement proper cleanup in close() method
- Never use global variables for state
- Use context managers for connection checkout/checkin

Example pattern:
- Initialize pool in __init__
- Get connection from pool for each operation
- Return connection to pool when done
- Clean up entire pool in close()

Why: Thread-safe code prevents race conditions and enables concurrent requests.
```

### Prompt 3: MongoDB Implementation with Connection Pooling
```
Create MongoDataAccess that implements the interface:
- Initialize MongoClient with connection pooling in __init__:
  * maxPoolSize=10
  * minPoolSize=2
  * maxIdleTimeMS=60000
  * connectTimeoutMS=5000
- Store pool as instance variable
- Implement all interface methods
- Use projection to fetch only needed fields
- Return consistent types (dict, list[dict], int, str)
- Handle MongoDB ObjectId properly (convert to string for IDs)
- Never print to console (use logging module)
- Close pool in close() method

Why: Connection pooling dramatically improves performance and resource usage.
```

### Prompt 4: Configuration Management
```
Create database_config.py with dataclasses:
- Define a dataclass for each database type's configuration
- Load all settings from environment variables
- Validate required settings in __post_init__
- Provide clear error messages for missing/invalid config
- Use pydantic if available for validation, else plain dataclasses
- Never hardcode credentials
- Support loading from .env file for development

Example classes:
- MongoConfig: host, port, database, username, password, pool_size
- RedisConfig: host, port, password, db_number, pool_size
- PostgresConfig: host, port, database, user, password, pool_size

Why: Centralized configuration is easier to manage and validate.
```

### Prompt 5: Error Handling with Custom Exceptions
```
Create exceptions.py with custom exception hierarchy:
- DatabaseConnectionError (when can't connect)
- DatabaseOperationError (when operation fails)
- ValidationError (when data is invalid)
- RecordNotFoundError (when find_by_id returns nothing)

All exceptions should:
- Inherit from a base DataAccessException
- Include original error in __cause__
- Have clear error messages
- Support additional context (operation, record_id, etc.)

Update all implementations to raise these specific exceptions instead of generic Exception.

Why: Specific exceptions make error handling more precise and informative.
```

### Prompt 6: Repository Pattern with Business Logic Separation
```
Create SensorDataRepository class that wraps the data access layer:
- Takes a SensorDataAccess implementation in __init__
- Adds business logic on top of data access:
  * validate_before_save(data: dict) -> validates and sanitizes
  * get_recent_readings(hours: int) -> filters by timestamp
  * get_by_location(location: str) -> filters by location
  * aggregate_by_sensor() -> groups data by sensor type
- Implements context manager (__enter__, __exit__)
- Handles all exceptions from data access layer
- Returns domain objects (not raw dicts)
- Never leaks database-specific details

Why: Separates business logic from data access, making both easier to test and maintain.
```

### Prompt 7: Factory with Dependency Injection
```
Create database_factory.py with a factory function:
- get_data_access(config: DatabaseConfig) -> SensorDataAccess
- Reads database type from config
- Creates appropriate implementation
- Configures connection pool based on config
- Registers cleanup handlers
- Validates connection before returning
- Raises clear error for unsupported database types

Create a dependency injection container:
- Singleton pattern for database connections
- Lazy initialization (create on first use)
- Thread-safe singleton implementation
- Support for testing (ability to inject mocks)

Why: Makes configuration-driven, makes testing easier, centralizes object creation.
```

---

## Quality Enhancement Prompts

### Prompt 8: Retry Logic with Exponential Backoff
```
Create a retry decorator for transient failures:
- Retry on connection errors only (not validation errors)
- Maximum 3 attempts
- Exponential backoff: 1s, 2s, 4s
- Log each retry attempt
- Re-raise original exception after max attempts
- Skip retries for write operations (only retry reads)
- Use tenacity library if available, else implement simple decorator

Apply to all find operations, not write operations.

Why: Handles temporary network issues gracefully without complicating business logic.
```

### Prompt 9: Prepared Statements for SQL Databases
```
For PostgreSQL and MySQL implementations:
- Never use string formatting for queries (SQL injection risk)
- Always use parameterized queries
- Prepare frequently-used queries once at initialization
- Cache prepared statements in instance variable
- Use proper parameter binding
- Example: cursor.execute("SELECT * FROM table WHERE id = %s", (id,))

Why: Prevents SQL injection and improves performance through query caching.
```

### Prompt 10: Monitoring and Metrics
```
Add performance monitoring to all implementations:
- Track operation counts (reads, writes, deletes)
- Track operation latencies (min, max, avg, p95, p99)
- Track connection pool stats (active, idle, waiting)
- Track error counts by type
- Expose metrics via get_metrics() method
- Log slow operations (>100ms) with WARNING
- Reset metrics on demand

Use a simple dataclass to store metrics, not a heavy framework.

Why: Observability is essential for production systems.
```

### Prompt 11: Transaction Support (where applicable)
```
For databases that support transactions (PostgreSQL, MySQL):
- Add begin_transaction() method
- Add commit() method
- Add rollback() method
- Support context manager for automatic rollback on error
- Example:
  with db_access.transaction():
      db_access.log(data1)
      db_access.log(data2)
      # auto-commits if no error, auto-rollbacks on error

For databases without transactions (MongoDB, Redis):
- Document that transactions are not supported
- Suggest alternative patterns (idempotency, compensation)

Why: Ensures data consistency for multi-step operations.
```

### Prompt 12: Caching Layer
```
Add optional caching to reduce database load:
- Use Redis or in-memory cache (configurable)
- Cache read operations (find_all, find_by_id)
- Set TTL on cache entries (60 seconds default)
- Invalidate cache on writes/deletes
- Cache key pattern: f"{table}:{operation}:{params}"
- Make caching opt-in via configuration
- Track cache hit rate in metrics

Why: Dramatically reduces database load for read-heavy workloads.
```

### Prompt 13: Bulk Operations for Performance
```
Add bulk operation methods:
- log_batch(data_list: list[dict]) -> list[str]  # Returns list of IDs
- delete_by_ids(ids: list[str]) -> int  # Returns count deleted
- Validate all items before inserting any (atomic)
- Use database-specific bulk operations (insert_many, executemany)
- 10-100x faster than individual operations
- Include progress callback for large batches

Why: Bulk operations are essential for high-throughput scenarios.
```

### Prompt 14: Health Check Implementation
```
Implement health_check() for each database:
- Attempt simple operation (SELECT 1 or ping)
- Set timeout to 5 seconds
- Return True if successful, False if failed
- Log failures but don't raise exceptions
- Include health check in metrics
- Check connection pool health (any connections available?)

Why: Enables monitoring tools to detect database issues.
```

---

## Testing & Quality Assurance Prompts

### Prompt 15: Unit Tests with Mocks
```
Create test suite for the interface implementation:
- Use pytest with fixtures
- Mock database connections (don't use real database)
- Test all methods with valid inputs
- Test all methods with invalid inputs
- Test connection failures
- Test timeout scenarios
- Test thread safety (concurrent operations)
- Test resource cleanup (connections closed)
- Achieve >90% code coverage

Why: Tests ensure correctness and prevent regressions.
```

### Prompt 16: Integration Tests
```
Create integration test suite:
- Use testcontainers to spin up real databases
- Test against actual MongoDB, PostgreSQL, etc.
- Test connection pooling under load
- Test transaction rollback
- Test bulk operations with large datasets
- Test retry logic with network interruption
- Clean up test data after each test

Why: Integration tests catch issues that unit tests miss.
```

### Prompt 17: Load Testing
```
Create load test script:
- Simulate 100 concurrent threads
- Each thread performs 1000 operations
- Mix of reads (70%), writes (25%), deletes (5%)
- Measure throughput (operations per second)
- Measure latency (p50, p95, p99)
- Monitor connection pool usage
- Ensure no connection leaks
- Fail if error rate > 0.1%

Why: Ensures the implementation can handle production load.
```

---

## Expected Quality Improvements

### Architecture Quality Metrics:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Thread Safety** | None (global state) | Excellent (thread-local) | ⬆️ 100% |
| **Resource Management** | Poor (leaks) | Excellent (pooling) | ⬆️ 95% |
| **Performance** | Poor (new conn each time) | Excellent (pooling) | ⬆️ 500% |
| **Error Handling** | Poor (generic exceptions) | Good (specific) | ⬆️ 85% |
| **Testability** | Poor (hard to mock) | Excellent (DI) | ⬆️ 95% |
| **Maintainability** | Fair (some duplication) | Excellent (DRY) | ⬆️ 80% |
| **Observability** | Poor (no metrics) | Good (detailed metrics) | ⬆️ 90% |
| **Security** | Poor (hardcoded creds) | Excellent (config-driven) | ⬆️ 95% |

---

## Key Architecture Principles Applied

### 1. **Dependency Inversion**
- Depend on abstractions (interface), not concrete implementations
- Factory pattern for object creation
- Dependency injection for flexibility

### 2. **Single Responsibility**
- Data access layer only handles database operations
- Repository layer adds business logic
- Configuration management is separate
- Each class has one reason to change

### 3. **Open/Closed**
- Easy to add new database implementations
- Easy to add new features without modifying existing code
- Extension points clearly defined

### 4. **Resource Management**
- Connection pooling for efficiency
- Proper cleanup in all cases
- Context managers for automatic resource management

### 5. **Fail Fast**
- Validate configuration at startup
- Check connections before using
- Raise specific exceptions for different failures

---

## Why This Produces Better Code

### Performance:
- Connection pooling reduces latency by 80-90%
- Bulk operations are 10-100x faster
- Caching reduces database load by 60-80%
- Prepared statements improve SQL performance

### Reliability:
- Thread safety prevents race conditions
- Retry logic handles transient failures
- Proper error handling prevents crashes
- Resource cleanup prevents leaks

### Security:
- Parameterized queries prevent SQL injection
- No hardcoded credentials
- Configuration validation prevents misconfiguration
- Least privilege principle applied

### Maintainability:
- Clean separation of concerns
- Single implementation pattern across databases
- Dependency injection makes testing easy
- Metrics show what's actually happening

### Scalability:
- Connection pooling supports high concurrency
- Thread safety enables parallel processing
- Bulk operations handle high throughput
- No global state to serialize access

---

## Production Readiness Checklist

After these prompts, the data access layer will have:
- ✅ Thread-safe connection management
- ✅ Connection pooling for performance
- ✅ Proper error handling with custom exceptions
- ✅ Retry logic for transient failures
- ✅ Transaction support where applicable
- ✅ Monitoring and metrics
- ✅ Health checks for monitoring
- ✅ Comprehensive test coverage
- ✅ Security best practices
- ✅ Documentation and type hints

---

## Balance Maintained

These prompts:
- ✅ Use clear, simple language
- ✅ Explain the WHY behind each pattern
- ✅ Focus on principles, not syntax
- ✅ Build incrementally (core → enhancements)
- ✅ Provide concrete examples when helpful
- ✅ Address real production concerns

Developers understand what production-grade architecture looks like without being given a textbook.
