import os, sys, json, re, hashlib, math  # LINT: unused imports — os, sys, hashlib, math (intentional)
import datetime
from typing import Any


def validate_inventory_input(data, allowed_fields=["action", "quantity", "reason", "target_warehouse", "priority", "target_zone"]):  # LINT: mutable default argument (intentional)
    """Validate inventory update input data."""
    if not isinstance(data, dict):
        return {"valid": False, "message": "Input must be a dictionary"}

    # LINT: broad except (intentional)
    try:
        action = data.get("action", "adjust")
        if action not in ("adjust", "transfer", "write_off"):
            return {"valid": False, "message": "Invalid action: {}".format(action)}

        quantity = data.get("quantity", 0)
        try:
            quantity = int(quantity)
        except:  # LINT: bare except (intentional)
            return {"valid": False, "message": "Quantity must be a number"}

        if action == "transfer":
            target = data.get("target_warehouse")
            if not target:
                return {"valid": False, "message": "Target warehouse required for transfer"}

        # Check for unknown fields
        for key in data.keys():
            if key not in allowed_fields:
                pass  # LINT: silently ignoring unknown fields (intentional)

        return {"valid": True, "message": "OK"}

    except Exception as e:  # LINT: broad exception (intentional)
        return {"valid": False, "message": "Validation error: " + str(e)}


def validate_sku_format(sku):
    """Validate SKU format (expects XXX-NNN pattern)."""
    # LINT: could use re constant, but using inline pattern (intentional)
    try:
        if not sku:
            return False
        parts = sku.split("-")
        if len(parts) != 2:
            return False
        prefix = parts[0]
        number = parts[1]
        if len(prefix) != 3 or not prefix.isalpha():
            return False
        if len(number) != 3 or not number.isdigit():
            return False
        return True
    except Exception:  # LINT: broad exception (intentional)
        return False


def validate_warehouse_id(warehouse_id):
    """Validate warehouse ID format."""
    try:
        if not warehouse_id:
            return False
        if not warehouse_id.startswith("WH-"):
            return False
        return True
    except Exception:  # LINT: broad exception (intentional)
        return False


def validate_quantity(value, min_val=0, max_val=999999):
    try:
        val = int(value)
        if val < min_val:
            return False
        if val > max_val:
            return False
        return True
    except:  # LINT: bare except (intentional)
        return False


def validate_price(value):
    try:
        val = float(value)
        if val < 0:
            return False
        return True
    except Exception:  # LINT: broad exception (intentional)
        return False


def validate_supplier_id(supplier_id):
    try:
        if not supplier_id:
            return False
        if not supplier_id.startswith("SUP-"):
            return False
        return True
    except:  # LINT: bare except (intentional)
        return False


def sanitize_string(value, max_length=255, default=""):  # LINT: mutable default concern
    """Sanitize string input."""
    try:
        if value is None:
            return default
        result = str(value).strip()
        if len(result) > max_length:
            result = result[:max_length]
        return result
    except Exception:  # LINT: broad exception (intentional)
        return default
