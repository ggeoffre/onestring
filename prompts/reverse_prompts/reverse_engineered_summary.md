# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

# Reverse-Engineered Prompts: Summary
## Six Projects from Actual Working Code

---

## Overview

These prompts are **reverse-engineered from the actual generated code** in all_python_files.txt. Each set of prompts will reproduce the working code with high sufficiency (75-85%) using simple, clear language.

---

## The Six Projects

### **PROJECT 1: python/data**
**Focus:** Standalone database CRUD operations

**Files:** 8 files (main.py + 5 database modules + helper + requirements)

**What it does:**
- Raw CRUD operations for 5 databases (MongoDB, Cassandra, MySQL, PostgreSQL, Redis)
- Each database in its own module with complete lifecycle (connect, setup, store, retrieve, delete, close)
- Helper module with data generation and conversion utilities
- Main script calls all databases in sequence

**Key Pattern:** Self-contained database operations with consistent structure

---

### **PROJECT 2: python/api/flask-app**
**Focus:** Minimal Flask API

**Files:** 2 files (app.py + requirements)

**What it does:**
- Simple Flask API with 5 routes (/, /echo, /log, /report, /purge)
- JSON-to-CSV conversion without external library
- Static sample data constant
- Basic error handling

**Key Pattern:** Standalone Flask app with manual CSV generation

---

### **PROJECT 3: python/api/django-app**
**Focus:** Minimal Django API

**Files:** 2 files (app.py + requirements)

**What it does:**
- Single-file Django app (settings.configure in same file)
- Same 5 routes as Flask version
- JSON-to-CSV conversion
- Request method validation

**Key Pattern:** Django without separate settings.py, mirrors Flask functionality

---

### **PROJECT 4: python/access**
**Focus:** Abstract base class with stub implementations

**Files:** 7 files (protocol + 5 implementations + main)

**What it does:**
- Abstract base class defines interface
- 5 stub implementations (just print, no real database)
- Factory pattern for selecting implementation
- Environment-based configuration
- Demonstrates polymorphism

**Key Pattern:** Interface-based design with factory pattern

---

### **PROJECT 5: python/legos/flask-app**
**Focus:** Flask + Real database implementations

**Files:** 9 files (app + protocol + 5 real databases + helper + requirements)

**What it does:**
- Combines Projects 1, 2, and 4
- Flask API routes that actually call real databases
- Factory pattern selects database implementation
- CSV export from database data
- Full schema creation and CRUD

**Key Pattern:** Integration of API layer + data access layer + utilities

---

### **PROJECT 6: python/legos/django-app**
**Focus:** Django + Real database implementations

**Files:** 9 files (app + protocol + 5 real databases + helper + requirements)

**What it does:**
- Same as Project 5 but with Django instead of Flask
- All database code is IDENTICAL to Project 5
- Only app.py differs (Django vs Flask specifics)
- Demonstrates framework independence of data layer

**Key Pattern:** Same integration as Project 5, different web framework

---

## Progression & Learning Path

### **Learning Sequence:**

1. **PROJECT 1 (data):** Learn database operations
   - Connect, CRUD, close
   - Environment variables
   - Try-except-finally pattern
   - Helper functions

2. **PROJECT 2 (api/flask):** Learn Flask basics
   - Routes and methods
   - JSON handling
   - Response types
   - Error codes

3. **PROJECT 3 (api/django):** Learn Django basics
   - Single-file app
   - Django views
   - Request/response handling
   - URL patterns

4. **PROJECT 4 (access):** Learn abstractions
   - Abstract base classes
   - Interface pattern
   - Factory pattern
   - Polymorphism

5. **PROJECT 5 (legos/flask):** Integrate everything with Flask
   - Combine data + API + abstractions
   - Real working system
   - Environment-based database selection

6. **PROJECT 6 (legos/django):** Show portability
   - Same data layer, different API layer
   - Framework independence
   - Reusability principle

---

## Sufficiency Ratings

| Project | Files | Prompts Needed | Estimated Sufficiency | Interactions Expected |
|---------|-------|----------------|----------------------|---------------------|
| 1. data | 8 | 8 | 75-80% | 2-3 |
| 2. api/flask | 2 | 2 | 80-85% | 1-2 |
| 3. api/django | 2 | 2 | 80-85% | 1-2 |
| 4. access | 7 | 7 | 85-90% | 1-2 |
| 5. legos/flask | 9 | 9 | 75-80% | 2-4 |
| 6. legos/django | 9 | 4* | 75-80% | 2-4 |

*Project 6 only needs 4 new prompts because database implementations are reused from Project 5.

---

## Why These Prompts Work Better

### **1. Based on Actual Code**
- Reverse-engineered from real, working code
- Not theoretical or aspirational
- Prompts describe what actually exists

### **2. Right Level of Detail**
- Not too vague (missing key details)
- Not too detailed (pseudo-code)
- Sweet spot for LLM generation

### **3. Consistent Patterns**
- Same structure across similar files
- Predictable naming conventions
- Familiar patterns repeated

### **4. Incremental Complexity**
- Starts simple (Project 1: one function per file)
- Builds up (Project 4: abstractions)
- Integrates (Project 5/6: full system)

### **5. Clear Specifications**
- Exact variable names where they matter
- Specific values (ports, defaults)
- Concrete error messages
- Explicit return types

---

## Common Patterns Across All Projects

### **Configuration Pattern:**
```
- Read from environment variable
- Provide sensible default
- Use .lower() for hostnames
- Store in module-level constants
```

### **Error Handling Pattern:**
```
- Try-except blocks around operations
- Print errors (not raise in these examples)
- Return None or empty list on error
- Finally blocks for cleanup
```

### **Database Operation Pattern:**
```
1. Connect (with configuration)
2. Setup (create schema if needed)
3. Store (insert data)
4. Retrieve (select data)
5. Delete (truncate/delete)
6. Close (cleanup)
```

### **API Route Pattern:**
```
- GET / for health check
- POST /echo for testing
- POST /log for storing
- GET /report for retrieving
- GET/POST /purge for deleting
```

---

## Key Differences Between Projects

### **Data Layer:**
| Aspect | Project 1 | Projects 5/6 |
|--------|-----------|-------------|
| Structure | Functions | Classes |
| Interface | None | Abstract base class |
| Selection | N/A | Factory pattern |
| Integration | Standalone | Called by API |

### **API Layer:**
| Aspect | Projects 2/3 | Projects 5/6 |
|--------|-------------|-------------|
| Data | Static constant | From database |
| Storage | None | Real database |
| CSV | From constant | From live data |

### **Architecture:**
| Aspect | Projects 1-4 | Projects 5-6 |
|--------|-------------|-------------|
| Separation | Single concern | Multiple layers |
| Integration | None | Full stack |
| Complexity | Low | Medium |

---

## Using These Prompts

### **For Each Project:**

1. **Read the prompts in order**
   - They build on each other within a project
   - Later prompts reference earlier ones

2. **Generate files one at a time**
   - Create each file individually
   - Test as you go
   - Don't try to generate everything at once

3. **Use the exact structure shown**
   - File names matter
   - Import paths matter
   - Variable names in some cases matter

4. **Test incrementally**
   - Test each module before integrating
   - Verify database connections work
   - Check API routes individually

5. **Refer to the "Key Patterns" sections**
   - These explain WHY the code is structured this way
   - Understanding patterns helps with troubleshooting

---

## Expected Generation Quality

With these prompts, you should get:

### **First Generation (No Follow-up):**
- ✅ Correct file structure
- ✅ Correct imports
- ✅ Correct function signatures
- ✅ Correct basic logic
- ⚠️ Possibly minor syntax issues
- ⚠️ Possibly import order differences

### **After 1-2 Clarifications:**
- ✅ All syntax correct
- ✅ All imports correct
- ✅ Code runs successfully
- ✅ Matches original behavior

### **What Might Need Tweaking:**
- Comment styles (not critical)
- Print message wording (not critical)
- Variable naming in local scopes (not critical)
- Error message phrasing (not critical)

---

## Validation Checklist

For each project, verify:

**Structure:**
- [ ] All files created
- [ ] File names match exactly
- [ ] Directory structure correct

**Imports:**
- [ ] No missing imports
- [ ] No circular imports
- [ ] All modules importable

**Functionality:**
- [ ] Database connections work
- [ ] CRUD operations succeed
- [ ] API routes respond correctly
- [ ] CSV export works
- [ ] Error handling graceful

**Configuration:**
- [ ] Environment variables read correctly
- [ ] Defaults work
- [ ] Factory pattern selects correctly (Projects 5/6)

---

## Troubleshooting Guide

### **Import Errors:**
- Check file names match imports exactly
- Ensure all files in same directory
- Verify Python path includes directory

### **Connection Errors:**
- Set DATA_HOSTNAME environment variable
- Check database is running
- Verify ports are correct
- Test connection manually

### **Type Errors:**
- Check function signatures match protocol
- Verify return types
- Ensure None vs [] consistency

### **Missing Attributes:**
- Check class names match exactly
- Verify abstract methods implemented
- Ensure all required methods present

---

## Success Metrics

You'll know the prompts worked if:

1. **All files generate on first try** (structure correct)
2. **Code runs with 0-2 fixes** (high sufficiency)
3. **Database operations work** (logic correct)
4. **API responds correctly** (integration works)
5. **CSV export produces valid CSV** (utilities work)
6. **Can switch databases via env var** (factory works)

---

## Conclusion

These 6 prompt sets represent:
- **Real code that actually runs**
- **Simple, clear language** (not pseudo-code)
- **High sufficiency** (75-90%)
- **Minimal follow-up** (1-4 interactions)
- **Incremental complexity** (learn as you go)
- **Complete working systems** (not just fragments)

They successfully **close the gap** identified in the original analysis by providing prompts that:
1. Start from actual working code (not aspirational)
2. Use appropriate level of detail (not too vague or too detailed)
3. Build incrementally (not overwhelming)
4. Demonstrate real patterns (not academic examples)

Most importantly, they maintain the **pedagogical value** of the original assignments:
- Students still discover through experimentation
- Prompts guide without prescribing
- Iteration is still required (but minimal)
- Learning happens through doing
