# QA Test Strategy Document

**Project:** E-Commerce Platform
**Version:** 1.0
**Date:** 2026-03-21
**Author:** QA Team

---

## 1. Project Scope & Objectives

### System Under Test

A full-stack e-commerce platform consisting of:

| Component | Technology | Repository |
|-----------|-----------|------------|
| **Backend API** | Go 1.24, Gin, MongoDB | github.com/Choppa2077/ecommerce-backend |
| **Frontend** | React 19, TypeScript, Vite | github.com/Choppa2077/ecommerce-frontend |

### Scope

**In scope:**
- All 28 REST API endpoints (auth, products, categories, profiles)
- Core user flows on frontend (login, browse, purchase, profile)
- JWT authentication and token refresh mechanism
- Purchase flow and stock management
- Recommendation engine correctness

**Out of scope:**
- Performance/load testing (future assignment)
- Security penetration testing (future assignment)
- Mobile responsiveness testing

### Objectives

1. Establish a complete test coverage baseline from 0%
2. Identify and document defects in HIGH-risk modules
3. Set up automated testing infrastructure for CI/CD
4. Provide reproducible test environment for future assignments

---

## 2. Risk Assessment Summary

See full document: [RISK_ASSESSMENT.md](./RISK_ASSESSMENT.md)

| Priority | Count | Key Modules |
|----------|-------|-------------|
| 🔴 HIGH | 5 | Auth, Token Refresh, Product CRUD, Purchase Flow, Recommendations |
| 🟡 MEDIUM | 3 | Profile, Categories, Product Filtering |
| 🟢 LOW | 1 | Like/View Interactions |

**Critical known defect:** Admin role check is not implemented (`// TODO` in source code) — all authenticated users can create/delete products.

---

## 3. Test Approach

### 3.1 Testing Levels

| Level | Type | Tool | Coverage Target |
|-------|------|------|----------------|
| **API Testing** | Manual + Automated | Postman | All 28 endpoints |
| **Unit Testing** | Automated | Go `testing` + `testify` | HIGH-risk services |
| **E2E Testing** | Automated | Playwright | Critical user flows |
| **Smoke Testing** | Manual | Postman / Browser | Core functionality |

### 3.2 Priority Order

Testing is prioritized by risk score (HIGH first):

```
1. Authentication (JWT) ────────────── Score: 15
2. Product CRUD ────────────────────── Score: 12
3. Recommendation Engine ───────────── Score: 12
4. Token Refresh Interceptor ───────── Score: 12
5. Purchase Flow + Stock ───────────── Score: 10
6. Profile Management ──────────────── Score: 6
7. Category Management ─────────────── Score: 6
8. Product Filtering ───────────────── Score: 6
9. Like/View Interactions ──────────── Score: 4
```

### 3.3 Test Approach by Module

#### Auth Module (HIGH)
- **Manual:** Postman — register, login, refresh with valid/invalid inputs
- **Automated:** Unit test `authService.go` — Register(), Login(), ValidateToken(), RefreshToken()
- **E2E:** Playwright — login flow, redirect on invalid credentials, token persistence

#### Product CRUD (HIGH)
- **Manual:** Postman — CRUD operations, verify admin-less access (known bug)
- **Automated:** Unit test `productService.go` — CreateProduct(), GetProduct(), UpdateProduct()
- **E2E:** Playwright — browse products, filter by category/price

#### Purchase Flow (HIGH)
- **Manual:** Postman — purchase with valid/invalid quantity, verify stock decrement
- **Automated:** Unit test `interactionService.go` — PurchaseProduct(), stock validation
- **E2E:** Playwright — full purchase flow from product page

#### Recommendation Engine (HIGH)
- **Manual:** Postman — hit endpoint after seeding interaction data
- **Automated:** Unit test `recommendationService.go` — scoring, cold start fallback

---

## 4. Tool Selection & Configuration

### 4.1 Postman
- **Purpose:** Manual API testing and smoke testing all 28 endpoints
- **Collection:** `tests/ecommerce.postman_collection.json`
- **Environment Variables:** `baseUrl`, `accessToken`, `refreshToken`, `productId`, `categoryId`
- **Usage:** Import collection → run Login → copy token → test endpoints

### 4.2 Go Testing + Testify (Backend Unit Tests)
- **Purpose:** Automated unit tests for service layer
- **Framework:** Go built-in `testing` package + `github.com/stretchr/testify`
- **Location:** `internal/service/*_test.go`
- **Run:** `go test ./...`

### 4.3 Playwright (Frontend E2E)
- **Purpose:** Automated end-to-end testing of user flows
- **Config:** `playwright.config.ts` (baseURL: http://localhost:5173)
- **Tests location:** `e2e/` directory
- **Run:** `npm run test:e2e`
- **Browser:** Chromium (Desktop)

### 4.4 GitHub Actions (CI/CD)
- **Backend pipeline:** `.github/workflows/ci.yml` — triggers on push/PR to main
  - Steps: `go mod download` → `go build ./...` → `go vet ./...`
- **Frontend pipeline:** `.github/workflows/ci.yml` — triggers on push/PR to main
  - Steps: `npm ci` → `npm run lint` → `npm run build`

---

## 5. Test Environment

### 5.1 Local Development Environment
| Component | Details |
|-----------|---------|
| Backend | `go run cmd/web/main.go` → http://localhost:8080 |
| Frontend | `npm run dev` → http://localhost:5173 |
| Database | MongoDB via Docker: `docker-compose up -d mongodb` |
| Swagger UI | http://localhost:8080/swagger/index.html |

### 5.2 Test Data
- Seed database before testing: `go run scripts/seed/main.go`
- Test credentials: `admin@example.com` / `password123`

---

## 6. Planned Metrics

| Metric | Baseline (Now) | Target (End of Course) |
|--------|---------------|----------------------|
| API endpoint coverage | 0% | 100% (all 28 endpoints) |
| Service unit test coverage | 0% | 80% (HIGH-risk services) |
| E2E scenario coverage | 0% | 5 critical user flows |
| Defects found (HIGH severity) | 0 | tracked per sprint |
| CI pipeline pass rate | N/A | ≥ 95% |

### 6.1 Defect Severity Classification

| Severity | Definition | Example |
|----------|-----------|---------|
| **Critical** | System unusable | Login completely broken |
| **High** | Major feature broken | Purchase doesn't decrement stock |
| **Medium** | Feature works with workaround | Filter returns wrong results |
| **Low** | Minor visual/UX issue | Button misaligned |

---

## 7. Roles & Responsibilities

| Role | Responsibility |
|------|---------------|
| QA Engineer | Write and maintain tests, report defects |
| Developer | Fix defects, maintain CI pipeline |

---

## 8. Entry & Exit Criteria

### Entry Criteria (start testing)
- Application running locally (backend + frontend + MongoDB)
- Database seeded with test data
- Postman collection imported

### Exit Criteria (assignment complete)
- All 28 endpoints tested manually via Postman
- CI/CD pipelines green in both repositories
- All 4 deliverable documents created
- Screenshots captured for research paper
