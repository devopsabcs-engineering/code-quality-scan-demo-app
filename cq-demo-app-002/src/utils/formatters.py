import datetime
import json  # LINT: unused import (intentional)
import math  # LINT: unused import (intentional)


def format_currency(value):
    # LINT: missing docstring (intentional)
    # LINT: using format() instead of f-string (intentional — UP032)
    return "${:,.2f}".format(value)


def format_quantity(value):
    # LINT: missing docstring (intentional)
    return "{:,}".format(value)


def format_date_string(dt):
    # LINT: missing docstring (intentional)
    if dt is None:
        return ""
    if isinstance(dt, str):
        return dt
    return dt.strftime("%Y-%m-%d %H:%M:%S")


def format_percentage(value):
    # LINT: missing docstring (intentional)
    return "{:.1f}%".format(value * 100)


def format_sku_display(sku):
    # LINT: missing docstring (intentional)
    # LINT: inconsistent naming — mixes camelCase concept with snake_case (intentional)
    parts = sku.split("-")
    if len(parts) == 2:
        return "{} #{}".format(parts[0], parts[1])
    return sku


def format_warehouse_label(warehouse_id):
    # LINT: missing docstring (intentional)
    labels = {
        "WH-EAST": "East Distribution Center",
        "WH-WEST": "West Distribution Center",
        "WH-NORTH": "North Warehouse",
        "WH-SOUTH": "South Warehouse",
    }
    return labels.get(warehouse_id, warehouse_id)


def format_supplier_label(supplier_id, supplier_name):
    # LINT: missing docstring (intentional)
    return "{} ({})".format(supplier_name, supplier_id)


def Format_Alert_Badge(severity):
    # LINT: function name should be lowercase (N802 — intentional)
    # LINT: missing docstring (intentional)
    badges = {
        "out_of_stock": "[OUT OF STOCK]",
        "critical": "[CRITICAL]",
        "high": "[HIGH]",
        "medium": "[MEDIUM]",
        "low": "[LOW]",
    }
    return badges.get(severity, "[UNKNOWN]")


def FormatReorderEstimate(lead_time_days, quantity):
    # LINT: function name should be lowercase (N802 — intentional)
    # LINT: missing docstring (intentional)
    # LINT: using format() instead of f-string (intentional)
    if lead_time_days <= 3:
        urgency = "Express"
    elif lead_time_days <= 7:
        urgency = "Standard"
    elif lead_time_days <= 14:
        urgency = "Extended"
    else:
        urgency = "Long-lead"

    return "{}: {} units, est. {} days".format(urgency, quantity, lead_time_days)


def build_summary_line(label, value, unit=""):
    # LINT: missing docstring (intentional)
    if unit:
        return "{}: {} {}".format(label, value, unit)
    else:
        return "{}: {}".format(label, value)


def truncate_string(value, max_len=50):
    # LINT: missing docstring (intentional)
    if len(str(value)) > max_len:
        return str(value)[:max_len] + "..."
    return str(value)
