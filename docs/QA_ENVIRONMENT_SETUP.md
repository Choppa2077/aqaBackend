# QA Environment Setup Report

**Project:** E-Commerce Platform
**Date:** 2026-03-21
**Author:** QA Team

---

## 1. Overview

This document describes the QA environment configured for testing the E-Commerce platform. The environment covers the full stack: Go REST API backend, React/TypeScript frontend, MongoDB database, API testing tools, E2E testing framework, and CI/CD pipelines.

---

## 2. System Requirements & Installed Tools

| Tool | Version | Purpose |
|------|---------|---------|
| Docker Desktop | 28.5.1 | Run MongoDB database in container |
| Go | 1.25.4 | Build and run backend API |
| Node.js | 24.9.0 | Build and run frontend |
| npm | bundled with Node | Frontend package management |
| Git | system | Version control |
| Postman | latest | Manual API testing |
| Playwright | 1.58.2 | Automated E2E frontend testing |

---

## 3. Repositories

| Repo | URL | Branch |
|------|-----|--------|
| Backend | https://github.com/Choppa2077/ecommerce-backend | `main` |
| Frontend | https://github.com/Choppa2077/ecommerce-frontend | `main` |

### Repository Structure

**Backend (`ecommerce-backend`):**
```
ecommerce-backend/
├── .github/workflows/ci.yml     ← CI/CD pipeline
├── cmd/web/main.go              ← Entry point
├── config/config.yaml           ← App configuration
├── docs/                        ← QA documentation
│   ├── RISK_ASSESSMENT.md
│   ├── TEST_STRATEGY.md
│   ├── QA_ENVIRONMENT_SETUP.md  ← this file
│   └── BASELINE_METRICS.md
├── internal/                    ← App source code
├── scripts/seed/main.go         ← DB seeder
├── tests/
│   └── ecommerce.postman_collection.json
└── docker-compose.yml
```

**Frontend (`ecommerce-frontend`):**
```
ecommerce-frontend/
├── .github/workflows/ci.yml     ← CI/CD pipeline
├── e2e/
│   └── auth.spec.ts             ← Playwright E2E tests
├── playwright.config.ts         ← Playwright configuration
├── src/                         ← App source code
└── package.json
```

---

## 4. Local Environment Setup

### Step 1 — Clone repositories

```bash
git clone https://github.com/Choppa2077/ecommerce-backend.git
git clone https://github.com/Choppa2077/ecommerce-frontend.git
```

### Step 2 — Start MongoDB

```bash
cd ecommerce-backend
docker-compose up -d mongodb
```

Verify MongoDB is running:
```bash
docker ps | grep mongodb
# Expected: ecommerce_mongodb   Up
```

### Step 3 — Seed the database

```bash
go run scripts/seed/main.go
```

This creates test users, categories, and products. Default credentials:
- `admin@example.com` / `password123`
- `user1@example.com` / `password123`

### Step 4 — Start the backend

```bash
go run cmd/web/main.go
```

Backend available at: `http://localhost:8080`
Swagger UI: `http://localhost:8080/swagger/index.html`

### Step 5 — Start the frontend

```bash
cd ecommerce-frontend
npm install
npm run dev
```

Frontend available at: `http://localhost:5173`

---

## 5. Postman Collection

### Location
`tests/ecommerce.postman_collection.json` in the backend repository.

### Import Steps
1. Open Postman
2. Click **Import** → select `ecommerce.postman_collection.json`
3. The collection **"E-Commerce API"** will appear with 4 folders and 28 requests

### Collection Variables

| Variable | Default Value | Description |
|----------|--------------|-------------|
| `baseUrl` | `http://localhost:8080/api/v1` | API base URL |
| `accessToken` | _(empty)_ | Set automatically after Login |
| `refreshToken` | _(empty)_ | Set automatically after Login |
| `productId` | `1` | Product ID for single-product requests |
| `categoryId` | `1` | Category ID for single-category requests |

### Usage

1. Run **Login** request first — it auto-saves `accessToken` and `refreshToken` to collection variables
2. All subsequent requests use `Bearer {{accessToken}}` automatically
3. Run requests in order: Auth → Products → Categories → Profiles

### Endpoint Coverage

| Folder | Requests |
|--------|----------|
| Auth | 3 |
| Products | 13 |
| Categories | 5 |
| Profiles | 9 |
| **Total** | **28** |

---

## 6. CI/CD Pipelines

### Backend Pipeline (GitHub Actions)

**File:** `.github/workflows/ci.yml` in `ecommerce-backend`

**Triggers:** Every push and pull request to `main`

**Steps:**
1. `actions/checkout@v4` — checkout code
2. `actions/setup-go@v5` with Go 1.24 — install Go
3. `go mod download` — download dependencies
4. `go build ./...` — compile entire project
5. `go vet ./...` — static analysis

**Purpose:** Ensures the backend compiles and passes static analysis on every commit.

### Frontend Pipeline (GitHub Actions)

**File:** `.github/workflows/ci.yml` in `ecommerce-frontend`

**Triggers:** Every push and pull request to `main`

**Steps:**
1. `actions/checkout@v4` — checkout code
2. `actions/setup-node@v4` with Node 20 — install Node.js
3. `npm ci` — install dependencies from lockfile
4. `npm run lint` — ESLint static analysis
5. `npm run build` — TypeScript compile + Vite production build

**Purpose:** Ensures the frontend builds cleanly and passes linting on every commit.

---

## 7. Playwright E2E Setup

### Installation

```bash
cd ecommerce-frontend
npm install -D @playwright/test   # already done
npx playwright install chromium   # install browser
```

### Configuration

**File:** `playwright.config.ts`

```typescript
{
  testDir: './e2e',
  use: {
    baseURL: 'http://localhost:5173',
  },
  projects: [{ name: 'chromium' }],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
  }
}
```

### Running E2E Tests

```bash
# Run all E2E tests
npm run test:e2e

# Run with UI (headed mode)
npx playwright test --headed

# View HTML report after run
npx playwright show-report
```

### Current Test Coverage

| File | Tests | Area |
|------|-------|------|
| `e2e/auth.spec.ts` | 4 | Login page, Register page, Auth redirect, Invalid login |

---

## 8. Verification Checklist

| Item | Status |
|------|--------|
| MongoDB running via Docker | ✅ |
| Backend starts on port 8080 | ✅ |
| Swagger UI accessible | ✅ |
| Frontend starts on port 5173 | ✅ |
| Postman collection imports (28 requests) | ✅ |
| Backend GitHub Actions pipeline | ✅ |
| Frontend GitHub Actions pipeline | ✅ |
| Playwright installed | ✅ |
| E2E test scaffold created | ✅ |

---

## 9. Screenshots

> See `docs/screenshots/` for:
> - Swagger UI with all endpoints
> - Postman collection with 28 requests
> - Backend CI/CD pipeline (green)
> - Frontend CI/CD pipeline (green)
> - Repository structure (backend and frontend)
