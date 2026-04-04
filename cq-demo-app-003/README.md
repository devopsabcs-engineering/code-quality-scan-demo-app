# cq-demo-app-003 — Employee & Payroll Management API (C# / ASP.NET Core)

Code Quality demo application with **intentional code quality violations** for scanner training and workshop demonstrations.

## Intentional Violations

| Category | Count | Files |
|----------|-------|-------|
| High Cyclomatic Complexity (CCN > 10) | 3 | `EmployeeService.cs`, `EmployeesController.cs`, `DataValidator.cs` |
| Code Duplication | 4 | `PayrollService.cs` ↔ `EmployeeService.cs`, `PayrollController.cs` ↔ `EmployeesController.cs` |
| .NET Analyzer Warnings | 10+ | `DataValidator.cs`, `ReportGenerator.cs` |
| Low Test Coverage (< 50%) | 1 | Only `EmployeesControllerTests.cs` — services and utils untested |

**Total expected findings: 18+**

## API Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/api/employees` | List all employees |
| GET | `/api/employees/{id}` | Get employee by ID |
| POST | `/api/employees` | Create employee |
| POST | `/api/employees/{id}/action` | Process employee action (hire, fire, promote, etc.) |
| GET | `/api/payroll` | List payroll records |
| POST | `/api/payroll/calculate` | Calculate payroll for employee |

## Run Locally

Build and run with Docker (works in GitHub Codespaces):

```bash
docker build -t cq-demo-app-003 .
docker run -p 8080:8080 cq-demo-app-003
```

Open http://localhost:8080 in your browser.

### Run with .NET CLI

```bash
cd src
dotnet run
```

Open http://localhost:8080 in your browser.

## Project Structure

```
cq-demo-app-003/
├── src/
│   ├── Controllers/        # API controllers
│   ├── Services/            # Business logic
│   ├── Models/              # Data models
│   └── Utils/               # Utility classes
├── tests/                   # Unit tests (intentionally incomplete)
├── infra/                   # Azure Bicep templates
├── Dockerfile               # Multi-stage container build
└── CqDemoApp003.sln         # Solution file
```

## Deployment

Deployed as a Docker container to Azure Web App for Containers via GitHub Actions. See `.github/workflows/deploy.yml`.
