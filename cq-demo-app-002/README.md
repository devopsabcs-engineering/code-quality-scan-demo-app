# cq-demo-app-002 — Inventory Management API (Python / Flask)

A Python/Flask REST API for inventory management. This demo app contains **intentional code quality violations** for use with the Code Quality Scanner in the Agentic Accelerator Framework.

## Intentional Violations

| Category | File(s) | Description |
|----------|---------|-------------|
| **High Complexity** | `src/services/inventory_service.py` | `process_inventory_update()` has CCN > 15 with deeply nested if/elif/else for stock levels, warehouse zones, reorder thresholds, and supplier routing |
| **High Complexity** | `src/routes/inventory.py` | `list_inventory()` has deeply nested filter logic with 7+ nesting levels |
| **Code Duplication** | `src/services/report_service.py` ↔ `src/services/inventory_service.py` | Nearly identical aggregation logic with minor field name changes |
| **Code Duplication** | `src/routes/reports.py` ↔ `src/routes/inventory.py` | Repeated filtering and sorting patterns |
| **Lint — Unused Imports** | `src/utils/validators.py`, `src/services/*.py`, `src/utils/formatters.py` | `os`, `sys`, `hashlib`, `math`, `json`, `re` imported but never used |
| **Lint — Broad Exceptions** | `src/utils/validators.py` | `except Exception:` and bare `except:` used throughout |
| **Lint — Missing Docstrings** | `src/utils/formatters.py` | All functions lack docstrings |
| **Lint — Naming** | `src/utils/formatters.py` | `Format_Alert_Badge()`, `FormatReorderEstimate()` violate snake_case convention |
| **Lint — Mutable Default** | `src/utils/validators.py` | `validate_inventory_input(data, allowed_fields=[...])` uses mutable default argument |
| **Lint — format() vs f-string** | `src/utils/formatters.py` | Uses `"{}".format()` instead of f-strings throughout |
| **Low Coverage** | `tests/test_app.py` | Only 3 tests covering health, index, and list. Services and utils are completely untested (~30% coverage) |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | App info |
| GET | `/health` | Health check |
| GET | `/api/inventory/` | List inventory (query params: `warehouse`, `zone`, `category`, `low_stock`, `sort_by`, `order`) |
| GET | `/api/inventory/<id>` | Get item by ID |
| POST | `/api/inventory/<id>/update` | Update stock (body: `action`, `quantity`, `reason`, `target_warehouse`, `priority`) |
| GET | `/api/inventory/summary` | Inventory summary by warehouse |
| GET | `/api/inventory/reorder` | Reorder recommendations |
| GET | `/api/reports/inventory` | Inventory status report |
| GET | `/api/reports/warehouse-utilization` | Warehouse utilization report |
| GET | `/api/reports/supplier-performance` | Supplier performance report |
| GET | `/api/reports/stock-alerts` | Stock alert report |

## Run Locally

Build and run with Docker (works in GitHub Codespaces):

```bash
docker build -t cq-demo-app-002 .
docker run -p 5000:5000 cq-demo-app-002
```

Then open [http://localhost:5000](http://localhost:5000) to see the app info, or try:

```bash
curl http://localhost:5000/health
curl http://localhost:5000/api/inventory/
curl http://localhost:5000/api/inventory/1
curl http://localhost:5000/api/inventory/summary
curl http://localhost:5000/api/inventory/reorder
curl http://localhost:5000/api/reports/stock-alerts
```

### Run without Docker

```bash
pip install -r requirements.txt
python -m flask --app src/app run --host=0.0.0.0 --port=5000
```

### Run Tests

```bash
pip install -r requirements.txt
pytest tests/ -v --cov=src --cov-report=term-missing
```

## Scanning

This app is designed to produce findings with:

- **Ruff** — unused imports, missing docstrings, broad exceptions, naming conventions, mutable defaults, format vs f-string
- **Lizard** — high cyclomatic complexity (CCN > 15) in `process_inventory_update()`
- **jscpd** — code duplication between `inventory_service.py` / `report_service.py` and `inventory.py` / `reports.py`
- **pytest-cov** — low test coverage (< 50%)

## Tech Stack

- Python 3.12
- Flask 3.x
- pytest + pytest-cov
- Docker (multi-stage build)
- Azure Web App for Containers (via Bicep + ACR)
