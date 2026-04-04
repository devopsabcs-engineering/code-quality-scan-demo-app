# CQ Demo App 001 — E-Commerce Dashboard (TypeScript / Next.js)

A simple e-commerce dashboard built with Next.js 15 and TypeScript. This app contains **intentional code quality violations** for scanner detection during workshop exercises.

## Intentional Violations

This demo app is designed to trigger findings across all four code quality scanner categories:

| Category | Files | Violation Type | Expected Findings |
|----------|-------|---------------|-------------------|
| **Cyclomatic Complexity** | `src/lib/calculations.ts` | Functions with CCN > 15 (`calculateOrderTotal`, `processPayment`) | 3+ |
| **Code Duplication** | `src/lib/validators.ts`, `src/lib/formatters.ts` | Nearly identical email, phone, and date logic copied across files | 3+ |
| **ESLint / Lint** | `src/lib/helpers.ts`, `src/app/api/*/route.ts` | `prefer-const`, `no-unused-vars`, `no-explicit-any` | 20+ |
| **Test Coverage** | `__tests__/calculations.test.ts` | Only 2 tests covering 1 function; most code untested (< 30% coverage) | 8+ files below 80% |

**Total expected findings: 30+**

## Tech Stack

- **Runtime**: Node.js 20
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript 5.7
- **Testing**: Jest 29 with ts-jest
- **Linting**: ESLint 9 (flat config) with `@typescript-eslint`
- **Container**: Multi-stage Docker build with standalone output

## Project Structure

```
cq-demo-app-001/
├── src/
│   ├── app/                    # Next.js App Router pages and API routes
│   │   ├── layout.tsx          # Root layout
│   │   ├── page.tsx            # Dashboard home page
│   │   ├── globals.css         # Global styles
│   │   ├── api/
│   │   │   ├── health/route.ts # Health check endpoint
│   │   │   ├── products/route.ts
│   │   │   └── orders/route.ts
│   │   ├── products/page.tsx   # Products catalog page
│   │   └── orders/page.tsx     # Orders management page
│   └── lib/
│       ├── calculations.ts     # ← HIGH COMPLEXITY (CCN > 15)
│       ├── validators.ts       # ← DUPLICATION (copied from formatters.ts)
│       ├── formatters.ts       # ← DUPLICATION (copied from validators.ts)
│       └── helpers.ts          # ← LINT VIOLATIONS (unused vars, any, let)
├── __tests__/
│   └── calculations.test.ts    # ← LOW COVERAGE (only 2 tests)
├── infra/main.bicep            # Azure ACR + Web App for Containers
├── Dockerfile                  # Multi-stage build
├── package.json
├── tsconfig.json
├── jest.config.ts
└── eslint.config.mjs
```

## Run Locally

Build and run with Docker (works in GitHub Codespaces):

```bash
docker build -t cq-demo-app-001 .
docker run -p 3000:3000 cq-demo-app-001
```

Then open [http://localhost:3000](http://localhost:3000) to view the dashboard.

### Without Docker

```bash
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check — returns app status |
| `/api/products` | GET | Returns product catalog |
| `/api/orders` | GET | Returns order list with revenue summary |

## Running Tests

```bash
npm test                 # Run tests
npm run test:coverage    # Run tests with coverage report
```

Expected coverage: **< 30%** (intentional — only `calculateOrderTotal` is tested).

## Running Lint

```bash
npm run lint             # Next.js ESLint
npx eslint src/          # Direct ESLint
```

Expected: **20+ warnings** from `prefer-const`, `no-unused-vars`, and `no-explicit-any` rules.

## Azure Deployment

This app deploys as a Docker container to Azure Web App for Containers via the GitHub Actions workflow in `.github/workflows/deploy.yml`.

Infrastructure (ACR + App Service Plan + Web App) is provisioned via Bicep in `infra/main.bicep` using `uniqueString(resourceGroup().id)` for globally unique resource names.

## License

This demo app is part of the [Agentic Accelerator Framework](https://github.com/devopsabcs-engineering/agentic-accelerator-framework) and is provided for educational purposes.
