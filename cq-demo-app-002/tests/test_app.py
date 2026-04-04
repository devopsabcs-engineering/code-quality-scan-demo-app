"""Minimal tests for cq-demo-app-002 — intentionally low coverage.

Only the health endpoint and one inventory route are tested.
Services (inventory_service, report_service) and utilities
(validators, formatters) are completely untested.
This is an INTENTIONAL code quality violation for demo purposes.
"""
import pytest
from src.app import create_app


@pytest.fixture
def client():
    """Create test client."""
    app = create_app()
    app.config["TESTING"] = True
    with app.test_client() as client:
        yield client


def test_health_endpoint(client):
    """Test the /health endpoint returns healthy status."""
    response = client.get("/health")
    assert response.status_code == 200
    data = response.get_json()
    assert data["status"] == "healthy"


def test_index_endpoint(client):
    """Test the / endpoint returns app info."""
    response = client.get("/")
    assert response.status_code == 200
    data = response.get_json()
    assert data["app"] == "cq-demo-app-002"
    assert data["status"] == "ok"


def test_list_inventory(client):
    """Test the /api/inventory/ endpoint returns items."""
    response = client.get("/api/inventory/")
    assert response.status_code == 200
    data = response.get_json()
    assert "items" in data
    assert "count" in data
    assert data["count"] > 0


# NOTE: No tests for:
# - POST /api/inventory/<id>/update (process_inventory_update)
# - GET /api/inventory/summary (get_inventory_summary)
# - GET /api/inventory/reorder (calculate_reorder_recommendations)
# - GET /api/reports/inventory (generate_inventory_report)
# - GET /api/reports/warehouse-utilization (generate_warehouse_utilization_report)
# - GET /api/reports/supplier-performance (generate_supplier_performance_report)
# - GET /api/reports/stock-alerts
# - src/utils/validators.py (validate_inventory_input, validate_sku_format, etc.)
# - src/utils/formatters.py (format_currency, format_quantity, etc.)
#
# This leaves approximately 70-80% of the codebase untested.
