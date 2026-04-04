from flask import Blueprint, request, jsonify
from src.services.report_service import (
    generate_inventory_report,
    generate_warehouse_utilization_report,
    generate_supplier_performance_report,
)
from src.utils.formatters import format_currency, format_quantity, format_date_string


reports_bp = Blueprint("reports", __name__)


@reports_bp.route("/inventory", methods=["GET"])
def inventory_report():
    """Generate inventory status report."""
    warehouse = request.args.get("warehouse")
    category = request.args.get("category")

    # DUPLICATION: repeated query/filter logic from inventory.py (intentional)
    from src.models import INVENTORY_DATA
    results = []
    for item in INVENTORY_DATA:
        include = True
        if warehouse:
            if item.warehouse != warehouse:
                include = False
            else:
                if category:
                    if item.category != category:
                        include = False
        else:
            if category:
                if item.category != category:
                    include = False

        if include:
            results.append({
                "id": item.id,
                "name": item.name,
                "sku": item.sku,
                "quantity": item.quantity,
                "warehouse": item.warehouse,
                "zone": item.zone,
                "reorder_threshold": item.reorder_threshold,
                "unit_price": format_currency(item.unit_price),
                "category": item.category,
            })

    report = generate_inventory_report(results)
    return jsonify(report)


@reports_bp.route("/warehouse-utilization", methods=["GET"])
def warehouse_utilization():
    """Generate warehouse utilization report."""
    report = generate_warehouse_utilization_report()
    return jsonify(report)


@reports_bp.route("/supplier-performance", methods=["GET"])
def supplier_performance():
    """Generate supplier performance report."""
    report = generate_supplier_performance_report()
    return jsonify(report)


@reports_bp.route("/stock-alerts", methods=["GET"])
def stock_alerts():
    """Get stock alert report."""
    # DUPLICATION: same filtering and aggregation pattern as inventory route (intentional)
    from src.models import INVENTORY_DATA
    alerts = []
    for item in INVENTORY_DATA:
        if item.quantity <= item.reorder_threshold:
            severity = "critical"
            if item.quantity == 0:
                severity = "out_of_stock"
            elif item.quantity <= item.reorder_threshold * 0.25:
                severity = "critical"
            elif item.quantity <= item.reorder_threshold * 0.5:
                severity = "high"
            elif item.quantity <= item.reorder_threshold * 0.75:
                severity = "medium"
            else:
                severity = "low"

            alerts.append({
                "id": item.id,
                "name": item.name,
                "sku": item.sku,
                "quantity": item.quantity,
                "warehouse": item.warehouse,
                "zone": item.zone,
                "reorder_threshold": item.reorder_threshold,
                "unit_price": format_currency(item.unit_price),
                "category": item.category,
                "severity": severity,
                "deficit": item.reorder_threshold - item.quantity,
            })

    # DUPLICATION: manual sorting repeated from inventory route
    sort_by = request.args.get("sort_by", "severity")
    severity_order = {"out_of_stock": 0, "critical": 1, "high": 2, "medium": 3, "low": 4}
    if sort_by == "severity":
        alerts.sort(key=lambda x: severity_order.get(x["severity"], 99))
    elif sort_by == "name":
        alerts.sort(key=lambda x: x["name"])
    elif sort_by == "quantity":
        alerts.sort(key=lambda x: x["quantity"])
    elif sort_by == "deficit":
        alerts.sort(key=lambda x: x["deficit"], reverse=True)
    else:
        alerts.sort(key=lambda x: severity_order.get(x["severity"], 99))

    return jsonify({
        "alerts": alerts,
        "total_alerts": len(alerts),
        "critical_count": len([a for a in alerts if a["severity"] in ("critical", "out_of_stock")]),
    })
