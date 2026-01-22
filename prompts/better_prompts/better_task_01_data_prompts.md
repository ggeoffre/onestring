# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 1: DATA - Better Quality Prompts
## Focus: Cleaner, More Efficient, Maintainable, Secure Code

---

## Quality-Focused Core Prompts

### Prompt 1: Connection Management with Resource Cleanup
```
Create a Python function for MongoDB connection that:
- Uses a context manager (with statement) to automatically close connections
- Returns client, database, and collection
- Gets hostname from DATA_HOSTNAME environment variable (default "localhost")
- Uses port 27017 and connects to "sensor_data_db" database, "sensor_data" collection
- Has a 5-second connection timeout
- Closes the connection when done (even if an error occurs)
- Example pattern: Use 'try/finally' or context manager protocol

Why: Prevents resource leaks and ensures connections are always cleaned up properly.
```

### Prompt 2: Secure Configuration Management
```
Create a config.py file with a class called DatabaseConfig that:
- Reads all database settings from environment variables
- Provides sensible defaults
- Validates that required settings are present
- Raises clear errors if critical settings are missing
- Never hardcodes passwords or credentials in the code
- Includes settings for: hostname, port, database name, collection name, timeout
- Uses a @property decorator for read-only access

Why: Separates configuration from code, makes settings easier to change, and improves security.
```

### Prompt 3: Store Function with Validation
```
Create a function called store_sensor_data() that:
- Takes a dictionary as input
- Validates the dictionary has all required fields: recorded, location, sensor, measurement, units, value
- Validates that 'recorded' is an integer and 'value' is a number
- Returns True if successful, False if validation fails
- Uses the connection context manager from Prompt 1
- Logs errors using Python's logging module (not print statements)
- Only inserts if validation passes
- Handles MongoDB errors gracefully and logs them

Why: Prevents invalid data from entering the database and provides better error visibility.
```

### Prompt 4: Efficient Retrieval with Projection
```
Create a function called retrieve_sensor_data() that:
- Takes an optional 'limit' parameter (default None for all records)
- Takes an optional 'fields' parameter to specify which fields to return
- Uses MongoDB projection to only fetch requested fields (don't use SELECT *)
- Automatically excludes the '_id' field
- Returns an empty list if no data (not None)
- Uses the connection context manager
- Logs how many records were retrieved using Python's logging module
- Handles errors and logs them, returning empty list on error

Why: Fetches only needed data (more efficient), supports pagination, consistent return types.
```

### Prompt 5: Safe Delete with Confirmation
```
Create a function called delete_all_sensor_data() that:
- Takes an optional 'confirm' parameter (default False)
- Only deletes if confirm=True (safety feature)
- Returns the count of deleted documents
- Uses the connection context manager
- Logs the deletion count using Python's logging module
- Handles errors and logs them, returning 0 on error

Why: Prevents accidental data loss, provides feedback on what was deleted.
```

### Prompt 6: Logging Configuration
```
Create a logging_config.py file that:
- Sets up Python's logging module
- Logs to both console and a file called 'sensor_app.log'
- Uses different log levels: INFO for normal operations, ERROR for failures, DEBUG for detailed info
- Formats logs with timestamp, level, and message
- Provides a function 'get_logger(name)' that other modules can use

Why: Professional logging is easier to debug and monitor than print statements.
```

### Prompt 7: Helper Module with Pure Functions
```
Create sensor_data_helper.py with these functions that have NO side effects:

1. generate_random_sensor_data() -> dict
   - Pure function that returns a new dictionary each time
   - Never modifies global state
   - Uses current timestamp and random temperature (22.4-32.1)

2. validate_sensor_data(data: dict) -> tuple[bool, str]
   - Returns (True, "") if valid, (False, error_message) if invalid
   - Checks all required fields exist
   - Checks data types are correct
   - Never raises exceptions, just returns validation result

3. sanitize_sensor_data(data: dict) -> dict
   - Returns a cleaned copy of the data
   - Removes any unexpected fields
   - Ensures types are correct (converts strings to numbers if needed)
   - Returns None if data cannot be sanitized

Why: Pure functions are easier to test, more reliable, and don't cause unexpected side effects.
```

### Prompt 8: Main with Proper Error Handling
```
Create main.py that:
- Configures logging first (using logging_config)
- Wraps all operations in try/except blocks
- Catches specific exceptions (not bare 'except Exception')
- Logs errors with context (what operation was being performed)
- Uses sys.exit(1) for failures (returns proper exit code)
- Demonstrates all operations: generate data, validate, store, retrieve, delete
- Never prints to console (use logging instead)

Why: Proper error handling and logging makes debugging much easier.
```

---

## Quality Enhancement Prompts

### Prompt 9: Connection Pooling
```
Modify the MongoDB connection to use connection pooling:
- Create a connection pool when the app starts (not on every request)
- Reuse connections from the pool
- Set maxPoolSize to 10 connections
- Set minPoolSize to 2 connections
- This makes the app much faster under load

Why: Connection pooling significantly improves performance by reusing connections.
```

### Prompt 10: Batch Operations
```
Add a function called store_sensor_data_batch() that:
- Takes a list of sensor data dictionaries
- Validates all of them first (fail fast if any are invalid)
- Uses MongoDB's insert_many() to insert all at once
- Returns a tuple: (success_count, failure_count, error_messages)
- Much faster than inserting one at a time

Why: Batch operations are 10-100x faster than individual inserts.
```

### Prompt 11: Data Access Layer Pattern
```
Create a class called SensorDataRepository that:
- Encapsulates all database operations in one class
- Uses the connection pool (from Prompt 9)
- Implements __enter__ and __exit__ for context manager support
- Provides methods: save(), find_all(), find_by_location(), delete_all()
- Handles all error logging internally
- Never leaks database-specific details to callers

Why: Encapsulation makes code more maintainable and easier to test.
```

### Prompt 12: Retry Logic for Reliability
```
Add a decorator called @retry_on_failure that:
- Retries a function up to 3 times if it fails
- Waits 1 second between retries (exponential backoff)
- Only retries on connection errors (not on validation errors)
- Logs each retry attempt
- Apply this decorator to all database operations

Why: Handles transient network issues automatically, making the app more reliable.
```

### Prompt 13: Input Sanitization for Security
```
Update the validation function to also sanitize inputs:
- Strip whitespace from string fields
- Ensure 'recorded' timestamp is within reasonable range (not in future, not before 2020)
- Ensure 'value' is within reasonable range (e.g., -50 to 100 for temperature)
- Escape any special characters in location/sensor names
- Reject data with unexpected fields (security hardening)

Why: Prevents injection attacks and ensures data integrity.
```

### Prompt 14: Performance Monitoring
```
Add a decorator called @measure_time that:
- Records how long each database operation takes
- Logs slow operations (> 100ms) with a WARNING
- Logs very slow operations (> 1000ms) with ERROR
- Apply to all database functions
- Helps identify performance bottlenecks

Why: Makes it easy to find and fix slow operations.
```

---

## Testing & Documentation Prompts

### Prompt 15: Unit Tests
```
Create test_sensor_data.py that:
- Uses pytest for testing
- Mocks the MongoDB connection (don't need real database for tests)
- Tests validation logic thoroughly (valid cases and invalid cases)
- Tests error handling (what happens when database is down)
- Tests the retry decorator
- Uses fixtures for common test data
- Achieves at least 80% code coverage

Why: Tests catch bugs early and make refactoring safer.
```

### Prompt 16: Documentation
```
Add docstrings to all functions using Google style:
- One-line summary of what the function does
- Args: parameter descriptions with types
- Returns: what the function returns
- Raises: what exceptions might be raised
- Example usage when helpful
- Keep it concise but complete

Why: Good documentation helps other developers (and your future self) understand the code.
```

---

## Expected Quality Improvements

### Code Quality Metrics:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Resource Management** | Poor (leaks) | Excellent (auto-cleanup) | ⬆️ 90% |
| **Error Handling** | Poor (bare except) | Good (specific exceptions) | ⬆️ 80% |
| **Security** | Poor (hardcoded creds) | Good (env vars, validation) | ⬆️ 85% |
| **Performance** | Poor (no pooling) | Good (pooling + batch) | ⬆️ 500% |
| **Maintainability** | Fair (duplicated code) | Good (DRY, encapsulated) | ⬆️ 70% |
| **Testability** | Poor (hard to test) | Excellent (mockable) | ⬆️ 90% |
| **Observability** | Poor (print statements) | Good (structured logging) | ⬆️ 85% |

---

## Key Quality Principles Applied

### 1. **SOLID Principles**
- Single Responsibility: Each function does one thing
- Open/Closed: Easy to extend without modifying existing code
- Dependency Inversion: Depend on abstractions, not concrete implementations

### 2. **DRY (Don't Repeat Yourself)**
- Shared configuration in one place
- Reusable decorators for cross-cutting concerns
- Common validation logic centralized

### 3. **Fail Fast**
- Validate inputs early
- Check configuration at startup
- Don't hide errors

### 4. **Defensive Programming**
- Validate all inputs
- Handle all error cases
- Never trust external data
- Use type hints

### 5. **Clean Code**
- Descriptive names
- Small, focused functions
- Consistent style
- Minimal complexity

---

## Why This Produces Better Code

### Efficiency:
- Connection pooling reduces latency by 80-90%
- Batch operations are 10-100x faster
- Only fetching needed fields reduces network transfer

### Reliability:
- Retry logic handles transient failures
- Proper resource cleanup prevents leaks
- Validation prevents bad data

### Maintainability:
- Logging makes debugging easier
- Tests prevent regressions
- Documentation helps onboarding
- Separation of concerns makes changes safer

### Security:
- No hardcoded credentials
- Input validation prevents injection
- Sanitization prevents XSS/injection attacks
- Configuration validation prevents misconfiguration

### Scalability:
- Connection pooling supports many concurrent users
- Batch operations handle high throughput
- No global state enables multi-threading
- Proper resource management prevents exhaustion

---

## Balance Maintained

These prompts:
- ✅ Use simple, clear language
- ✅ Explain the "why" behind each improvement
- ✅ Don't prescribe exact syntax
- ✅ Focus on patterns and principles
- ✅ Build incrementally (core → enhancements)
- ✅ Provide practical benefits for each change

Developers understand WHAT quality means and WHY it matters, without being told HOW to write every line.
