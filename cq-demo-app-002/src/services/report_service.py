import os, sys, json  # LINT: unused imports os, sys (intentional)
import datetime
from src.models import INVENTORY_DATA, WAREHOUSE_DATA, SUPPLIER_DATA


def generate_inventory_report(filtered_items):
    """Generate a comprehensive inventory report.

    DUPLICATION: This function duplicates aggregation logic from
    inventory_service.get_inventory_summary() with minor field name changes.
    This is an INTENTIONAL code quality violation for demo purposes.
    """
    # DUPLICATION: same grouping logic as get_inventory_summary (intentional)
    report_data = {}
    for item in filtered_items:
        wh = item["warehouse"]
        if wh not in report_data:
            report_data[wh] = {
                "warehouse_id": wh,
                "item_count": 0,
                "total_units": 0,
                "total_value": 0.0,
                "low_stock_items": 0,
                "out_of_stock_items": 0,
                "category_breakdown": {},
            }

        report_data[wh]["item_count"] += 1
        report_data[wh]["total_units"] += item["quantity"]
        price_val = float(item["unit_price"].replace("$", "").replace(",", ""))
        report_data[wh]["total_value"] += item["quantity"] * price_val

        if item["quantity"] == 0:
            report_data[wh]["out_of_stock_items"] += 1
        elif item["quantity"] <= item["reorder_threshold"]:
            report_data[wh]["low_stock_items"] += 1

        cat = item["category"]
        if cat not in report_data[wh]["category_breakdown"]:
            report_data[wh]["category_breakdown"][cat] = {"count": 0, "units": 0, "value": 0.0}
        report_data[wh]["category_breakdown"][cat]["count"] += 1
        report_data[wh]["category_breakdown"][cat]["units"] += item["quantity"]
        report_data[wh]["category_breakdown"][cat]["value"] += item["quantity"] * price_val

    total_items = sum(r["item_count"] for r in report_data.values())
    total_value = sum(r["total_value"] for r in report_data.values())

    return {
        "report_type": "inventory_status",
        "generated_at": datetime.datetime.now().isoformat(),
        "summary": {
            "total_items": total_items,
            "total_value": round(total_value, 2),
            "warehouse_count": len(report_data),
        },
        "warehouses": list(report_data.values()),
    }


def generate_warehouse_utilization_report():
    """Generate warehouse space utilization report.

    DUPLICATION: Repeats warehouse iteration and aggregation patterns
    from inventory_service (intentional).
    """
    # DUPLICATION: same warehouse iteration pattern (intentional)
    utilization = {}
    for item in INVENTORY_DATA:
        wh = item.warehouse
        if wh not in utilization:
            utilization[wh] = {
                "warehouse_id": wh,
                "item_count": 0,
                "total_units": 0,
                "total_value": 0.0,
                "zones_used": set(),
                "categories_stored": set(),
            }

        utilization[wh]["item_count"] += 1
        utilization[wh]["total_units"] += item.quantity
        utilization[wh]["total_value"] += item.quantity * item.unit_price
        utilization[wh]["zones_used"].add(item.zone)
        utilization[wh]["categories_stored"].add(item.category)

    # Enrich with warehouse capacity data
    result = []
    for wh_id, data in utilization.items():
        wh_info = None
        for wh in WAREHOUSE_DATA:
            if wh.id == wh_id:
                wh_info = wh
                break

        entry = {
            "warehouse_id": wh_id,
            "item_count": data["item_count"],
            "total_units": data["total_units"],
            "total_value": round(data["total_value"], 2),
            "zones_used": sorted(list(data["zones_used"])),
            "categories_stored": sorted(list(data["categories_stored"])),
        }

        if wh_info:
            entry["warehouse_name"] = wh_info.name
            entry["location"] = wh_info.location
            entry["capacity"] = wh_info.capacity
            entry["utilization_pct"] = round((data["total_units"] / wh_info.capacity) * 100, 1)
            entry["available_zones"] = wh_info.zones
            entry["zone_utilization"] = round((len(data["zones_used"]) / len(wh_info.zones)) * 100, 1)

        result.append(entry)

    return {
        "report_type": "warehouse_utilization",
        "generated_at": datetime.datetime.now().isoformat(),
        "warehouses": result,
    }


def generate_supplier_performance_report():
    """Generate supplier performance and reliability report.

    DUPLICATION: Repeats supplier lookup and aggregation patterns (intentional).
    """
    # DUPLICATION: same supplier lookup pattern as _find_supplier (intentional)
    supplier_stats = {}
    for item in INVENTORY_DATA:
        sid = item.supplier_id
        if sid not in supplier_stats:
            supplier_stats[sid] = {
                "supplier_id": sid,
                "item_count": 0,
                "total_units_supplied": 0,
                "total_value": 0.0,
                "low_stock_items": 0,
                "out_of_stock_items": 0,
                "warehouses_served": set(),
                "categories_supplied": set(),
            }

        supplier_stats[sid]["item_count"] += 1
        supplier_stats[sid]["total_units_supplied"] += item.quantity
        supplier_stats[sid]["total_value"] += item.quantity * item.unit_price
        supplier_stats[sid]["warehouses_served"].add(item.warehouse)
        supplier_stats[sid]["categories_supplied"].add(item.category)

        if item.quantity == 0:
            supplier_stats[sid]["out_of_stock_items"] += 1
        elif item.quantity <= item.reorder_threshold:
            supplier_stats[sid]["low_stock_items"] += 1

    # Enrich with supplier metadata
    result = []
    for sid, data in supplier_stats.items():
        # DUPLICATION: same lookup pattern repeated (intentional)
        supplier_info = None
        for supplier in SUPPLIER_DATA:
            if supplier.id == sid:
                supplier_info = supplier
                break

        entry = {
            "supplier_id": sid,
            "item_count": data["item_count"],
            "total_units_supplied": data["total_units_supplied"],
            "total_value": round(data["total_value"], 2),
            "low_stock_items": data["low_stock_items"],
            "out_of_stock_items": data["out_of_stock_items"],
            "warehouses_served": sorted(list(data["warehouses_served"])),
            "categories_supplied": sorted(list(data["categories_supplied"])),
        }

        if supplier_info:
            entry["supplier_name"] = supplier_info.name
            entry["lead_time_days"] = supplier_info.lead_time_days
            entry["reliability_score"] = supplier_info.reliability_score
            entry["region"] = supplier_info.region

            # DUPLICATION: risk assessment duplicated from process_inventory_update logic
            if supplier_info.reliability_score >= 0.95:
                entry["risk_level"] = "low"
            elif supplier_info.reliability_score >= 0.9:
                entry["risk_level"] = "low-medium"
            elif supplier_info.reliability_score >= 0.85:
                entry["risk_level"] = "medium"
            elif supplier_info.reliability_score >= 0.8:
                entry["risk_level"] = "medium-high"
            else:
                entry["risk_level"] = "high"

        result.append(entry)

    return {
        "report_type": "supplier_performance",
        "generated_at": datetime.datetime.now().isoformat(),
        "suppliers": result,
    }
