# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# TASK 4: LEGOS - Better Quality Prompts
## Focus: Production-Ready Integration with Best Practices

---

## Quality-Focused Integration Prompts

### Prompt 1: Application Architecture
```
Design the application with clear separation of layers:

1. **API Layer** (app.py):
   - Handles HTTP requests/responses
   - Validates and sanitizes inputs
   - Returns consistent response formats
   - No business logic here

2. **Service Layer** (sensor_service.py):
   - Contains business logic
   - Orchestrates data access operations
   - Handles validation and transformation
   - Independent of HTTP framework

3. **Data Access Layer** (database implementations):
   - Only database operations
   - No business logic
   - Uses repository pattern

4. **Configuration Layer** (config.py):
   - Centralized configuration
   - Environment variable management
   - Validation of settings

Each layer only depends on the layer below it. No circular dependencies.

Why: Clean architecture makes the system maintainable, testable, and scalable.
```

### Prompt 2: Application Configuration
```
Create app_config.py that manages all configuration:
- Load settings from environment variables
- Support .env file for development
- Validate all required settings at startup
- Fail fast if critical settings missing
- Include settings for:
  * Flask (SECRET_KEY, DEBUG, ALLOWED_HOSTS)
  * Database (type, connection params, pool size)
  * Logging (level, format, handlers)
  * Security (CORS origins, rate limits)
  * Performance (cache TTL, timeouts)
- Use dataclasses or pydantic for type safety
- Provide sensible defaults for development
- No defaults for secrets (must be explicit)

Why: Centralized configuration is easier to manage and validate.
```

### Prompt 3: Application Factory Pattern
```
Create the Flask app using factory pattern:
- Define create_app(config_name='default') function
- Load appropriate configuration
- Initialize extensions (logging, CORS, rate limiter)
- Register blueprints/routes
- Register error handlers
- Initialize database connection pool
- Register cleanup handlers
- Return configured app

Support different configurations:
- 'development': DEBUG=True, verbose logging
- 'testing': In-memory database, no rate limits
- 'production': DEBUG=False, strict settings

Why: Factory pattern makes testing easier and supports multiple environments.
```

### Prompt 4: Service Layer for Business Logic
```
Create sensor_service.py with SensorService class:
- Takes data access repository in __init__ (dependency injection)
- Provides high-level operations:
  * store_reading(data: dict) -> Result[str, Error]
  * get_readings(filters: dict, page: int, size: int) -> Result[Page, Error]
  * export_readings(format: str) -> Result[bytes, Error]
  * clear_readings(confirm: bool) -> Result[int, Error]
- Handles validation before calling data access
- Handles business rules (e.g., deduplicate readings)
- Converts between domain models and DTOs
- Returns Result type (success or error, never exceptions)
- Logs all operations with context

Why: Service layer keeps business logic out of routes and data access.
```

### Prompt 5: Request/Response DTOs
```
Create models.py with data transfer objects:
- Define Pydantic models (or dataclasses) for:
  * SensorReadingRequest (input validation)
  * SensorReadingResponse (output format)
  * PagedResponse (pagination wrapper)
  * ErrorResponse (error format)
  * HealthCheckResponse (health check format)
- Include field validation rules
- Automatic JSON serialization/deserialization
- Type hints for all fields
- Example values in docstrings

Use these DTOs in routes to validate inputs and format outputs.

Why: DTOs ensure consistent data shapes and enable automatic validation.
```

### Prompt 6: API Routes with Clean Code
```
Create routes using Flask blueprints:

api_v1_bp = Blueprint('api_v1', __name__, url_prefix='/v1')

Each route should:
- Use a single service instance (injected)
- Validate input with DTOs
- Call service layer (no direct database access)
- Handle Result type from service
- Return consistent response format
- Include proper HTTP status codes
- Log request/response
- Be < 20 lines of code

Routes:
- GET /v1/health: Health check
- POST /v1/readings: Store reading
- GET /v1/readings: List readings (paginated)
- GET /v1/readings/{id}: Get specific reading
- DELETE /v1/readings: Clear all (with confirmation)
- GET /v1/export: Export as CSV/JSON

Why: Clean, focused routes are easier to understand and maintain.
```

### Prompt 7: Global Error Handler
```
Create error_handlers.py with handlers for:
- ValidationError (400): Input validation failed
- NotFoundError (404): Resource not found
- UnauthorizedError (401): Authentication failed
- ForbiddenError (403): Not allowed to access
- DatabaseError (503): Database unavailable
- TimeoutError (504): Operation timed out
- Generic Exception (500): Unexpected error

All handlers should:
- Return consistent JSON error format
- Log error with full context and stack trace
- Include request ID in response
- Never expose sensitive details in production
- Set appropriate HTTP status code
- Include error code for client handling

Why: Consistent error handling improves API usability and debugging.
```

### Prompt 8: Middleware for Cross-Cutting Concerns
```
Add middleware for:
1. Request ID: Generate and track unique ID per request
2. Request Logging: Log all requests with timing
3. Response Logging: Log all responses with size
4. Security Headers: Add security headers to all responses
5. CORS: Handle cross-origin requests
6. Rate Limiting: Enforce rate limits per route
7. Authentication: Verify API keys (if needed)
8. Compression: Gzip compress large responses

Each middleware should:
- Be composable and reusable
- Have single responsibility
- Be testable independently
- Log its actions
- Handle errors gracefully

Why: Middleware keeps cross-cutting concerns separate from business logic.
```

---

## Quality Enhancement Prompts

### Prompt 9: Graceful Shutdown
```
Implement graceful shutdown:
- Register signal handlers (SIGTERM, SIGINT)
- Finish processing current requests
- Close database connections
- Flush logs
- Exit with appropriate code
- Don't accept new requests during shutdown
- Timeout after 30 seconds (force exit)

Why: Prevents data loss and ensures clean resource cleanup.
```

### Prompt 10: Health Check Endpoints
```
Create comprehensive health checks:
- GET /health: Simple up/down check
- GET /health/live: Liveness probe (is app running?)
- GET /health/ready: Readiness probe (can app serve traffic?)
- GET /health/detail: Detailed component health

Check:
- Database connectivity (with timeout)
- Disk space availability
- Memory usage
- Active connections
- Recent error rate

Return 200 if healthy, 503 if not ready.

Why: Enables proper orchestration with Kubernetes/load balancers.
```

### Prompt 11: Structured Logging
```
Configure logging with:
- JSON format for structured logs
- Fields: timestamp, level, message, context, request_id
- Log to stdout in production (for container logging)
- Rotate file logs daily in development
- Different levels per module (DEBUG for development, INFO for production)
- Correlation IDs to trace requests across services
- Performance metrics (request duration, database time)

Create logging utilities:
- log_request(request): Log incoming request
- log_response(response): Log outgoing response
- log_error(error, context): Log error with full context
- log_metric(name, value): Log performance metric

Why: Structured logs are much easier to search and analyze.
```

### Prompt 12: Metrics and Monitoring
```
Add Prometheus-style metrics:
- Request count by route and status code
- Request duration histogram
- Database operation count and duration
- Error count by type
- Active connections
- Cache hit rate

Expose metrics at /metrics endpoint.

Track custom business metrics:
- Readings stored per minute
- Most common sensor types
- Average reading value by location

Why: Metrics enable proactive monitoring and alerting.
```

### Prompt 13: Circuit Breaker Pattern
```
Implement circuit breaker for database operations:
- Track failure rate for database operations
- Open circuit after 5 consecutive failures
- Return cached data or error when circuit open
- Half-open after 30 seconds (try one request)
- Close circuit if request succeeds
- Log circuit state changes

Use a library like pybreaker or implement simple version.

Why: Prevents cascading failures when database is down.
```

### Prompt 14: Caching Strategy
```
Implement multi-level caching:
1. In-memory cache (for hot data):
   - Cache readings for last 5 minutes
   - LRU eviction policy
   - TTL: 60 seconds

2. Redis cache (for warm data):
   - Cache aggregations and exports
   - TTL: 300 seconds
   - Invalidate on writes

3. HTTP caching:
   - Cache-Control headers
   - ETag support
   - Conditional requests

Cache key format: f"{operation}:{params}:{version}"

Why: Caching dramatically reduces database load and improves response times.
```

### Prompt 15: Background Jobs
```
Add background job processing:
- Export large datasets asynchronously
- Aggregate data in background
- Clean up old data periodically
- Use Celery or APScheduler for job queue

Jobs:
- export_to_csv(job_id, filters): Export data to file
- aggregate_daily_stats(): Calculate daily statistics
- cleanup_old_data(days): Delete old readings

Track job status:
- POST /jobs: Submit job
- GET /jobs/{id}: Check status
- GET /jobs/{id}/result: Get result when done

Why: Background jobs keep API responsive and handle long-running operations.
```

---

## Testing & Quality Assurance

### Prompt 16: Comprehensive Test Suite
```
Create test suite with:
1. Unit tests (test_service.py):
   - Test service layer logic
   - Mock data access layer
   - Test validation rules
   - Test error handling

2. Integration tests (test_integration.py):
   - Test API → Service → Database flow
   - Use test database
   - Test with real HTTP requests
   - Test error scenarios

3. API tests (test_api.py):
   - Test all routes
   - Test authentication
   - Test rate limiting
   - Test error responses

4. Load tests (test_load.py):
   - Simulate 100 concurrent users
   - Test for 5 minutes
   - Ensure p95 latency < 200ms
   - Ensure error rate < 0.1%

Use pytest, fixtures, and mocks effectively.
Achieve >85% code coverage.

Why: Comprehensive testing catches bugs before production.
```

### Prompt 17: Performance Benchmarks
```
Create performance benchmarks:
- Benchmark each route under load
- Measure database operation latency
- Measure cache hit rates
- Measure memory usage
- Compare with/without caching
- Compare different database implementations

Set performance targets:
- P95 latency < 200ms
- Throughput > 1000 req/sec
- Memory < 512MB
- CPU < 50% average

Fail CI/CD if targets not met.

Why: Performance benchmarks prevent regressions.
```

### Prompt 18: Security Testing
```
Add security tests:
- SQL injection attempts
- XSS attempts in inputs
- CSRF protection
- Rate limit bypass attempts
- Authentication bypass attempts
- Unauthorized access attempts
- Input fuzzing

Use tools like OWASP ZAP or safety.

Document security assumptions and threat model.

Why: Security testing finds vulnerabilities before attackers do.
```

---

## Deployment & Operations

### Prompt 19: Docker Configuration
```
Create production-ready Dockerfile:
- Use specific Python version (not :latest)
- Multi-stage build (build stage + runtime stage)
- Run as non-root user
- Include only necessary files
- Set appropriate environment variables
- Health check instruction
- Expose port 8080
- Use gunicorn (not Flask dev server)

Create docker-compose.yml for local development:
- App container
- Database containers (MongoDB, PostgreSQL, Redis)
- Shared networks
- Volume mounts for data persistence

Why: Containerization ensures consistency across environments.
```

### Prompt 20: Production Deployment Guide
```
Create deployment documentation:
- Environment variables reference
- Database setup instructions
- Migration procedure
- Scaling recommendations
- Monitoring setup
- Alerting rules
- Backup strategy
- Disaster recovery plan

Include runbook for common issues:
- High error rate
- High latency
- Database connection issues
- Memory leaks

Why: Good documentation enables reliable operations.
```

---

## Expected Quality Improvements

### System Quality Metrics:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Architecture** | Monolithic | Layered | ⬆️ 85% |
| **Performance** | Baseline | Optimized + Cached | ⬆️ 400% |
| **Reliability** | Poor (crashes) | Excellent (resilient) | ⬆️ 95% |
| **Security** | Poor (vulnerabilities) | Good (hardened) | ⬆️ 90% |
| **Observability** | Poor (blind) | Excellent (full visibility) | ⬆️ 95% |
| **Testability** | Poor (hard to test) | Excellent (DI + mocks) | ⬆️ 90% |
| **Maintainability** | Fair (tightly coupled) | Good (loosely coupled) | ⬆️ 85% |
| **Scalability** | Poor (resource leaks) | Good (efficient) | ⬆️ 500% |

---

## Production Readiness Checklist

After these prompts, the system will have:
- ✅ Clean layered architecture
- ✅ Dependency injection
- ✅ Comprehensive error handling
- ✅ Structured logging and metrics
- ✅ Security hardening
- ✅ Performance optimization
- ✅ Circuit breakers and resilience
- ✅ Background job processing
- ✅ Health checks and monitoring
- ✅ Comprehensive test coverage
- ✅ Documentation and runbooks
- ✅ Container deployment

---

## Key Integration Principles Applied

### 1. **Separation of Concerns**
- API layer handles HTTP
- Service layer handles business logic
- Data layer handles persistence
- Each layer has single responsibility

### 2. **Dependency Injection**
- Services receive dependencies
- Easy to mock for testing
- Configuration-driven
- No hard-coded dependencies

### 3. **Fail Fast and Safe**
- Validate configuration at startup
- Validate inputs at entry points
- Circuit breakers prevent cascading failures
- Graceful degradation

### 4. **Observable and Debuggable**
- Structured logging
- Request tracing
- Performance metrics
- Health checks

### 5. **Secure by Default**
- No hardcoded secrets
- Input validation and sanitization
- Rate limiting
- Secure headers
- Least privilege

---

## Why This Produces Production-Ready Systems

### Reliability:
- Circuit breakers handle failures
- Graceful shutdown prevents data loss
- Health checks enable monitoring
- Retry logic handles transient issues

### Performance:
- Multi-level caching
- Connection pooling
- Background jobs for heavy work
- Efficient database queries

### Security:
- Input validation at boundaries
- No SQL injection vulnerabilities
- Secrets from environment
- Security headers
- Rate limiting

### Maintainability:
- Clean architecture is easier to change
- Comprehensive tests catch regressions
- Logging makes debugging easier
- Documentation helps new team members

### Scalability:
- Stateless design enables horizontal scaling
- Connection pooling handles high concurrency
- Caching reduces database load
- Background jobs handle spikes

---

## Balance Maintained

These prompts:
- ✅ Use simple, clear language
- ✅ Explain the WHY behind architecture choices
- ✅ Focus on patterns and principles, not code
- ✅ Build incrementally from basics to advanced
- ✅ Address real production concerns
- ✅ Provide practical examples

Developers understand what production-ready means without getting a computer science lecture.

---

## Success Metrics

A system built with these prompts should achieve:
- **Uptime**: 99.9% (< 9 hours downtime/year)
- **Performance**: P95 latency < 200ms
- **Scalability**: Handle 1000+ req/sec
- **Security**: Pass security audit
- **Maintainability**: New features in days, not weeks
- **Reliability**: Graceful degradation when things fail
- **Observability**: Debug issues in minutes, not hours

This is what "better code" means in production.
