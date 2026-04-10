# Midterm Project - QA Implementation and Empirical Analysis

**Course:** Advanced Quality Assurance
**Project:** E-Commerce Platform
**Author:** QA Team
**Date:** 2026-04-10
**Repository:** https://github.com/Choppa2077/aqaBackend

---

## 1. System Description

### Architecture

The system is a monolithic REST API backend written in Go (version 1.24), paired with a React/TypeScript frontend. The backend follows a layered architecture:

- **Delivery layer** (`internal/delivery/rest/v1/`) - HTTP handlers using Gin framework, DTO validation
- **Service layer** (`internal/service/`) - business logic, the primary target of unit tests
- **Repository layer** (`internal/repository/`) - database access interfaces
- **Domain layer** (`internal/domain/`) - shared domain models and error types

### Technologies

- **Backend:** Go 1.24, Gin HTTP framework, JWT authentication (golang-jwt/jwt), bcrypt password hashing, MongoDB (via driver), PostgreSQL migration in progress
- **Frontend:** React 19, TypeScript, Vite, Axios
- **Testing:** Go `testing` package, `testify/mock` and `testify/assert`, Playwright (Chromium) for E2E
- **CI/CD:** GitHub Actions
- **Documentation:** Swagger (swaggo/gin-swagger)

### Key Functionalities

1. **Authentication** - user registration, login, JWT token generation and refresh
2. **Product Catalog** - CRUD for products and categories, filtering, search, statistics
3. **User Interactions** - likes, views, purchase history
4. **Purchase Flow** - stock validation, atomic purchase recording
5. **User Profile** - profile management, password change, account deletion
6. **Recommendations** - collaborative filtering based on user interaction history

---

## 2. Methodology

### Risk-Based Testing Approach

Testing priority was determined by a risk matrix from Assignment 1. The three highest-risk modules were:

1. **Auth** (risk score 9/10) - security-critical, all user access depends on JWT validity
2. **Purchase Flow** (risk score 9/10) - financial impact, inventory accuracy
3. **Product CRUD** (risk score 8/10) - data integrity, core business function

These modules received the most test coverage in Assignments 1 and 2. In the Midterm, coverage was extended to edge cases, concurrency scenarios, and HTTP handler integration tests.

### Test Design Strategy

Tests are designed following the principle of isolation - each unit test mocks all repository dependencies, testing only service-layer logic. This approach:

- Eliminates database dependency (tests run without MongoDB or PostgreSQL)
- Provides deterministic, repeatable results
- Runs in under 1 second for 46 unit tests

The Midterm extended this with integration tests at the HTTP handler level, using `httptest.NewRecorder()` and mock service implementations to test the full request/response cycle.

### Automation Tools

| Tool | Purpose | Justification |
|------|---------|---------------|
| Go `testing` package | Unit and integration test runner | Native Go tool, no external dependency needed for assertions |
| `testify/mock` | Mock repository generation | Allows interface-based mocking without code generation |
| `testify/assert` | Assertion library | Cleaner assertions than raw `t.Fatal` calls |
| `net/http/httptest` | Integration test HTTP server | Go stdlib, no extra dependency needed |
| Gin (test mode) | Router for integration tests | Same router as production, real routing logic tested |
| Playwright (Chromium) | Frontend E2E tests | Supports modern React apps; headless CI-compatible |

---

## 3. Automation Implementation

### Test Structure

```
internal/
  service/
    authService_test.go       - 15 unit tests (Register, Login, JWT, Refresh, Concurrent)
    productService_test.go    - 19 unit tests (CRUD, Stock, Category, Edge Cases, Concurrent)
    interactionService_test.go - 12 unit tests (Purchase, Like, View, History, Unlike, HasPurchased)
  delivery/rest/v1/
    auth_handler_test.go      - 5 integration tests (HTTP level, mock service)
e2e/
  auth.spec.ts               - 9 Playwright E2E tests
  products.spec.ts           - 3 Playwright E2E tests
  profile.spec.ts            - 9 Playwright E2E tests
```

**Total: 46 unit tests + 5 integration tests + 21 E2E tests = 72 tests**

### CI/CD Pipeline

The GitHub Actions pipeline (`.github/workflows/ci.yml`) runs on every push to main:

1. Checkout code
2. Set up Go 1.24
3. Download dependencies (`go mod download`)
4. Build (`go build ./...`) - QG04
5. Static analysis (`go vet ./...`) - QG05
6. Run unit tests (`go test ./internal/service/... -v -cover`) - QG01
7. Run integration tests (`go test ./internal/delivery/... -v`) - QG02
8. Generate coverage report (`go test ./internal/... -coverprofile=coverage.out`)
9. Check coverage threshold (70%) - QG07
10. Upload test results and coverage artifacts

### Quality Gates Definition

| ID | Gate | Threshold | Enforcement |
|----|------|-----------|-------------|
| QG01 | Unit test pass rate | 100% | Pipeline fails if any unit test fails |
| QG02 | Integration test pass rate | 100% | Pipeline fails if any integration test fails |
| QG03 | Critical test failures | 0 | Any failure blocks merge to main |
| QG04 | Build success | 100% | Compilation failure stops pipeline immediately |
| QG05 | Static analysis errors | 0 | `go vet` errors block pipeline |
| QG06 | Test execution time | Less than 2 minutes | Monitored; currently at 1.2s |
| QG07 | Coverage threshold | 70% overall | Added in Midterm; currently partially met |

---

## 4. Results

### Test Execution Results

| Test Suite | Count | PASS | FAIL | Execution Time |
|------------|-------|------|------|----------------|
| Unit tests (service layer) | 46 | 46 | 0 | 0.62s |
| Integration tests (handlers) | 5 | 5 | 0 | 0.55s |
| E2E tests (Playwright) | 21 | 21 | 0 | ~30-45s |
| **Total** | **72** | **72** | **0** | **~32-47s** |

### Coverage Metrics

| Service | Coverage % | Threshold | Status |
|---------|-----------|-----------|--------|
| authService | 75% avg | 70% | PASS |
| productService (tested methods) | 82% avg | 70% | PASS |
| interactionService (tested methods) | 72% avg | 70% | PASS |
| delivery/rest/v1 handlers | 10.8% | 70% | FAIL |
| UserService | 0% | 70% | FAIL |
| RecommendationService | 0% | 70% | FAIL |
| Category operations | 0% | 70% | FAIL |
| **Overall (internal/...)** | **11.7%** | **70%** | **FAIL** |

The overall coverage is low because large portions of the codebase (UserService, RecommendationService, all Category methods, middleware) have no tests yet. The tested methods in the three core services all meet the 70% threshold.

### Defect Summary

| Severity | Count | Description |
|----------|-------|-------------|
| Critical | 2 | BUG-001, BUG-002: missing admin role checks on Product and Category CRUD |
| High | 1 | BUG-003: race condition on concurrent purchases (no atomic transaction) |
| Medium | 2 | BUG-004: tokens not invalidated on password change; BUG-005: no login rate limiting |
| Low | 1 | BUG-006: no max-length validation on product name in service layer |
| **Total** | **6** | All documented in MIDTERM_TABLES.md |

### Planned vs Actual (Summary)

| Metric | Planned | Actual | Delta |
|--------|---------|--------|-------|
| Unit test count | 20-30 | 46 | +16 to +26 |
| Integration tests | 0 (not planned) | 5 | +5 |
| E2E tests | 10-15 | 21 | +6 to +11 |
| Defects found | 3-5 | 6 | +1 more than expected |
| Pipeline runtime | Under 5 min | ~3-5 min | On target |
| Flaky tests | 1-2 expected | 0 | Better than expected |

---

## 5. Discussion

### What Worked

**Mock-based unit testing** was highly effective. Using `testify/mock` to inject fake repositories allowed 46 tests to run in 0.62 seconds without any database. This approach is maintainable and produces zero flaky tests since all dependencies return deterministic values.

**GitHub Actions integration** worked reliably from the first push. Every commit triggers the full pipeline automatically, providing immediate feedback on regressions.

**Playwright E2E tests without a live backend** - designing tests to verify only redirect behavior and UI structure means they run in CI without spinning up the full stack. All 21 tests pass consistently.

**Concurrency tests** - implementing `TestRegister_Concurrent` and `TestPurchaseProduct_Concurrent` using `sync.WaitGroup` provided deterministic results and revealed BUG-003 conceptually, even though the mock-level test cannot replicate the actual database race condition.

### What Did Not Work

**Overall coverage is low (11.7%)** because UserService, RecommendationService, all Category service methods, and most HTTP handlers have no tests. The initial plan from A1 was too focused on the three core modules and did not schedule time for the remaining services.

**Integration tests took effort to set up** due to the need to create a `MockAuthService` that satisfies the `service.AuthService` interface at the handler level. The Handler struct requires a full `*service.Service` object, meaning all service fields must be present even if unused.

**Coverage threshold in CI** - the 70% threshold added in Midterm currently fails because untested packages (0% coverage) drag the average below 70%. In practice, the threshold needs to be applied per-module or per-package rather than globally.

### Unexpected Findings

1. **BUG-006** was discovered unexpectedly during `TestCreateProduct_VeryLongName`. The test was written to check if the service would reject a 300-character name - it did not, revealing a missing validation that could cause database truncation errors.

2. **Expired token test** required creating a second `AuthService` instance with a negative duration (`-1s`). This revealed that `time.ParseDuration("-1s")` is valid in Go and generates tokens with past expiry timestamps, which was a non-obvious but correct approach.

3. **The `done` channel pattern failed** in the initial concurrency test implementation because goroutines completed out of order - goroutine `N-1` could close the channel before other goroutines wrote their results. Replacing with `sync.WaitGroup` fixed this determinism issue.

### Improvements for Next Phase

1. **Add tests for UserService** (GetProfile, UpdateProfile, ChangePassword, DeleteAccount) - these are 0% covered and medium-risk
2. **Add tests for Category operations** - CreateCategory, UpdateCategory, DeleteCategory are all 0%
3. **Fix BUG-001 and BUG-002** - add admin role middleware and test it
4. **Address BUG-003** - add database-level locking or optimistic concurrency control for PurchaseProduct
5. **Add load testing** - use Locust or k6 to measure performance under concurrent users
6. **Tighten coverage gate** - apply 70% threshold per-package rather than globally
