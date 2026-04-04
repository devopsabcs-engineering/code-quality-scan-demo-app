---
name: code-quality-scan
description: "Code quality scanning methodology — 4-tool architecture, MegaLinter orchestration, SARIF output, coverage and complexity mapping, and CI/CD integration"
---

# Code Quality Scan Skill

Domain knowledge for code quality scanning across multiple languages. Agents load this skill to understand the scanner architecture, tool stack, output format, severity classification, and compliance thresholds.

## Scanner Architecture

The code quality scanner uses a **4-tool architecture** with a MegaLinter orchestrator for maximum language coverage:

| Tool | Role | Technology |
|---|---|---|
| **Per-language linters** | Static analysis for style, correctness, and best practices | ESLint, Ruff, golangci-lint, .NET Analyzers, Checkstyle/PMD |
| **jscpd** | Code duplication detection across files and languages | jscpd (native SARIF output) |
| **Lizard** | Cyclomatic complexity and function length analysis | Lizard CLI (converter: `lizard-to-sarif.py`) |
| **Per-language coverage** | Test coverage measurement and gap identification | jest, pytest-cov, Coverlet, JaCoCo, go test -cover (converter: `coverage-to-sarif.py`) |
| **MegaLinter** (orchestrator) | Multi-language aggregation and native SARIF output | MegaLinter v8+ |

### Scan Flow

1. Run MegaLinter with per-language linter configuration (`.mega-linter.yml`).
2. Run jscpd with `.jscpd.json` config for duplication detection.
3. Run Lizard for complexity analysis and pipe output through `lizard-to-sarif.py`.
4. Run per-language coverage tools and pipe output through `coverage-to-sarif.py`.
5. Aggregate all SARIF outputs into a unified result set.
6. Upload to GitHub Security tab (`codeql-action/upload-sarif@v4`) or ADO Advanced Security (`AdvancedSecurity-Publish@1`).

## SARIF Output Format

The scanner produces SARIF v2.1.0 compliant output.

### Required SARIF Fields

| Field | Value |
|---|---|
| `$schema` | `https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json` |
| `version` | `2.1.0` |
| `tool.driver.name` | `code-quality-scanner` |
| `automationDetails.id` | `code-quality/coverage/<project>` |
| `partialFingerprints` | Hash of `ruleId:file:function` for deduplication |
| `results[].ruleId` | Unique rule ID per finding type (e.g., `coverage-below-threshold`, `ccn-exceeded`) |
| `results[].level` | Mapped from severity (see severity mapping) |
| `help.markdown` | Rule description, threshold reference, remediation guidance |
| `properties.tags` | Includes `code-quality` plus category tags (`coverage`, `complexity`, `duplication`, `lint`) |

### SARIF Enrichment

- `help.markdown` — Embeds rule description, threshold reference, and remediation guidance (GitHub does not render `helpUri`; URLs must be embedded in markdown).
- `properties.precision` — `very-high` for direct measurements (coverage, CCN), `high` for lint findings.
- `properties.tags` — Includes `code-quality` tag plus category tags for GitHub filtering.
- `partialFingerprints` — Hash of `ruleId:file:function` for deduplication across runs.
- `automationDetails.id` — Category prefix `code-quality/coverage/` for multi-tool scenarios.

## Severity Mapping

| Condition | SARIF Level | Description |
|---|---|---|
| Coverage < 50%, CCN > 20, critical lint errors | `error` | Immediate action required — block merge |
| Coverage 50–79%, CCN 11–20, moderate lint warnings | `warning` | Should be addressed in current sprint |
| Coverage 80–89%, CCN 6–10, minor style issues | `note` | Track for future improvement |
| Coverage ≥ 90%, CCN ≤ 5 | Pass | No finding generated |

### CWE Mapping

| Finding Type | CWE ID | Description |
|---|---|---|
| High cyclomatic complexity | CWE-1121 | Excessive McCabe Cyclomatic Complexity |
| Code duplication | CWE-1041 | Use of Redundant Code |
| Missing error handling | CWE-754 | Improper Check for Unusual Conditions |

## Thresholds

All projects must maintain the following minimum quality levels, as defined in `instructions/code-quality.instructions.md`:

| Metric | Threshold | Enforcement |
|---|---|---|
| Line coverage | ≥ 80% | CI gate — block merge if below |
| Branch coverage | ≥ 80% | CI gate — block merge if below |
| Function coverage | ≥ 80% | CI gate — block merge if below |
| New code coverage | ≥ 90% | PR check — warn if below |
| Cyclomatic complexity | ≤ 10 per function | CI gate — block merge if exceeded |
| Nesting depth | ≤ 4 levels | CI gate — block merge if exceeded |
| Function length | ≤ 50 lines | Review — warn if exceeded |
| Duplication | < 10 similar consecutive lines | Review — flag for extraction |

## SARIF Converter Patterns

Two converters transform tool-native output into SARIF v2.1.0:

### `coverage-to-sarif.py`

```text
Usage: coverage-to-sarif.py --input <coverage-file> --format <format> --output <sarif-file> [--threshold 80]

Supported formats: cobertura, json-summary, lcov, jacoco, gocover
```

- Accepts Cobertura XML (pytest-cov, Coverlet), JSON summary (jest), lcov, JaCoCo XML, and Go cover profiles.
- Emits one SARIF result per file below the coverage threshold.
- Sets `automationDetails.id` to `code-quality/coverage/`.

### `lizard-to-sarif.py`

```text
Usage: lizard-to-sarif.py --input <lizard-csv> --output <sarif-file> [--ccn-threshold 10] [--length-threshold 50]
```

- Accepts Lizard CSV output (`lizard --csv`).
- Maps each function exceeding the CCN or length threshold to a SARIF result.
- Sets `automationDetails.id` to `code-quality/complexity/`.

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run MegaLinter
  uses: oxsecurity/megalinter@v8
  env:
    VALIDATE_ALL_CODEBASE: true
    SARIF_REPORTER: true

- name: Run complexity analysis
  run: |
    lizard --csv src/ > lizard-output.csv
    python scripts/lizard-to-sarif.py --input lizard-output.csv --output complexity.sarif

- name: Run coverage
  run: |
    npm test -- --coverage --coverageReporters=json-summary
    python scripts/coverage-to-sarif.py --input coverage/coverage-summary.json --format json-summary --output coverage.sarif

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v4
  with:
    sarif_file: ./results
    category: code-quality/coverage
```

### Azure DevOps

```yaml
- script: |
    npx mega-linter-runner --flavor dotnetweb
  displayName: 'Run MegaLinter'

- script: |
    lizard --csv src/ > lizard-output.csv
    python scripts/lizard-to-sarif.py --input lizard-output.csv --output $(Build.ArtifactStagingDirectory)/complexity.sarif
  displayName: 'Complexity Analysis'

- task: AdvancedSecurity-Publish@1
  inputs:
    SarifFileDirectory: '$(Build.ArtifactStagingDirectory)'
```

## Per-Language Tool Reference

| Language | Linter | Coverage | Config |
|---|---|---|---|
| TypeScript/JavaScript | ESLint | jest / vitest | `eslint.config.mjs`, `jest.config.ts` |
| Python | Ruff | pytest-cov | `pyproject.toml` |
| C# | .NET Analyzers | Coverlet | `*.csproj`, `.editorconfig` |
| Java | Checkstyle / PMD | JaCoCo | `pom.xml`, `checkstyle.xml` |
| Go | golangci-lint | go test -cover | `.golangci.yml`, `go.mod` |

## References

- [ESLint Documentation](https://eslint.org/docs/latest/)
- [Ruff Documentation](https://docs.astral.sh/ruff/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [MegaLinter Documentation](https://megalinter.io/)
- [jscpd Documentation](https://github.com/kucherenko/jscpd)
- [Lizard Documentation](https://github.com/terryyin/lizard)
- [SARIF v2.1.0 Specification](https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
