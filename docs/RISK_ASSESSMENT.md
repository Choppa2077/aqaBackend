# Risk Assessment Document

**Project:** E-Commerce Platform (Backend API + Frontend)
**System Type:** Web Application / REST API
**Assessment Date:** 2026-03-21
**Author:** QA Team

---

## 1. System Overview

The system is a full-stack e-commerce platform consisting of:

- **Backend:** Go REST API (Gin framework) with MongoDB, JWT authentication, collaborative filtering recommendations
- **Frontend:** React 19 + TypeScript (Vite), Zustand state management, TanStack Query, Playwright-ready

**Total API Endpoints:** 28
**Frontend Pages:** 9
**Database:** MongoDB (no SQL transactions — risk factor)

---

## 2. Risk Assessment Methodology

Each module is evaluated on two dimensions:

| Dimension | Scale | Description |
|-----------|-------|-------------|
| **Probability** | 1–5 | How likely is a failure in this module? |
| **Impact** | 1–5 | How severe are the consequences if it fails? |
| **Risk Score** | 1–25 | Probability × Impact |

**Priority thresholds:**
- **HIGH:** Score ≥ 10
- **MEDIUM:** Score 5–9
- **LOW:** Score < 5

---

## 3. Module Risk Table

| # | Module | Description | Probability | Impact | Score | Priority |
|---|--------|-------------|-------------|--------|-------|----------|
| 1 | **Authentication (JWT)** | Register, login, token generation/validation | 3 | 5 | **15** | 🔴 HIGH |
| 2 | **Token Refresh Interceptor** | Frontend auto-refresh on 401, redirect on failure | 3 | 4 | **12** | 🔴 HIGH |
| 3 | **Product CRUD** | Create/update/delete products, category linking | 3 | 4 | **12** | 🔴 HIGH |
| 4 | **Recommendation Engine** | Collaborative filtering, cosine similarity scoring | 4 | 3 | **12** | 🔴 HIGH |
| 5 | **Purchase Flow + Stock** | Purchase validation, stock decrement, quantity checks | 2 | 5 | **10** | 🔴 HIGH |
| 6 | **Profile Management** | Update profile, change password, soft delete account | 2 | 3 | **6** | 🟡 MEDIUM |
| 7 | **Category Management** | CRUD for product categories, parent-child hierarchy | 2 | 3 | **6** | 🟡 MEDIUM |
| 8 | **Product Filtering & Search** | Pagination, price range, category filter, text search | 3 | 2 | **6** | 🟡 MEDIUM |
| 9 | **Like/View Interactions** | Recording views and likes, interaction history | 2 | 2 | **4** | 🟢 LOW |

---

## 4. Detailed Risk Analysis

### 🔴 Module 1: Authentication (JWT)
**Endpoints:** `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`
**Critical Files:** `internal/service/authService.go`, `internal/delivery/middleware/auth.go`

**Risks:**
- Weak/missing input validation on register (no email format check enforced at service level)
- JWT secret is configurable — misconfiguration could expose all tokens
- No rate limiting on login endpoint → brute force vulnerability
- Refresh token has 7-day expiry — if compromised, attacker has long window

**Assumptions:** bcrypt cost=10 is sufficient; JWT library (golang-jwt) is trusted

---

### 🔴 Module 2: Token Refresh Interceptor (Frontend)
**Files:** `src/shared/api/apiInstance.ts`

**Risks:**
- Race condition: multiple concurrent requests all receive 401 simultaneously, all attempt refresh
- If refresh token expires during session, user gets silently redirected without clear message
- localStorage token storage is vulnerable to XSS

**Assumptions:** Single-tab usage assumed; XSS protection via framework-level escaping

---

### 🔴 Module 3: Product CRUD
**Endpoints:** `POST/GET/PUT/DELETE /api/v1/products`, `GET /api/v1/products/:id/statistics`
**Critical Files:** `internal/service/productService.go`, `internal/repository/productRepository.go`

**Risks:**
- Admin role check is **not implemented** (`// TODO: Check if user has admin role` in handlers) — any authenticated user can create/delete products
- No validation that `category_id` exists before creating a product
- Stock can go negative if concurrent purchases happen (no MongoDB transaction)

**Assumptions:** Admin role enforcement is a known gap, to be tested

---

### 🔴 Module 4: Recommendation Engine
**Endpoints:** `GET /api/v1/profiles/me/recommendations`, `GET /api/v1/profiles/me/similar`
**Critical Files:** `internal/service/recommendationService.go`

**Risks:**
- Algorithm fetches ALL user interactions from DB without pagination — performance risk at scale
- Cold start problem: new users with no interactions receive empty or random recommendations
- Cosine similarity calculation done in-memory — no caching

**Assumptions:** Dataset is small (dev/test environment); production scale not yet considered

---

### 🔴 Module 5: Purchase Flow + Stock Validation
**Endpoint:** `POST /api/v1/products/:id/purchase`
**Critical Files:** `internal/service/interactionService.go`, `internal/repository/productRepository.go`

**Risks:**
- No MongoDB transaction: stock check and decrement are two separate operations → race condition on concurrent purchases
- No rollback mechanism if purchase record is saved but stock update fails
- Quantity validation exists (> 0) but no max quantity check

**Assumptions:** Low concurrency in test environment masks race condition

---

### 🟡 Module 6: Profile Management
**Endpoints:** `GET/PUT /api/v1/profiles/me`, `PUT /api/v1/profiles/me/password`, `DELETE /api/v1/profiles/me/account`

**Risks:**
- Password change does not invalidate existing JWT tokens
- Soft delete doesn't revoke active tokens — deleted user can still use the API until token expires

---

### 🟡 Module 7: Category Management
**Endpoints:** `GET/POST/PUT/DELETE /api/v1/categories`

**Risks:**
- Admin role check not implemented (same `// TODO` issue as products)
- Deleting a category that has products doesn't validate orphaned products

---

### 🟡 Module 8: Product Filtering & Search
**Endpoint:** `GET /api/v1/products`

**Risks:**
- No input sanitization on `search` query parameter (potential MongoDB injection via `$where`)
- Sorting by unsupported fields silently falls back to default (no error returned)

---

### 🟢 Module 9: Like/View Interactions
**Endpoints:** `POST/DELETE /products/:id/like`, `POST /products/:id/view`

**Risks:**
- Duplicate views can be recorded (no deduplication within a session)
- Low business impact if this data is slightly inaccurate

---

## 5. Risk Priority Summary

| Priority | Count | Modules |
|----------|-------|---------|
| 🔴 HIGH | 5 | Auth, Token Refresh, Product CRUD, Recommendations, Purchase Flow |
| 🟡 MEDIUM | 3 | Profile, Categories, Product Filtering |
| 🟢 LOW | 1 | Like/View Interactions |

---

## 6. Key Findings

1. **Critical security gap:** Admin role is not enforced on write endpoints (products, categories) — any logged-in user can modify catalog data
2. **Race condition risk:** Purchase flow lacks atomic stock operations (no MongoDB transactions)
3. **Frontend token storage:** localStorage usage is an XSS risk in production
4. **Zero existing tests:** No automated tests exist in either backend or frontend — 0% coverage baseline

---

## 7. Assumptions & Constraints

- System is evaluated in a **development/local environment** (Docker + MongoDB)
- No production load or real user data exists yet
- Admin role enforcement is documented as `TODO` in source code and is a known gap
- MongoDB replica set (required for transactions) is not configured in current Docker setup
