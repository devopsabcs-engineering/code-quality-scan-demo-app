from flask import Blueprint, request, jsonify
from src.services.inventory_service import (
    process_inventory_update,
    get_inventory_summary,
    calculate_reorder_recommendations,
)
from src.utils.validators import validate_inventory_input, validate_sku_format
from src.utils.formatters import format_currency, format_quantity, format_date_string


inventory_bp = Blueprint("inventory", __name__)


@inventory_bp.route("/", methods=["GET"])
def list_inventory():
    """List all inventory items with optional filtering."""
    warehouse = request.args.get("warehouse")
    zone = request.args.get("zone")
    category = request.args.get("category")
    low_stock = request.args.get("low_stock")
    sort_by = request.args.get("sort_by", "name")
    order = request.args.get("order", "asc")

    # COMPLEXITY: deeply nested filtering logic (intentional violation)
    from src.models import INVENTORY_DATA
    results = []
    for item in INVENTORY_DATA:
        include = True
        if warehouse:
            if item.warehouse != warehouse:
                include = False
            else:
                if zone:
                    if item.zone != zone:
                        include = False
                    else:
                        if category:
                            if item.category != category:
                                include = False
                            else:
                                if low_stock and low_stock.lower() == "true":
                                    if item.quantity >= item.reorder_threshold:
                                        include = False
                                else:
                                    pass
                        else:
                            if low_stock and low_stock.lower() == "true":
                                if item.quantity >= item.reorder_threshold:
                                    include = False
                else:
                    if category:
                        if item.category != category:
                            include = False
                        else:
                            if low_stock and low_stock.lower() == "true":
                                if item.quantity >= item.reorder_threshold:
                                    include = False
                    else:
                        if low_stock and low_stock.lower() == "true":
                            if item.quantity >= item.reorder_threshold:
                                include = False
        else:
            if zone:
                if item.zone != zone:
                    include = False
                else:
                    if category:
                        if item.category != category:
                            include = False
                        else:
                            if low_stock and low_stock.lower() == "true":
                                if item.quantity >= item.reorder_threshold:
                                    include = False
                    else:
                        if low_stock and low_stock.lower() == "true":
                            if item.quantity >= item.reorder_threshold:
                                include = False
            else:
                if category:
                    if item.category != category:
                        include = False
                    else:
                        if low_stock and low_stock.lower() == "true":
                            if item.quantity >= item.reorder_threshold:
                                include = False
                else:
                    if low_stock and low_stock.lower() == "true":
                        if item.quantity >= item.reorder_threshold:
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

    # COMPLEXITY: manual sorting with nested conditions
    if sort_by == "name":
        if order == "asc":
            results.sort(key=lambda x: x["name"])
        else:
            results.sort(key=lambda x: x["name"], reverse=True)
    elif sort_by == "quantity":
        if order == "asc":
            results.sort(key=lambda x: x["quantity"])
        else:
            results.sort(key=lambda x: x["quantity"], reverse=True)
    elif sort_by == "price":
        if order == "asc":
            results.sort(key=lambda x: float(x["unit_price"].replace("$", "").replace(",", "")))
        else:
            results.sort(key=lambda x: float(x["unit_price"].replace("$", "").replace(",", "")), reverse=True)
    else:
        results.sort(key=lambda x: x["name"])

    return jsonify({"items": results, "count": len(results)})


@inventory_bp.route("/<int:item_id>", methods=["GET"])
def get_item(item_id):
    """Get a single inventory item by ID."""
    from src.models import INVENTORY_DATA
    for item in INVENTORY_DATA:
        if item.id == item_id:
            return jsonify({
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
    return jsonify({"error": "Item not found"}), 404


@inventory_bp.route("/<int:item_id>/update", methods=["POST"])
def update_item(item_id):
    """Update inventory item stock levels."""
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "No data provided"}), 400

        validation_result = validate_inventory_input(data)
        if not validation_result["valid"]:
            return jsonify({"error": validation_result["message"]}), 400

        result = process_inventory_update(item_id, data)
        return jsonify(result)
    except Exception:  # LINT: broad exception (intentional)
        return jsonify({"error": "Internal server error"}), 500


@inventory_bp.route("/summary", methods=["GET"])
def inventory_summary():
    """Get inventory summary by warehouse."""
    summary = get_inventory_summary()
    return jsonify(summary)


@inventory_bp.route("/reorder", methods=["GET"])
def reorder_recommendations():
    """Get items that need reordering."""
    recommendations = calculate_reorder_recommendations()
    return jsonify({"recommendations": recommendations})
