from dataclasses import dataclass
from typing import Optional
import datetime


@dataclass
class InventoryItem:
    id: int
    name: str
    sku: str
    quantity: int
    warehouse: str
    zone: str
    reorder_threshold: int
    supplier_id: str
    unit_price: float
    category: str
    last_updated: Optional[str] = None


@dataclass
class Warehouse:
    id: str
    name: str
    location: str
    capacity: int
    zones: list


@dataclass
class Supplier:
    id: str
    name: str
    lead_time_days: int
    reliability_score: float
    region: str


# In-memory data store (intentional: no database)
INVENTORY_DATA = [
    InventoryItem(1, "Widget A", "WGT-001", 150, "WH-EAST", "A", 50, "SUP-001", 12.99, "electronics"),
    InventoryItem(2, "Widget B", "WGT-002", 30, "WH-EAST", "B", 100, "SUP-001", 24.50, "electronics"),
    InventoryItem(3, "Gadget X", "GDG-001", 5, "WH-WEST", "A", 20, "SUP-002", 89.99, "hardware"),
    InventoryItem(4, "Gadget Y", "GDG-002", 200, "WH-WEST", "C", 75, "SUP-002", 45.00, "hardware"),
    InventoryItem(5, "Part Alpha", "PRT-001", 0, "WH-NORTH", "A", 10, "SUP-003", 5.25, "components"),
    InventoryItem(6, "Part Beta", "PRT-002", 500, "WH-NORTH", "B", 150, "SUP-003", 3.10, "components"),
    InventoryItem(7, "Assembly K", "ASM-001", 12, "WH-EAST", "C", 25, "SUP-004", 199.99, "assemblies"),
    InventoryItem(8, "Assembly L", "ASM-002", 80, "WH-SOUTH", "A", 30, "SUP-004", 149.50, "assemblies"),
    InventoryItem(9, "Module M", "MOD-001", 45, "WH-SOUTH", "B", 60, "SUP-005", 75.00, "modules"),
    InventoryItem(10, "Module N", "MOD-002", 3, "WH-EAST", "A", 15, "SUP-005", 110.00, "modules"),
]

WAREHOUSE_DATA = [
    Warehouse("WH-EAST", "East Distribution Center", "New York", 10000, ["A", "B", "C"]),
    Warehouse("WH-WEST", "West Distribution Center", "Los Angeles", 15000, ["A", "B", "C", "D"]),
    Warehouse("WH-NORTH", "North Warehouse", "Chicago", 8000, ["A", "B"]),
    Warehouse("WH-SOUTH", "South Warehouse", "Houston", 12000, ["A", "B", "C"]),
]

SUPPLIER_DATA = [
    Supplier("SUP-001", "TechParts Inc", 7, 0.95, "domestic"),
    Supplier("SUP-002", "Global Hardware Ltd", 14, 0.88, "international"),
    Supplier("SUP-003", "Component World", 3, 0.99, "domestic"),
    Supplier("SUP-004", "Assembly Masters", 10, 0.92, "domestic"),
    Supplier("SUP-005", "ModuleTech", 21, 0.85, "international"),
]
