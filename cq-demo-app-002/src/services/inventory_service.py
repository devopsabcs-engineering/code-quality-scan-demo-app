import os, sys, json, re, datetime  # LINT: unused imports os, sys, re (intentional)
from src.models import INVENTORY_DATA, WAREHOUSE_DATA, SUPPLIER_DATA


def process_inventory_update(item_id, data):
    """Process an inventory update with complex business rules.

    COMPLEXITY: CCN > 15 — deeply nested if/elif/else for stock levels,
    warehouse zones, reorder thresholds, and supplier routing.
    This is an INTENTIONAL code quality violation for demo purposes.
    """
    item = None
    for inv_item in INVENTORY_DATA:
        if inv_item.id == item_id:
            item = inv_item
            break

    if item is None:
        return {"success": False, "error": "Item not found"}

    action = data.get("action", "adjust")
    quantity_change = data.get("quantity", 0)
    reason = data.get("reason", "manual")
    target_warehouse = data.get("target_warehouse", None)
    priority = data.get("priority", "normal")

    result = {"success": True, "item_id": item_id, "actions_taken": []}

    # --- BEGIN COMPLEX NESTED LOGIC (intentional high CCN) ---
    if action == "adjust":
        new_quantity = item.quantity + quantity_change
        if new_quantity < 0:
            return {"success": False, "error": "Cannot reduce below zero"}
        else:
            item.quantity = new_quantity
            result["actions_taken"].append("quantity_adjusted")

            # Check reorder thresholds after adjustment
            if item.quantity <= item.reorder_threshold:
                if item.quantity == 0:
                    # Out of stock — emergency reorder
                    if priority == "emergency":
                        supplier = _find_supplier(item.supplier_id)
                        if supplier:
                            if supplier.reliability_score >= 0.9:
                                result["actions_taken"].append("emergency_reorder_placed")
                                result["reorder_quantity"] = item.reorder_threshold * 3
                                result["estimated_days"] = supplier.lead_time_days
                            elif supplier.reliability_score >= 0.8:
                                result["actions_taken"].append("emergency_reorder_placed_backup")
                                result["reorder_quantity"] = item.reorder_threshold * 2
                                result["estimated_days"] = supplier.lead_time_days + 3
                            else:
                                result["actions_taken"].append("emergency_reorder_manual_review")
                                result["warning"] = "Supplier reliability too low for auto-reorder"
                        else:
                            result["actions_taken"].append("emergency_no_supplier")
                    elif priority == "high":
                        supplier = _find_supplier(item.supplier_id)
                        if supplier:
                            if supplier.region == "domestic":
                                result["actions_taken"].append("high_priority_domestic_reorder")
                                result["reorder_quantity"] = item.reorder_threshold * 2
                            else:
                                result["actions_taken"].append("high_priority_international_reorder")
                                result["reorder_quantity"] = item.reorder_threshold * 3
                                result["warning"] = "International supplier — longer lead time"
                        else:
                            result["actions_taken"].append("high_priority_no_supplier")
                    else:
                        result["actions_taken"].append("standard_reorder_triggered")
                        result["reorder_quantity"] = item.reorder_threshold
                elif item.quantity <= item.reorder_threshold * 0.25:
                    # Critical low stock
                    if item.zone == "A":
                        result["actions_taken"].append("critical_zone_a_alert")
                        result["priority"] = "critical"
                    elif item.zone == "B":
                        result["actions_taken"].append("critical_zone_b_alert")
                        result["priority"] = "high"
                    elif item.zone == "C":
                        result["actions_taken"].append("critical_zone_c_alert")
                        result["priority"] = "medium"
                    else:
                        result["actions_taken"].append("critical_unknown_zone_alert")
                        result["priority"] = "high"
                elif item.quantity <= item.reorder_threshold * 0.5:
                    result["actions_taken"].append("low_stock_warning")
                    result["priority"] = "medium"
                elif item.quantity <= item.reorder_threshold * 0.75:
                    result["actions_taken"].append("stock_watch")
                    result["priority"] = "low"
                else:
                    result["actions_taken"].append("approaching_reorder_point")

    elif action == "transfer":
        if target_warehouse is None:
            return {"success": False, "error": "Target warehouse required for transfer"}
        else:
            if quantity_change <= 0:
                return {"success": False, "error": "Transfer quantity must be positive"}
            elif quantity_change > item.quantity:
                return {"success": False, "error": "Insufficient stock for transfer"}
            else:
                # Validate target warehouse
                target_wh = None
                for wh in WAREHOUSE_DATA:
                    if wh.id == target_warehouse:
                        target_wh = wh
                        break

                if target_wh is None:
                    return {"success": False, "error": "Invalid target warehouse"}
                else:
                    # Check zone compatibility
                    target_zone = data.get("target_zone", item.zone)
                    if target_zone in target_wh.zones:
                        item.quantity -= quantity_change
                        result["actions_taken"].append("transfer_completed")
                        result["transferred_to"] = target_warehouse
                        result["transferred_quantity"] = quantity_change
                        result["target_zone"] = target_zone

                        # Check if source needs reorder after transfer
                        if item.quantity <= item.reorder_threshold:
                            if item.quantity == 0:
                                result["actions_taken"].append("source_depleted_after_transfer")
                                result["source_reorder_needed"] = True
                            elif item.quantity <= item.reorder_threshold * 0.5:
                                result["actions_taken"].append("source_low_after_transfer")
                                result["source_reorder_needed"] = True
                            else:
                                result["actions_taken"].append("source_below_threshold_after_transfer")
                    else:
                        return {"success": False, "error": "Target zone not available in destination warehouse"}

    elif action == "write_off":
        if reason == "damaged":
            if quantity_change > item.quantity:
                item.quantity = 0
                result["actions_taken"].append("full_write_off_damaged")
            else:
                item.quantity -= quantity_change
                result["actions_taken"].append("partial_write_off_damaged")
        elif reason == "expired":
            item.quantity = 0
            result["actions_taken"].append("full_write_off_expired")
        elif reason == "recall":
            item.quantity = 0
            result["actions_taken"].append("full_write_off_recall")
            result["requires_supplier_notification"] = True
        else:
            if quantity_change > 0:
                item.quantity = max(0, item.quantity - quantity_change)
                result["actions_taken"].append("write_off_other")
            else:
                return {"success": False, "error": "Write-off quantity required"}
    else:
        return {"success": False, "error": "Unknown action: {}".format(action)}

    result["new_quantity"] = item.quantity
    return result


def get_inventory_summary():
    """Get inventory summary grouped by warehouse."""
    summary = {}
    for item in INVENTORY_DATA:
        wh = item.warehouse
        if wh not in summary:
            summary[wh] = {
                "warehouse": wh,
                "total_items": 0,
                "total_quantity": 0,
                "total_value": 0.0,
                "low_stock_count": 0,
                "out_of_stock_count": 0,
                "categories": {},
            }

        summary[wh]["total_items"] += 1
        summary[wh]["total_quantity"] += item.quantity
        summary[wh]["total_value"] += item.quantity * item.unit_price

        if item.quantity == 0:
            summary[wh]["out_of_stock_count"] += 1
        elif item.quantity <= item.reorder_threshold:
            summary[wh]["low_stock_count"] += 1

        cat = item.category
        if cat not in summary[wh]["categories"]:
            summary[wh]["categories"][cat] = {"count": 0, "quantity": 0, "value": 0.0}
        summary[wh]["categories"][cat]["count"] += 1
        summary[wh]["categories"][cat]["quantity"] += item.quantity
        summary[wh]["categories"][cat]["value"] += item.quantity * item.unit_price

    return {"warehouses": list(summary.values())}


def calculate_reorder_recommendations():
    """Calculate which items need reordering and recommended quantities."""
    recommendations = []
    for item in INVENTORY_DATA:
        if item.quantity <= item.reorder_threshold:
            supplier = _find_supplier(item.supplier_id)
            reorder_qty = item.reorder_threshold * 2

            rec = {
                "item_id": item.id,
                "name": item.name,
                "sku": item.sku,
                "current_quantity": item.quantity,
                "reorder_threshold": item.reorder_threshold,
                "recommended_quantity": reorder_qty,
                "estimated_cost": reorder_qty * item.unit_price,
                "warehouse": item.warehouse,
            }

            if supplier:
                rec["supplier"] = supplier.name
                rec["lead_time_days"] = supplier.lead_time_days
                rec["supplier_reliability"] = supplier.reliability_score

            recommendations.append(rec)

    return recommendations


def _find_supplier(supplier_id):
    """Find supplier by ID."""
    for supplier in SUPPLIER_DATA:
        if supplier.id == supplier_id:
            return supplier
    return None
