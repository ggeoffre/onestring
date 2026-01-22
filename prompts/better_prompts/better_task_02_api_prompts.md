# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 2: API - Better Quality Prompts
## Focus: Secure, Performant, Production-Ready Web API

---

## Quality-Focused Core Prompts

### Prompt 1: Secure Flask Application Setup
```
Create a Flask application with security best practices:
- Generate a real SECRET_KEY from environment variable (use os.urandom() as fallback)
- Set DEBUG based on environment variable (default False for production)
- Configure ALLOWED_HOSTS from environment variable (don't allow all hosts)
- Enable CORS only for specific origins (not *)
- Set secure headers (X-Content-Type-Options, X-Frame-Options)
- Disable Flask's development server in production (use proper WSGI)
- Bind to 0.0.0.0:8080
- Log all configuration at startup

Why: Security must be built in from the start, not added later.
```

### Prompt 2: Request Validation Middleware
```
Create a decorator called @validate_json that:
- Checks if request has JSON content type
- Validates JSON is well-formed
- Returns 415 (Unsupported Media Type) if not JSON
- Returns 400 (Bad Request) if JSON is malformed
- Logs validation failures with request details
- Can be applied to any route that expects JSON

Example: @validate_json decorator can wrap route functions.

Why: Validates all inputs at the entry point, preventing bad data from reaching business logic.
```

### Prompt 3: Error Handling with Consistent Responses
```
Create an error handler system that:
- Catches all exceptions globally
- Returns consistent JSON error responses with this structure:
  {"error": "error_type", "message": "human readable", "details": {}}
- Uses proper HTTP status codes (400, 404, 500, etc.)
- Logs errors with full stack trace and request context
- Never exposes internal details in production (no stack traces to client)
- Has specific handlers for: ValidationError, DatabaseError, NotFoundError

Why: Consistent error handling makes APIs easier to use and debug.
```

### Prompt 4: Routes with Proper HTTP Semantics
```
Create Flask routes following REST best practices:

GET / (health check):
- Returns {"status": "healthy", "timestamp": current_time}
- Status 200 if database is reachable, 503 if not
- Includes app version in response

POST /echo:
- Validates JSON input (use @validate_json)
- Returns the exact JSON received
- Status 200 on success
- Never modifies the input

POST /data (instead of /log):
- Validates JSON has required sensor fields
- Stores in database using repository pattern
- Returns {"id": generated_id, "message": "stored"} with status 201 (Created)
- Returns 400 if validation fails
- Returns 503 if database unavailable

GET /data (instead of /report):
- Supports query parameters: ?limit=100&offset=0&location=den
- Returns paginated results
- Sets proper Cache-Control headers (cache for 60 seconds)
- Returns CSV if Accept header is text/csv, JSON otherwise
- Returns 200 with empty list if no data
- Includes pagination metadata in response

DELETE /data (instead of /purge):
- Requires confirmation query parameter: ?confirm=true
- Returns count of deleted records
- Returns 403 if confirm not provided (safety feature)
- Status 200 on success

Why: Proper HTTP methods and status codes make the API intuitive and standard-compliant.
```

### Prompt 5: Response Formatting with Consistency
```
Create helper functions for consistent responses:

1. success_response(data, status=200) -> Flask Response
   - Wraps data in consistent structure
   - Sets proper content-type
   - Adds timing header (how long request took)

2. error_response(error, status=400) -> Flask Response
   - Creates consistent error structure
   - Logs error details
   - Returns safe message to client

3. csv_response(data, filename) -> Flask Response
   - Converts data to CSV
   - Sets proper headers for download
   - Sets cache headers
   - Handles empty data gracefully

Use these helpers in all routes for consistency.

Why: Consistent responses make the API easier to consume and debug.
```

### Prompt 6: Structured Logging
```
Set up Flask logging that:
- Logs all requests with: method, path, IP, user agent, response time
- Logs all responses with: status code, response size
- Uses JSON format for logs (easier to parse)
- Logs to both file and console
- Rotates log files daily (keep 7 days)
- Different log levels for different environments (DEBUG in dev, INFO in prod)
- Includes request ID in all logs (for tracing)

Why: Structured logging makes debugging and monitoring much easier.
```

### Prompt 7: Rate Limiting for Protection
```
Add rate limiting using flask-limiter:
- Limit POST /data to 100 requests per minute per IP
- Limit DELETE /data to 10 requests per hour per IP
- Limit GET /data to 1000 requests per minute per IP
- Return 429 (Too Many Requests) when limit exceeded
- Include Retry-After header in response
- Log rate limit violations

Why: Protects the API from abuse and ensures fair usage.
```

---

## Quality Enhancement Prompts

### Prompt 8: Input Sanitization
```
Add input sanitization to all routes:
- Strip leading/trailing whitespace from all strings
- Validate string lengths (location max 100 chars, sensor max 50 chars)
- Ensure numeric values are within reasonable ranges
- Reject payloads larger than 10KB
- Remove any HTML/script tags from strings (prevent XSS)
- Validate timestamps are reasonable (not in distant future/past)

Why: Sanitization prevents injection attacks and ensures data quality.
```

### Prompt 9: Request ID Tracing
```
Add request ID tracking:
- Generate unique ID for each request (UUID)
- Add request ID to all log entries
- Return request ID in response headers (X-Request-ID)
- Include request ID in error responses
- Use request ID to trace request through entire system

Why: Makes debugging much easier by tracking requests end-to-end.
```

### Prompt 10: Health Check Endpoint
```
Enhance GET / to be a proper health check:
- Check database connectivity (with timeout)
- Check available memory
- Check disk space
- Return detailed status in development, simple status in production
- Return 200 if healthy, 503 if any check fails
- Include response time for each check
- Cache health check results for 10 seconds (don't overload database)

Why: Load balancers and monitoring tools need health checks.
```

### Prompt 11: API Versioning
```
Add API versioning support:
- Version routes under /v1/ prefix (e.g., /v1/data)
- Include version in response headers (X-API-Version: 1.0)
- Support version in Accept header (Accept: application/vnd.api.v1+json)
- Document which version each route belongs to
- Keep health check unversioned (always at /)

Why: Allows API evolution without breaking existing clients.
```

### Prompt 12: Response Compression
```
Enable gzip compression for responses:
- Compress responses larger than 1KB
- Support gzip and deflate encoding
- Only compress text responses (JSON, CSV)
- Add Content-Encoding header
- Improves performance by reducing response size by 70-90%

Why: Reduces bandwidth usage and improves response times.
```

### Prompt 13: Request Timeout Protection
```
Add timeout protection:
- Set maximum request processing time to 30 seconds
- Return 504 (Gateway Timeout) if exceeded
- Cancel database operations if timeout reached
- Log slow requests (>1 second) with WARNING
- Use a decorator @timeout(seconds=30) for routes

Why: Prevents slow requests from tying up resources.
```

### Prompt 14: CORS Configuration
```
Configure CORS properly:
- Read allowed origins from environment variable
- Default to empty list (no CORS) if not configured
- Support preflight OPTIONS requests
- Set proper Access-Control headers:
  * Allow-Origin (specific origins, not *)
  * Allow-Methods (only methods you support)
  * Allow-Headers (only headers you need)
  * Max-Age (cache preflight for 1 hour)
- Log CORS rejections

Why: Enables secure cross-origin requests for web clients.
```

---

## Testing & Documentation Prompts

### Prompt 15: API Integration Tests
```
Create test_api.py using pytest and Flask test client:
- Test all routes with valid inputs (happy path)
- Test all routes with invalid inputs (error cases)
- Test rate limiting (make too many requests)
- Test authentication if added
- Test response formats (JSON, CSV)
- Test error handling (database down, invalid JSON)
- Use fixtures for common test data
- Mock database calls (test API logic, not database)

Why: Tests ensure the API works correctly and catches regressions.
```

### Prompt 16: OpenAPI/Swagger Documentation
```
Add OpenAPI documentation using flask-swagger-ui:
- Document all routes with parameters and responses
- Include example requests and responses
- Document all error codes
- Show authentication requirements
- Make documentation available at /api/docs
- Generate documentation from code annotations (not manual docs)
- Keep documentation up to date automatically

Why: Good documentation makes the API easy to use and reduces support burden.
```

### Prompt 17: Performance Testing
```
Create a performance test script (locust or similar):
- Simulates 100 concurrent users
- Tests POST /data with random data
- Tests GET /data with various parameters
- Measures response times (p50, p95, p99)
- Identifies bottlenecks
- Fails if response time > 200ms for 95% of requests

Why: Ensures the API can handle production load.
```

---

## Expected Quality Improvements

### API Quality Metrics:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Security** | Poor (DEBUG=True, weak key) | Excellent (hardened) | ⬆️ 95% |
| **Error Handling** | Poor (inconsistent) | Excellent (standardized) | ⬆️ 90% |
| **Performance** | Fair (no caching) | Good (cached, compressed) | ⬆️ 70% |
| **Reliability** | Poor (no rate limits) | Good (protected) | ⬆️ 85% |
| **Observability** | Poor (print only) | Excellent (structured logs) | ⬆️ 95% |
| **Usability** | Fair (works) | Excellent (documented, consistent) | ⬆️ 80% |
| **Testability** | Poor (hard to test) | Excellent (easily testable) | ⬆️ 90% |

---

## Production Readiness Checklist

After these prompts, the API will have:
- ✅ Proper security (secrets, headers, validation)
- ✅ Error handling (consistent, logged)
- ✅ Monitoring (health checks, structured logs)
- ✅ Performance (caching, compression, rate limits)
- ✅ Documentation (OpenAPI/Swagger)
- ✅ Testing (unit, integration, performance)
- ✅ Standards compliance (REST, HTTP semantics)
- ✅ Scalability considerations (stateless, no resource leaks)

---

## Key API Design Principles Applied

### 1. **REST Principles**
- Resources, not actions (GET /data not GET /report)
- Proper HTTP methods (GET, POST, DELETE)
- Correct status codes (200, 201, 400, 404, 500)
- Stateless (no server-side sessions)

### 2. **Security First**
- Validate all inputs
- Sanitize all outputs
- Rate limiting
- Secure defaults
- No information leakage

### 3. **Fail Gracefully**
- Never expose stack traces
- Return helpful error messages
- Log errors for debugging
- Always return JSON (even for errors)

### 4. **Performance**
- Cache when possible
- Compress responses
- Paginate large results
- Efficient queries

### 5. **Observability**
- Structured logging
- Request tracing
- Performance metrics
- Health checks

---

## Why This Produces Better APIs

### Security:
- Environment-based secrets prevent credential leaks
- Input validation prevents injection attacks
- Rate limiting prevents abuse
- Secure headers prevent common vulnerabilities

### Reliability:
- Timeouts prevent hanging requests
- Error handling prevents crashes
- Health checks enable monitoring
- Rate limits ensure fair usage

### Performance:
- Compression reduces bandwidth by 70-90%
- Caching reduces database load
- Pagination prevents memory issues
- Connection pooling improves throughput

### Maintainability:
- Structured logging makes debugging easier
- Tests catch bugs early
- Documentation reduces onboarding time
- Consistent patterns make changes safer

### Usability:
- Consistent responses make API predictable
- Good error messages help developers debug
- Documentation shows examples
- Proper HTTP semantics feel natural

---

## Balance Maintained

These prompts:
- ✅ Use simple language (not pseudo-code)
- ✅ Explain WHY each improvement matters
- ✅ Build incrementally (basics → enhancements)
- ✅ Focus on patterns, not syntax
- ✅ Provide practical examples
- ✅ Address real production concerns

Developers understand what production-ready means without being given line-by-line code.
