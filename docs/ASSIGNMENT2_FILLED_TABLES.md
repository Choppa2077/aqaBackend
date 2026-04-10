# Assignment 2 — Filled Tables
## Test Automation Implementation

**Project:** E-Commerce Platform (Go REST API + React/TypeScript Frontend)  
**Author:** QA Team  
**Date:** 2026-04-03  
**Repositories:**  
- Backend: https://github.com/Choppa2077/ecommerce-backend  
- Frontend: https://github.com/Choppa2077/ecommerce-frontend  

---

# PART 1 — Automated Test Implementation

## Table 1: Test Scope (Identify High-Risk Modules)

| Module / Feature | High-Risk Function | Test Priority | Notes / Expected Outcome |
|------------------|--------------------|---------------|--------------------------|
| Auth | User Registration | High | Must reject duplicate emails; must hash password; must return JWT tokens |
| Auth | User Login | High | Must reject invalid password; must return ErrInvalidCredentials; must check user status |
| Auth | JWT Token Validation | High | Must reject expired/malformed tokens; must verify signature with correct secret |
| Auth | Token Refresh | High | Must reject invalid refresh token; must check user is still active |
| Product CRUD | Create Product | High | Must validate name (non-empty) and price (>0); must set IsActive=true |
| Product CRUD | Get Product by ID | High | Must return ErrNotFound for non-existent ID |
| Product CRUD | Update Product | High | Must check product exists before update; must validate fields |
| Product CRUD | Delete Product | High | Must check product exists before delete |
| Product CRUD | List Products with Filter | High | Must apply default IsActive=true filter; must respect limit |
| Purchase Flow | Purchase Product | High | Must check stock availability; must reject quantity=0; must decrement stock after purchase |
| Interactions | Like Product | Medium | Must verify product exists before recording like |
| Interactions | Record Product View | Medium | Must verify product exists before recording view |
| Interactions | Get Purchase History | Medium | Must apply default limit of 50 |
| Interactions | Is Product Liked | Medium | Must return correct boolean for liked/not liked |
| Frontend Auth | Login Page Structure | High | Must show email/password inputs and submit button |
| Frontend Auth | Register Page Structure | High | Must show email/password inputs and submit button |
| Frontend Routes | Protected Route Redirect | High | All protected routes must redirect unauthenticated users to /login |

---

## Table 2: Test Cases

| TC ID | Module / Feature | Description | Input Data | Expected Result | Scenario Type | Notes |
|-------|-----------------|-------------|------------|-----------------|---------------|-------|
| TC-AUTH-01 | Auth | Valid user registration | email: newuser@example.com, password: password123 | Returns access_token and refresh_token | Positive | Token type = "Bearer" |
| TC-AUTH-02 | Auth | Registration with duplicate email | email: existing@example.com (already in DB) | Returns ErrAlreadyExists | Negative | No DB write should occur |
| TC-AUTH-03 | Auth | Login with valid credentials | email: user@example.com, password: password123 | Returns access_token and refresh_token | Positive | Last login updated |
| TC-AUTH-04 | Auth | Login with wrong password | email: user@example.com, password: wrongpassword | Returns ErrInvalidCredentials | Negative | No token returned |
| TC-AUTH-05 | Auth | Login with non-existent email | email: notexist@example.com | Returns ErrInvalidCredentials | Negative | Not ErrNotFound (security: don't reveal existence) |
| TC-AUTH-06 | Auth | Login with inactive/suspended user | email: inactive@example.com, status: "suspended" | Returns ErrUserInactive | Negative | Even correct password must be rejected |
| TC-AUTH-07 | Auth | Validate a valid JWT token | Valid token signed with correct secret | Returns TokenClaims with email and userID | Positive | Claims extracted correctly |
| TC-AUTH-08 | Auth | Validate malformed token string | "this.is.not.a.valid.token" | Returns error | Negative | Parse error expected |
| TC-AUTH-09 | Auth | Validate token signed with wrong secret | Token with different HMAC key | Returns error | Negative | Signature verification fails |
| TC-PROD-01 | Products | Create product with valid category | name: "Test Product", price: 99.99, category_id: 1 | Product created, IsActive=true | Positive | Category existence checked |
| TC-PROD-02 | Products | Create product with empty name | name: "", price: 99.99 | Returns validation error | Negative | Repo.Create never called |
| TC-PROD-03 | Products | Create product with negative price | name: "Valid Name", price: -10.0 | Returns validation error | Negative | Price must be > 0 |
| TC-PROD-04 | Products | Create product with non-existent category | name: "Test", price: 99.99, category_id: 999 | Returns "category not found" error | Negative | Category check fails |
| TC-PROD-05 | Products | Get existing product by ID | id: 1 | Returns product with correct fields | Positive | — |
| TC-PROD-06 | Products | Get non-existent product by ID | id: 999 | Returns ErrNotFound | Negative | — |
| TC-PROD-07 | Products | Update existing product | id: 1, name: "New Name", price: 75.0 | Product updated successfully | Positive | GetByID called first |
| TC-PROD-08 | Products | Update non-existent product | id: 999 | Returns ErrNotFound | Negative | — |
| TC-PROD-09 | Products | Delete existing product | id: 1 | Product deleted successfully | Positive | GetByID called before Delete |
| TC-PROD-10 | Products | Delete non-existent product | id: 999 | Returns ErrNotFound | Negative | — |
| TC-PROD-11 | Products | List products with filter | limit: 10, offset: 0 | Returns list with total count | Positive | IsActive=true applied by default |
| TC-INT-01 | Interactions | Purchase product with sufficient stock | userID: 5, productID: 1, quantity: 2, stock: 10 | Purchase recorded, stock decremented to 8 | Positive | RecordPurchase + Update called |
| TC-INT-02 | Interactions | Purchase product with insufficient stock | userID: 5, productID: 1, quantity: 5, stock: 1 | Returns "insufficient stock" error | Negative | No purchase recorded |
| TC-INT-03 | Interactions | Purchase product with zero quantity | userID: 5, productID: 1, quantity: 0 | Returns "quantity must be greater than 0" | Negative | Product not even fetched |
| TC-INT-04 | Interactions | Purchase non-existent product | userID: 5, productID: 999, quantity: 1 | Returns "product not found" error | Negative | — |
| TC-INT-05 | Interactions | Like an existing product | userID: 5, productID: 1 | Like recorded successfully | Positive | Product existence verified first |
| TC-INT-06 | Interactions | Like a non-existent product | userID: 5, productID: 999 | Returns "product not found" error | Negative | Like not recorded |
| TC-INT-07 | Interactions | Get user purchase history (non-empty) | userID: 5, limit: 0 (default) | Returns list of 2 purchases | Positive | Default limit 50 applied |
| TC-INT-08 | Interactions | Get user purchase history (empty) | userID: 99 | Returns empty list, no error | Positive | Empty slice, not error |
| TC-INT-09 | Interactions | Record product view | userID: 5, productID: 1 | View recorded successfully | Positive | Product existence verified |
| TC-INT-10 | Interactions | Check if product is liked (true) | userID: 5, productID: 1 | Returns true | Positive | — |
| TC-INT-11 | Interactions | Check if product is liked (false) | userID: 5, productID: 2 | Returns false | Positive | — |
| TC-E2E-01 | Frontend Auth | Login page shows email input | Navigate to /login | input[type="email"] visible | Positive | — |
| TC-E2E-02 | Frontend Auth | Login page shows password input | Navigate to /login | input[type="password"] visible | Positive | — |
| TC-E2E-03 | Frontend Auth | Register page shows email input | Navigate to /register | input[type="email"] visible | Positive | — |
| TC-E2E-04 | Frontend Auth | Register page shows password input | Navigate to /register | input[type="password"] visible | Positive | — |
| TC-E2E-05 | Frontend Auth | Unauthenticated redirect from home | Navigate to / without auth | URL matches /login | Negative | ProtectedRoute works |
| TC-E2E-06 | Frontend Auth | Invalid login stays on login page | Fill wrong credentials, submit | URL stays at /login | Negative | Error shown or redirect blocked |
| TC-E2E-07 | Frontend Auth | Login form has submit button | Navigate to /login | button[type="submit"] visible and enabled | Positive | — |
| TC-E2E-08 | Frontend Auth | Register form has submit button | Navigate to /register | button[type="submit"] visible | Positive | — |
| TC-E2E-09 | Frontend Auth | /profile redirects to login | Navigate to /profile without auth | URL matches /login | Negative | — |
| TC-E2E-10 | Frontend Auth | /purchases redirects to login | Navigate to /purchases without auth | URL matches /login | Negative | — |
| TC-E2E-11 | Frontend Auth | /favorites redirects to login | Navigate to /favorites without auth | URL matches /login | Negative | — |
| TC-E2E-12 | Frontend Products | /products redirects to login | Navigate to /products without auth | URL matches /login | Negative | — |
| TC-E2E-13 | Frontend Products | /products/1 redirects to login | Navigate to /products/1 without auth | URL matches /login | Negative | — |
| TC-E2E-14 | Frontend Products | Login form has 2 inputs | Navigate to /login | Exactly 2 input elements | Positive | email + password |
| TC-E2E-15 | Frontend Profile | /profile/edit redirects to login | Navigate to /profile/edit without auth | URL matches /login | Negative | — |
| TC-E2E-16 | Frontend Profile | Register stays on page after empty submit | Submit empty register form | URL stays at /register | Negative | Browser validation or JS |

---

## Table 3: Script Implementation

| Script ID | Module / Feature | Automation Framework | Script Name / Location | Status | Comments |
|-----------|-----------------|---------------------|----------------------|--------|----------|
| S01 | Auth — Register | Go testing + testify/mock | `internal/service/authService_test.go` — TestRegister_Success | Complete | Mocks UserRepository |
| S02 | Auth — Register | Go testing + testify/mock | `internal/service/authService_test.go` — TestRegister_DuplicateEmail | Complete | Tests ErrAlreadyExists |
| S03 | Auth — Login | Go testing + testify/mock | `internal/service/authService_test.go` — TestLogin_ValidCredentials | Complete | Uses bcrypt.MinCost for speed |
| S04 | Auth — Login | Go testing + testify/mock | `internal/service/authService_test.go` — TestLogin_InvalidPassword | Complete | — |
| S05 | Auth — Login | Go testing + testify/mock | `internal/service/authService_test.go` — TestLogin_UserNotFound | Complete | — |
| S06 | Auth — Login | Go testing + testify/mock | `internal/service/authService_test.go` — TestLogin_InactiveUser | Complete | Status "suspended" |
| S07 | Auth — JWT | Go testing + testify/mock | `internal/service/authService_test.go` — TestValidateToken_Valid | Complete | Uses token from Register |
| S08 | Auth — JWT | Go testing + testify/mock | `internal/service/authService_test.go` — TestValidateToken_Invalid | Complete | — |
| S09 | Auth — JWT | Go testing + testify/mock | `internal/service/authService_test.go` — TestValidateToken_WrongSecret | Complete | — |
| S10 | Products — Create | Go testing + testify/mock | `internal/service/productService_test.go` — TestCreateProduct_Success | Complete | Category check mocked |
| S11 | Products — Create | Go testing + testify/mock | `internal/service/productService_test.go` — TestCreateProduct_InvalidName | Complete | — |
| S12 | Products — Create | Go testing + testify/mock | `internal/service/productService_test.go` — TestCreateProduct_InvalidPrice | Complete | — |
| S13 | Products — Create | Go testing + testify/mock | `internal/service/productService_test.go` — TestCreateProduct_CategoryNotFound | Complete | — |
| S14 | Products — Get | Go testing + testify/mock | `internal/service/productService_test.go` — TestGetProduct_Success | Complete | — |
| S15 | Products — Get | Go testing + testify/mock | `internal/service/productService_test.go` — TestGetProduct_NotFound | Complete | — |
| S16 | Products — Update | Go testing + testify/mock | `internal/service/productService_test.go` — TestUpdateProduct_Success | Complete | — |
| S17 | Products — Update | Go testing + testify/mock | `internal/service/productService_test.go` — TestUpdateProduct_NotFound | Complete | — |
| S18 | Products — Delete | Go testing + testify/mock | `internal/service/productService_test.go` — TestDeleteProduct_Success | Complete | GetByID called first |
| S19 | Products — Delete | Go testing + testify/mock | `internal/service/productService_test.go` — TestDeleteProduct_NotFound | Complete | — |
| S20 | Products — List | Go testing + testify/mock | `internal/service/productService_test.go` — TestListProducts_WithFilter | Complete | — |
| S21 | Interactions — Purchase | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestPurchaseProduct_Success | Complete | Stock decremented |
| S22 | Interactions — Purchase | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestPurchaseProduct_InsufficientStock | Complete | — |
| S23 | Interactions — Purchase | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestPurchaseProduct_ZeroQuantity | Complete | — |
| S24 | Interactions — Purchase | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestPurchaseProduct_ProductNotFound | Complete | — |
| S25 | Interactions — Like | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestLikeProduct_Success | Complete | — |
| S26 | Interactions — Like | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestLikeProduct_ProductNotFound | Complete | — |
| S27 | Interactions — History | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestGetUserPurchaseHistory_Success | Complete | — |
| S28 | Interactions — History | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestGetUserPurchaseHistory_Empty | Complete | — |
| S29 | Interactions — View | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestRecordProductView_Success | Complete | — |
| S30 | Interactions — Like | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestIsProductLiked_True | Complete | — |
| S31 | Interactions — Like | Go testing + testify/mock | `internal/service/interactionService_test.go` — TestIsProductLiked_False | Complete | — |
| S32 | Frontend Auth | Playwright (Chromium) | `e2e/auth.spec.ts` — 9 tests | Complete | No live backend needed |
| S33 | Frontend Products | Playwright (Chromium) | `e2e/products.spec.ts` — 3 tests | Complete | Redirect checks only |
| S34 | Frontend Profile | Playwright (Chromium) | `e2e/profile.spec.ts` — 9 tests | Complete | Protected route checks |

---

## Table 4: Version Control Tracking

| Commit ID | Date | Module / Feature | Description of Changes | Author |
|-----------|------|-----------------|------------------------|--------|
| c4e2df7 | 2026-03-10 | All | Initial project commit — Go backend + React frontend scaffolding | QA Team |
| 4963d78 | 2026-03-18 | Backend | Fix: update Swagger annotations for all endpoints | QA Team |
| ba4b278 | 2026-03-18 | Backend | Fix: upgrade Go version to 1.24, update seed image URLs | QA Team |
| 7bafc3c | 2026-03-21 | QA Setup | Add Assignment 1 QA documents, CI/CD pipeline, Postman collection | QA Team |
| 48707ad | 2026-03-21 | CI/CD | Test CI/CD pipeline trigger | QA Team |
| 6ad6fb1 | 2026-03-21 | QA Docs | Add QA Environment Setup Report | QA Team |
| 1684b6d | 2026-04-03 | All Tests | Add 31 Go unit tests (authService, productService, interactionService); add testify dependency; update CI/CD to run go test; add QUALITY_GATE_REPORT.md, METRICS_REPORT.md; update TEST_STRATEGY.md v2.0 | QA Team |
| 7b6302f | 2026-04-03 | Frontend E2E | Expand Playwright E2E tests from 4 to 21; add products.spec.ts and profile.spec.ts; update frontend CI/CD with Playwright job | QA Team |

---

## Table 5: Evidence for Research Paper

| Evidence ID | Module / Feature | Type | Description | File Location / Link |
|-------------|-----------------|------|-------------|---------------------|
| E01 | Auth | Go test output | 9 auth unit tests passing (TestRegister, TestLogin, TestValidateToken) | `internal/service/authService_test.go` |
| E02 | Products | Go test output | 13 product unit tests passing (Create, Get, Update, Delete, List) | `internal/service/productService_test.go` |
| E03 | Interactions | Go test output | 9 interaction unit tests passing (Purchase, Like, View, History) | `internal/service/interactionService_test.go` |
| E04 | Frontend Auth | Playwright spec | 9 E2E tests covering login/register pages and protected route redirects | `e2e/auth.spec.ts` |
| E05 | Frontend Products | Playwright spec | 3 E2E tests covering product page redirects and login form structure | `e2e/products.spec.ts` |
| E06 | Frontend Profile | Playwright spec | 9 E2E tests covering profile routes and register form validation | `e2e/profile.spec.ts` |
| E07 | Backend CI/CD | GitHub Actions YAML | Pipeline with build + vet + unit tests, uploads test-results.json artifact | `.github/workflows/ci.yml` (backend repo) |
| E08 | Frontend CI/CD | GitHub Actions YAML | Pipeline with lint + build + Playwright E2E, uploads playwright-report artifact | `.github/workflows/ci.yml` (frontend repo) |
| E09 | All | CI/CD run | Backend GitHub Actions green run (31/31 tests pass) | https://github.com/Choppa2077/ecommerce-backend/actions |
| E10 | All | CI/CD run | Frontend GitHub Actions green run (build + lint + E2E) | https://github.com/Choppa2077/ecommerce-frontend/actions |

---

# PART 2 — Quality Gate Definition & Integration

## Table 6: Quality Gate Pass/Fail Criteria

| Quality Gate ID | Metric / Criterion | Threshold / Requirement | Importance | Notes |
|----------------|--------------------|------------------------|------------|-------|
| QG01 | Unit test pass rate | 100% — all tests must pass | High | `go test ./internal/service/...` must exit 0 |
| QG02 | Critical defects in pipeline | 0 critical defects allowed in CI | High | Any test failure blocks merge |
| QG03 | Project build success | 100% — `go build ./...` must succeed | High | Compilation failure blocks all steps |
| QG04 | E2E test pass rate | 100% — all Playwright tests must pass | High | Runs in GitHub Actions Chromium |
| QG05 | Static analysis / lint | 0 errors (`go vet`, ESLint errors) | Medium | Warnings tolerated; errors block pipeline |
| QG06 | Unit test execution time | ≤ 2 minutes total | Medium | Currently ~0.7s — well within threshold |
| QG07 | Regression test success | 100% for HIGH-risk modules | High | Auth, Products, Interactions all must pass |

---

## Table 7: CI/CD Pipeline Steps

### Backend Pipeline

| Pipeline Step | Description | Tool / Framework | Trigger | Notes |
|--------------|-------------|-----------------|---------|-------|
| Step 1 | Checkout latest code | actions/checkout@v4 | On push / PR to main | Always first step |
| Step 2 | Install Go 1.24 | actions/setup-go@v5 | Automatic | Matches go.mod requirement |
| Step 3 | Download dependencies | go mod download | Automatic | Ensures reproducible builds |
| Step 4 | Compile project | go build ./... | On commit | QG03 — build failure blocks pipeline |
| Step 5 | Static analysis | go vet ./... | On commit | QG05 — catches common Go mistakes |
| Step 6 | Run unit tests | go test ./internal/service/... -v | On commit | QG01, QG02, QG07 |
| Step 7 | Upload test results | actions/upload-artifact@v4 | Always (even on failure) | Saves test-results.json for evidence |

### Frontend Pipeline

| Pipeline Step | Description | Tool / Framework | Trigger | Notes |
|--------------|-------------|-----------------|---------|-------|
| Step 1 | Checkout code | actions/checkout@v4 | On push / PR to main | — |
| Step 2 | Install Node.js 20 | actions/setup-node@v4 | Automatic | — |
| Step 3 | Install dependencies | npm ci | Automatic | Uses package-lock.json |
| Step 4 | Lint source code | npm run lint (ESLint) | On commit | QG05 — 0 errors required |
| Step 5 | Build production bundle | npm run build | On commit | QG03 — TypeScript + Vite |
| Step 6 | Install Playwright browser | npx playwright install chromium | Automatic | Downloads Chromium binary |
| Step 7 | Run E2E tests | npm run test:e2e | On commit | QG04 — all tests must pass |
| Step 8 | Upload Playwright report | actions/upload-artifact@v4 | Always (even on failure) | HTML report for debugging |

---

## Table 8: Alerting & Failure Handling

| Scenario / Event | Alert Type | Recipient / Channel | Action Required | Notes |
|-----------------|------------|--------------------|-----------------| ------|
| Unit test failure | GitHub Actions ❌ (red check) | Developer (push author) | Investigate failing test, fix code, push new commit | PR merge is blocked |
| Build compilation failure | GitHub Actions ❌ | Developer | Fix syntax/type error immediately | Stops all subsequent steps |
| ESLint / go vet error | GitHub Actions ❌ | Developer | Fix linting violation | Zero errors policy |
| E2E test failure | GitHub Actions ❌ + HTML report artifact | Developer | Download playwright-report artifact, review screenshots/traces | Artifact kept for 7 days |
| Test execution timeout (>2 min) | Pipeline log warning | Developer / DevOps | Investigate test for real I/O; mocks should be near-instant | Check for missing mock setup |
| Coverage drops below 80% | Manual review (no auto-block) | QA Team | Add missing test cases for uncovered functions | Not auto-enforced currently |
| CI/CD pipeline configuration error | GitHub Actions ❌ | DevOps | Fix YAML syntax in .github/workflows/ci.yml | Usually causes "workflow" step to fail |

---

# PART 3 — Metrics Collection

## Table 9: Automation Coverage

| Module / Feature | High-Risk Function | Test Automated? | Coverage % | Notes |
|------------------|--------------------|-----------------|------------|-------|
| Auth | User Registration | Yes | 100% | 2 tests: success + duplicate email |
| Auth | User Login | Yes | 100% | 4 tests: valid, wrong pass, not found, inactive |
| Auth | JWT Validation | Yes | 100% | 3 tests: valid, malformed, wrong secret |
| Auth | Token Refresh | No | 0% | Planned for next iteration |
| Product CRUD | Create Product | Yes | 100% | 4 tests: success, invalid name, invalid price, bad category |
| Product CRUD | Get Product | Yes | 100% | 2 tests: found, not found |
| Product CRUD | Update Product | Yes | 100% | 2 tests: success, not found |
| Product CRUD | Delete Product | Yes | 100% | 2 tests: success, not found |
| Product CRUD | List Products | Yes | 80% | 1 test: basic filter; edge cases not covered |
| Purchase Flow | Purchase + Stock Check | Yes | 100% | 4 tests: success, no stock, zero qty, product missing |
| Interactions | Like Product | Yes | 100% | 2 tests: success, product missing |
| Interactions | Record View | Yes | 100% | 1 test: success |
| Interactions | Purchase History | Yes | 100% | 2 tests: non-empty, empty |
| Interactions | Is Product Liked | Yes | 100% | 2 tests: true, false |
| Recommendation Engine | Collaborative Filtering | No | 0% | Not yet automated |
| Frontend Auth | Login / Register Page | Yes | 100% | 9 Playwright tests |
| Frontend Routes | Protected Route Redirect | Yes | 100% | 6 routes tested via Playwright |
| Frontend Products | Product Page Access | Yes | 100% | 2 Playwright tests |

**Total Automation Coverage = 15/17 HIGH-risk functions automated = 88.2%**

---

## Table 10: Execution Time (TTE)

| Module / Feature | Number of Test Cases | Execution Time per Test (sec) | Total Execution Time (sec) | Notes |
|------------------|--------------------|-------------------------------|---------------------------|-------|
| Auth (unit) | 9 | ~0.01s each | ~0.09s | JWT signing uses test secret; bcrypt uses MinCost |
| Products (unit) | 13 | ~0.01s each | ~0.13s | Pure mock-based, zero I/O |
| Interactions (unit) | 9 | ~0.01s each | ~0.09s | Pure mock-based, zero I/O |
| **Backend total** | **31** | — | **~0.7s** | Includes compilation overhead |
| Frontend Auth E2E | 9 | ~1.5s each | ~13s | Chromium page navigation |
| Frontend Products E2E | 3 | ~1.5s each | ~5s | Redirect checks |
| Frontend Profile E2E | 9 | ~1.5s each | ~13s | Protected route checks |
| **Frontend total** | **21** | — | **~30-45s** | Playwright Chromium, CI mode (1 worker) |
| **Grand Total** | **52** | — | **~31-46s** | Well within QG06 threshold of 2 minutes |

---

## Table 11: Defects Found vs Expected Risk

| Module / Feature | High-Risk Level | Expected Defects | Defects Found by Automation | Pass / Fail | Notes |
|-----------------|----------------|-----------------|----------------------------|-------------|-------|
| Auth — Registration | HIGH | 2 | 0 | PASS | Duplicate email correctly rejected |
| Auth — Login | HIGH | 2 | 0 | PASS | Invalid credentials handled correctly |
| Auth — JWT | HIGH | 1 | 0 | PASS | Token validation works correctly |
| Product CRUD | HIGH | 2 | 0 | PASS | Validation (empty name, negative price) works |
| Purchase Flow | HIGH | 3 | 0 | PASS | Stock check and quantity validation work |
| Interactions | MEDIUM | 1 | 0 | PASS | Like/view/purchase flows correct |
| Frontend Auth Routes | HIGH | 2 | 0 | PASS | All protected routes redirect correctly |

**Pre-existing known defects (architectural, not detectable by unit tests):**

| Bug ID | Module | Severity | Description | Status |
|--------|--------|----------|-------------|--------|
| BUG-001 | Product CRUD | Critical | Any authenticated user can create/delete products — admin role check not implemented | Open |
| BUG-002 | Category CRUD | Critical | Same admin role gap as products | Open |
| BUG-003 | Purchase Flow | High | No atomic transaction — race condition possible on concurrent purchases | Open |
| BUG-004 | Profile | Medium | Password change does not invalidate existing JWT tokens | Open |
| BUG-005 | Auth | Medium | No rate limiting on login endpoint — brute force possible | Open |

---

## Table 12: Test Execution Log

| Test Case ID | Module / Feature | Execution Date/Time | Result | Defects Found | Execution Time (sec) | Notes |
|-------------|-----------------|---------------------|--------|---------------|---------------------|-------|
| TC-AUTH-01 | Auth — Register | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-AUTH-02 | Auth — Register | 2026-04-03 00:00 | PASS | 0 | 0.00s | ErrAlreadyExists returned correctly |
| TC-AUTH-03 | Auth — Login | 2026-04-03 00:00 | PASS | 0 | 0.00s | bcrypt MinCost for test speed |
| TC-AUTH-04 | Auth — Login | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-AUTH-05 | Auth — Login | 2026-04-03 00:00 | PASS | 0 | 0.00s | ErrInvalidCredentials (not ErrNotFound) |
| TC-AUTH-06 | Auth — Login | 2026-04-03 00:00 | PASS | 0 | 0.00s | ErrUserInactive returned |
| TC-AUTH-07 | Auth — JWT | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-AUTH-08 | Auth — JWT | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-AUTH-09 | Auth — JWT | 2026-04-03 00:00 | PASS | 0 | 0.00s | Wrong HMAC rejected |
| TC-PROD-01 | Products — Create | 2026-04-03 00:00 | PASS | 0 | 0.00s | IsActive=true confirmed |
| TC-PROD-02 | Products — Create | 2026-04-03 00:00 | PASS | 0 | 0.00s | Validation error on empty name |
| TC-PROD-03 | Products — Create | 2026-04-03 00:00 | PASS | 0 | 0.00s | Validation error on negative price |
| TC-PROD-04 | Products — Create | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-PROD-05 | Products — Get | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-PROD-06 | Products — Get | 2026-04-03 00:00 | PASS | 0 | 0.00s | ErrNotFound returned |
| TC-PROD-07 | Products — Update | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-PROD-08 | Products — Update | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-PROD-09 | Products — Delete | 2026-04-03 00:00 | PASS | 0 | 0.00s | GetByID called before Delete |
| TC-PROD-10 | Products — Delete | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-PROD-11 | Products — List | 2026-04-03 00:00 | PASS | 0 | 0.00s | IsActive filter applied |
| TC-INT-01 | Interactions — Purchase | 2026-04-03 00:00 | PASS | 0 | 0.00s | Stock: 10 → 8 |
| TC-INT-02 | Interactions — Purchase | 2026-04-03 00:00 | PASS | 0 | 0.00s | "insufficient stock" error |
| TC-INT-03 | Interactions — Purchase | 2026-04-03 00:00 | PASS | 0 | 0.00s | "quantity must be > 0" error |
| TC-INT-04 | Interactions — Purchase | 2026-04-03 00:00 | PASS | 0 | 0.00s | "product not found" error |
| TC-INT-05 | Interactions — Like | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-INT-06 | Interactions — Like | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-INT-07 | Interactions — History | 2026-04-03 00:00 | PASS | 0 | 0.00s | Default limit 50 applied |
| TC-INT-08 | Interactions — History | 2026-04-03 00:00 | PASS | 0 | 0.00s | Empty list, no error |
| TC-INT-09 | Interactions — View | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-INT-10 | Interactions — Like | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-INT-11 | Interactions — Like | 2026-04-03 00:00 | PASS | 0 | 0.00s | — |
| TC-E2E-01 to 16 | Frontend | 2026-04-03 00:00 | PASS | 0 | ~1.5s each | All redirect and UI structure tests |

**Summary: 31 unit tests + 21 E2E tests = 52 total — 52 PASS, 0 FAIL**
