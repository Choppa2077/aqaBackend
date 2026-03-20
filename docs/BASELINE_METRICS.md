# Baseline Metrics

**Project:** E-Commerce Platform
**Date:** 2026-03-21
**Purpose:** Establish starting metrics for the research paper. These numbers represent the state **before** any automated testing is implemented.

---

## 1. System Size Metrics

| Metric | Value |
|--------|-------|
| Total API endpoints | **28** |
| Public endpoints (no auth) | **3** |
| Protected endpoints (auth required) | **25** |
| Frontend pages | **9** |
| Backend service modules | **5** |
| Backend repository modules | **4** |
| Lines of code (backend) | ~2,500 |
| Lines of code (frontend) | ~1,800 |

---

## 2. Endpoint Inventory

| Module | Endpoints | Risk |
|--------|-----------|------|
| Auth | 3 | 🔴 HIGH |
| Products (CRUD) | 5 | 🔴 HIGH |
| Product Interactions | 6 | 🟡 MEDIUM / 🟢 LOW |
| Categories | 5 | 🟡 MEDIUM |
| Profiles | 9 | 🟡 MEDIUM |
| **Total** | **28** | — |

---

## 3. Test Coverage Baseline

| Area | Existing Tests | Coverage |
|------|---------------|----------|
| Backend unit tests (`*_test.go`) | **0** | **0%** |
| Frontend unit tests (`*.test.tsx`) | **0** | **0%** |
| E2E tests (Playwright) | **0** | **0%** |
| API tests (Postman automated) | **0** | **0%** |

> **Note:** Zero test coverage is the baseline. All future assignments will show improvement from this starting point.

---

## 4. Risk Module Count

| Priority | Count | Modules |
|----------|-------|---------|
| 🔴 HIGH | **5** | Auth, Token Refresh, Product CRUD, Purchase Flow, Recommendations |
| 🟡 MEDIUM | **3** | Profile, Categories, Product Filtering |
| 🟢 LOW | **1** | Like/View Interactions |
| **Total** | **9** | — |

---

## 5. Known Defects at Baseline

| ID | Module | Severity | Description |
|----|--------|----------|-------------|
| BUG-001 | Product CRUD | **Critical** | Admin role check not implemented — any authenticated user can create/update/delete products |
| BUG-002 | Category CRUD | **Critical** | Same admin role gap as products |
| BUG-003 | Purchase Flow | **High** | No atomic transaction — race condition possible on concurrent purchases |
| BUG-004 | Profile | **Medium** | Password change does not invalidate existing JWT tokens |
| BUG-005 | Auth | **Medium** | No rate limiting on login endpoint |

---

## 6. Estimated Testing Effort

| Area | Estimated Hours |
|------|----------------|
| Manual API testing (Postman, all 28 endpoints) | 4h |
| Backend unit tests (HIGH-risk services) | 16h |
| Frontend E2E tests (5 critical flows) | 8h |
| CI/CD pipeline setup & verification | 2h |
| Documentation | 4h |
| **Total estimate** | **~34h** |

---

## 7. QA Environment Setup Status

| Tool | Status | Notes |
|------|--------|-------|
| Docker + MongoDB | ✅ Running | `docker-compose up -d` |
| Go backend | ✅ Running | http://localhost:8080 |
| Swagger UI | ✅ Available | http://localhost:8080/swagger/index.html |
| Postman collection | ✅ Created | `tests/ecommerce.postman_collection.json` (28 endpoints) |
| GitHub Actions — Backend | ✅ Configured | `.github/workflows/ci.yml` (build + vet) |
| GitHub Actions — Frontend | ✅ Configured | `.github/workflows/ci.yml` (lint + build) |
| Playwright | ✅ Installed | `playwright.config.ts`, `e2e/auth.spec.ts` |

---

## 8. Screenshots

> Screenshots are saved in `docs/screenshots/` directory.
>
> Required screenshots:
> - [ ] Swagger UI showing all endpoints
> - [ ] Postman collection imported with 28 requests
> - [ ] GitHub Actions — backend pipeline (green)
> - [ ] GitHub Actions — frontend pipeline (green)
> - [ ] Repository structure (backend)
> - [ ] Repository structure (frontend)

---

## 9. Research Paper Connection

These baseline metrics feed into the **Introduction and Methodology** chapters:

- **System description:** E-commerce REST API (Go + MongoDB) + React/TS frontend
- **Initial state:** 0 tests, 0% coverage, 5 known defects
- **Risk methodology:** Probability × Impact matrix (see RISK_ASSESSMENT.md)
- **Future comparison:** Later assignments will show test count, coverage %, and defect detection rate improvements against this baseline
