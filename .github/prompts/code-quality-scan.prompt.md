---
description: "Run a comprehensive code quality scan across all demo apps"
---

Scan all 5 demo apps (cq-demo-app-001 through cq-demo-app-005) for code quality issues:

1. Run ESLint/Ruff/golangci-lint for lint violations
2. Run jscpd for code duplication
3. Run Lizard for cyclomatic complexity
4. Run coverage tools for test coverage gaps
5. Aggregate results into SARIF v2.1.0
6. Upload to GitHub Security tab

Use matrix strategy to scan all apps in parallel.
Report total finding count by severity.
