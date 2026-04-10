# Midterm Project - Required Tables
## QA Implementation and Empirical Analysis

**Project:** E-Commerce Platform (Go REST API + React/TypeScript Frontend)
**Author:** QA Team
**Date:** 2026-04-10
**Repository:** https://github.com/Choppa2077/aqaBackend

---

# Task 1.1 - Risk Re-evaluation Table

| Module | Original Risk Score (A1) | Observed Issues (A2) | Updated Risk Score | Justification |
|--------|--------------------------|----------------------|--------------------|---------------|
| Auth (Register/Login/JWT) | 9/10 | 0 failures in 15 unit tests; all edge cases passed including expired token, empty email, concurrent registration | 8/10 | Score slightly reduced because automation confirmed core logic is solid. Risk remains HIGH due to BUG-005 (no rate limiting on login - brute force possible). Coverage of authService: 75% average across functions. |
| Product CRUD | 8/10 | 0 failures in 19 unit tests; new UpdateStock and CheckStock tests passed; VeryLongName edge case passed without error (no max-length validation in service layer) | 8/10 | Score unchanged. Discovering that service has no max-length validation is an unexpected gap. Category operations (CreateCategory, UpdateCategory, DeleteCategory) remain at 0% coverage - this is an active risk. |
| Purchase Flow (Stock Check) | 9/10 | 0 failures in 4 unit tests; concurrent purchase test confirmed race condition is a real risk (only mock-level, real DB has no atomic lock) | 9/10 | Score unchanged. Concurrency test revealed that PurchaseProduct has no atomic transaction - BUG-003 confirmed as active. This is a HIGH risk in production. |
| Frontend Auth Routes | 7/10 | 0 failures in 9 Playwright E2E tests; all 6 protected routes redirect correctly | 6/10 | Score reduced. Automation confirmed all redirect logic works correctly. Risk is now lower since routes are verified end-to-end via CI/CD. |
| Recommendation Engine | 6/10 | 0 automated tests exist; 0% coverage confirmed | 7/10 | Score increased. No tests were added in A2 or Midterm. This untested module processes user interaction data and if broken, affects core product value. Detectability is very low. |

---

# Task 1.2A - Failed Test Cases

| Test Name / ID | Module | Failure Type | Frequency | Status |
|----------------|--------|--------------|-----------|--------|
| No failures recorded | - | - | - | All 51 tests pass consistently across 5+ CI/CD runs |

**Note:** Zero test failures were observed across all automation runs. All 46 unit tests and 5 integration tests pass on every push to main branch. This confirms the core business logic is implemented correctly for tested scenarios.

---

# Task 1.2B - Flaky Tests (Instability Analysis)

| Test Name | Module | Passes | Failures | Suspected Cause | Status |
|-----------|--------|--------|----------|-----------------|--------|
| No flaky tests identified | - | - | - | - | Stable |

**Analysis:** All 51 tests have been executed across multiple CI/CD pipeline runs with consistent PASS results. Tests use mock repositories with deterministic return values, eliminating timing-based flakiness. The concurrency tests (TestRegister_Concurrent, TestPurchaseProduct_Concurrent) use sync.WaitGroup for deterministic completion - no race detector violations observed.

---

# Task 1.2C - Coverage Gaps

| Module / Function | Coverage % | Status | Risk Level |
|-------------------|------------|--------|------------|
| authService - Register | 87.5% | Good | Low |
| authService - Login | 83.3% | Good | Low |
| authService - ValidateToken | 72.2% | Acceptable | Medium |
| authService - RefreshToken | 64.3% | Below threshold | HIGH |
| productService - CreateProduct | 91.7% | Good | Low |
| productService - GetProduct | 100% | Full | None |
| productService - DeleteProduct | 100% | Full | None |
| productService - UpdateProduct | 50.0% | Below threshold | HIGH |
| productService - UpdateStock | 87.5% | Good | Low |
| productService - CheckStock | 75.0% | Acceptable | Medium |
| productService - CreateCategory | 0% | Not covered | CRITICAL |
| productService - UpdateCategory | 0% | Not covered | CRITICAL |
| productService - DeleteCategory | 0% | Not covered | CRITICAL |
| productService - SearchProducts | 0% | Not covered | HIGH |
| productService - GetProductStatistics | 0% | Not covered | HIGH |
| interactionService - PurchaseProduct | 80.0% | Good | Low |
| interactionService - LikeProduct | 75.0% | Acceptable | Medium |
| interactionService - UnlikeProduct | 40.0% | Below threshold | HIGH |
| interactionService - GetUserViewHistory | 0% | Not covered | HIGH |
| interactionService - GetUserInteractionSummary | 0% | Not covered | HIGH |
| delivery/rest/v1 - Auth handlers | 10.8% | Below threshold | HIGH |
| UserService (all methods) | 0% | Not covered | CRITICAL |
| RecommendationService (all methods) | 0% | Not covered | CRITICAL |

**Modules with coverage below 70%:** RefreshToken (64.3%), UpdateProduct (50.0%), UnlikeProduct (40.0%), all Category operations (0%), UserService (0%), RecommendationService (0%), delivery handlers (10.8%).

**Overall service-layer coverage: 26.4% of statements**

---

# Task 1.2D - Unexpected System Behavior

| ID | Module | Type | Description | Discovered In |
|----|--------|------|-------------|---------------|
| BUG-001 | Product CRUD | Security - Critical | Any authenticated user can create/update/delete products. No admin role check exists in service or handler layer. | Code review during A2 |
| BUG-002 | Category CRUD | Security - Critical | Same admin role gap as products - any user can manage categories | Code review during A2 |
| BUG-003 | Purchase Flow | Race Condition - High | PurchaseProduct has no atomic transaction. Two concurrent buyers of the last item can both pass the stock check before either decrements. TestPurchaseProduct_Concurrent confirmed this at mock level. | Midterm concurrency test |
| BUG-004 | Auth | Security - Medium | Password change (ChangePassword in UserService) does not invalidate existing JWT tokens. Old tokens remain valid until natural expiry. | Code review during A2 |
| BUG-005 | Auth | Security - Medium | No rate limiting on /auth/login endpoint. Unlimited brute-force attempts are possible. | Code review during A2 |
| BUG-006 | Product | Validation Gap | Service layer has no maximum length validation on product name. A 300-character name passes CreateProduct without error. TestCreateProduct_VeryLongName confirmed this unexpectedly passed. | Midterm edge case test |

---

# Task 1.3 - Risk Dimensions Matrix

| Module | Likelihood | Impact | Detectability | Overall Risk | Notes |
|--------|-----------|--------|---------------|--------------|-------|
| Auth (Register/Login/JWT) | Low (2/5) | Critical (5/5) | High (4/5) | HIGH | 15 tests, 75% avg coverage; main risk is BUG-005 (brute force) |
| Product CRUD | Low (2/5) | High (4/5) | Medium (3/5) | HIGH | Category ops at 0% coverage lower detectability |
| Purchase Flow | Medium (3/5) | Critical (5/5) | Low (2/5) | CRITICAL | BUG-003 race condition not detectable via current mock tests |
| Frontend Auth Routes | Low (1/5) | Medium (3/5) | High (5/5) | LOW | All 9 E2E tests pass; routes verified in CI/CD |
| Recommendation Engine | High (4/5) | High (4/5) | Very Low (1/5) | CRITICAL | 0% test coverage; complex algorithm; no monitoring |
| UserService | Medium (3/5) | High (4/5) | Very Low (1/5) | HIGH | 0% coverage; password change, profile update untested |
| Category Management | Low (2/5) | Medium (3/5) | Very Low (1/5) | HIGH | 0% coverage on all 5 category service methods |

---

# Task 2.1 - New Test Cases (Midterm Extensions)

| TC ID | Module | Scenario Type | Input Data | Expected Output | Actual Result | Status |
|-------|--------|---------------|------------|-----------------|---------------|--------|
| TC-EDGE-01 | Auth - Login | Invalid Input | email: "", password: "password123" | ErrInvalidCredentials | ErrInvalidCredentials returned | PASS |
| TC-EDGE-02 | Auth - Register | Invalid Input | email: "edge@example.com", PasswordHash: "" - repo returns error | Error propagated from repo | Error returned | PASS |
| TC-EDGE-03 | Auth - JWT | Failure Scenario | Token generated with -1s duration (already expired at creation) | Error on ValidateToken | Error returned | PASS |
| TC-EDGE-04 | Auth - Refresh | Positive | Valid refresh token + active user in DB | New access token returned | New token returned | PASS |
| TC-EDGE-05 | Auth - Refresh | Failure Scenario | "this.is.invalid.refresh.token" | Error returned | Error returned | PASS |
| TC-EDGE-06 | Auth - Register | Concurrency | 10 goroutines register same email; goroutine 0 sees free slot, 1-9 see taken slot | 1 success, 9 ErrAlreadyExists | 1 success, 9 ErrAlreadyExists | PASS |
| TC-EDGE-07 | Product - Stock | Positive | productID: 1, stock: 10, add quantity: +5 | Stock updated to 15 | Stock updated to 15 | PASS |
| TC-EDGE-08 | Product - Stock | Failure Scenario | productID: 999, quantity: 5 | ErrNotFound | ErrNotFound returned | PASS |
| TC-EDGE-09 | Product - CheckStock | Positive (Edge) | productID: 1, stock: 20, requested: 10 | Returns true | true returned | PASS |
| TC-EDGE-10 | Product - CheckStock | Edge Case | productID: 1, stock: 3, requested: 10 | Returns false (insufficient) | false returned | PASS |
| TC-EDGE-11 | Product - Create | Edge Case (Large Input) | name: 300 chars, price: 9.99 | Product created (no max-length in service) | Created successfully | PASS |
| TC-EDGE-12 | Purchase - Concurrent | Concurrency / Race Condition | 5 goroutines buy last item; goroutine 0 has stock=1, others have stock=0 | 1 success, 4 failures | 1 success, 4 failures | PASS |
| TC-EDGE-13 | Interaction - Unlike | Positive | userID: 5, productID: 1, RemoveLike succeeds | No error | No error | PASS |
| TC-INT-01 | Auth Handler | Integration - Positive | POST /auth/register with valid JSON body | HTTP 201 + access_token in response | HTTP 201 returned | PASS |
| TC-INT-02 | Auth Handler | Integration - Invalid Input | POST /auth/register with malformed JSON | HTTP 400 | HTTP 400 returned | PASS |
| TC-INT-03 | Auth Handler | Integration - Positive | POST /auth/login with valid credentials | HTTP 200 + access_token + refresh_token | HTTP 200 returned | PASS |
| TC-INT-04 | Auth Handler | Integration - Failure | POST /auth/login with wrong password | HTTP 401 + error message | HTTP 401 returned | PASS |
| TC-INT-05 | Auth Handler | Integration - Invalid Input | POST /auth/login missing password field | HTTP 400 | HTTP 400 returned | PASS |

---

# Task 2.4 - Quality Gates

| QG ID | Metric | Threshold | Actual Value | Status | Critical Analysis |
|-------|--------|-----------|--------------|--------|-------------------|
| QG01 | Unit test pass rate | 100% | 46/46 = 100% | PASS | Threshold appropriate. All business logic tests pass. |
| QG02 | Integration test pass rate | 100% | 5/5 = 100% | PASS | Threshold appropriate. HTTP layer verified end-to-end. |
| QG03 | Critical test failures | 0 | 0 | PASS | Threshold appropriate. Zero tolerance for critical failures is correct. |
| QG04 | Build success | 100% | Passing | PASS | Non-negotiable gate. Any compilation failure blocks everything. |
| QG05 | Static analysis (go vet) | 0 errors | 0 errors | PASS | Threshold appropriate. Go vet catches common mistakes. |
| QG06 | Unit test execution time | Less than 2 minutes | 0.62s (unit) + 0.55s (integration) = ~1.2s total | PASS | Threshold too lenient - current tests run in under 2 seconds. Could tighten to 30s. |
| QG07 | Service layer coverage | 70% per high-risk module | Auth: 75% avg, Products (tested methods): 82% avg, Interactions (tested): 72% avg | PARTIAL | Coverage gate was not enforced in CI/CD in A2. Added in Midterm. Many untested methods (Category, UserService, Recommendation) bring overall coverage to 26.4%. Threshold is too lenient if applied per-method; needs to be tightened in next iteration. |

**Critical Analysis:**
- QG06 threshold of 2 minutes is too lenient for a test suite that runs in under 2 seconds. Should be tightened to 30 seconds.
- QG07 was not enforced in CI/CD pipeline during A2. Added in Midterm. Overall 26.4% is below 70% when considering untested services. The gate threshold needs to either exclude untested packages or be applied per-module rather than overall.
- No failures were caused by poor code quality. Gaps are due to insufficient tests (category operations, UserService, RecommendationService not yet covered).

---

# Task 3.1 - Coverage Metrics

| Module | Functions Tested | Functions Total | Coverage % | Tool |
|--------|-----------------|-----------------|------------|------|
| authService | 4/4 | 4 | 75% avg | go test -cover |
| productService (core CRUD) | 7/18 | 18 | 82% avg (tested only) | go test -cover |
| interactionService (core) | 8/10 | 10 | 72% avg (tested only) | go test -cover |
| delivery/rest/v1 (auth handlers) | 2/3 handlers | ~20 handlers | 10.8% | go test -cover |
| UserService | 0/4 | 4 | 0% | go test -cover |
| RecommendationService | 0/2 | 2 | 0% | go test -cover |
| Overall internal/... | - | - | 11.7% statements | go test -cover |

---

# Task 3.1 - Defect Detection

| Bug ID | Module | Severity | Description | Risk Level | Detected By |
|--------|--------|----------|-------------|------------|-------------|
| BUG-001 | Product CRUD | Critical | Any user can create/delete products (no admin check) | HIGH | Code review |
| BUG-002 | Category CRUD | Critical | Same admin gap for categories | HIGH | Code review |
| BUG-003 | Purchase Flow | High | No atomic transaction - race condition on concurrent purchases | CRITICAL | Midterm concurrency test |
| BUG-004 | Auth | Medium | Password change does not invalidate JWT tokens | MEDIUM | Code review |
| BUG-005 | Auth | Medium | No rate limiting on /auth/login | MEDIUM | Code review |
| BUG-006 | Product | Low | No max-length validation on product name in service layer | LOW | Midterm edge case test (TC-EDGE-11) |

**Total defects found: 6 (2 Critical, 1 High, 2 Medium, 1 Low)**

---

# Task 3.1 - Efficiency Metrics (Execution Time)

| Phase | Test Count | Execution Time | Notes |
|-------|------------|----------------|-------|
| A2 - Unit tests (service layer only) | 31 | ~0.7s | No integration tests |
| Midterm - Unit tests (service layer) | 46 | ~0.62s | +15 new tests, still under 1s |
| Midterm - Integration tests (handlers) | 5 | ~0.55s | First time running handler tests |
| Midterm - Total (unit + integration) | 51 | ~1.17s | Well within QG06 threshold |
| Frontend E2E (Playwright, CI) | 21 | ~30-45s | Unchanged from A2 |
| Full pipeline runtime (CI/CD) | - | ~3-5 minutes | Build + vet + tests + coverage |

**Improvement:** A2 had 31 tests in 0.7s. Midterm added 20 new tests with total time remaining under 1.2s. Pipeline efficiency is well maintained.

---

# Task 4.1 - Planned vs Actual

| Aspect | Planned (A1) | Actual (A2 + Midterm) | Gap |
|--------|-------------|----------------------|-----|
| Test count | 20-30 unit tests | 46 unit + 5 integration + 21 E2E = 72 total | Exceeded plan by 2.4x |
| Coverage target | 80% for HIGH-risk modules | 75-92% for tested methods; 0% for untested modules (UserService, RecommendationService, Category ops) | Partially met - tested modules meet target, but scope too narrow |
| Integration tests | Not planned explicitly | 5 handler-level integration tests added in Midterm | Gap filled in Midterm |
| Concurrency tests | Not planned in A1 | 2 concurrency tests (TestRegister_Concurrent, TestPurchaseProduct_Concurrent) added in Midterm | Gap filled in Midterm |
| CI/CD quality gates | 5 gates defined | 7 gates defined + coverage threshold added in Midterm | Exceeded plan |
| Flaky tests | Expected 1-2 | 0 flaky tests observed across 5+ runs | Better than expected |
| Defects found | Expected 3-5 | 6 defects found (2 critical security gaps, 1 race condition, 2 security medium, 1 validation) | More security issues than expected |
| Test execution time | Under 5 minutes | 1.2s unit + integration; 3-5min full pipeline | Well within target |
| E2E test count | 10-15 | 21 Playwright tests | Exceeded plan |

**Required Insights:**

1. **Incorrect assumptions in planning:** A1 assumed that unit tests alone would achieve 80% coverage per module. In reality, many service methods (CreateCategory, SearchProducts, UserService, RecommendationService) were not covered because they were deprioritized as "medium risk." The actual coverage of the entire service package is only 26.4%.

2. **Missing test scenarios in A1 plan:** Concurrency and race condition tests were not planned in A1. The Midterm requirement to add these revealed BUG-003 (purchase race condition). Edge case testing (very long names, empty fields at service level) also found BUG-006 (no max-length validation).

3. **Inefficient automation design:** The initial test structure covered only the "happy path" and basic negative cases for each service. Edge cases and boundary conditions were added only in Midterm as a requirement. A better initial design would have included edge cases from the start.
