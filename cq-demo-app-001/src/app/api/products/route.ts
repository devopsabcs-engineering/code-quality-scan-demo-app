import { NextResponse } from "next/server";

// VIOLATION: unused variable (no-unused-vars)
let defaultCategory = "all";

const products = [
  { id: 1, name: "Wireless Headphones", price: 79.99, stock: 142, category: "Electronics" },
  { id: 2, name: "Running Shoes", price: 129.95, stock: 56, category: "Sports" },
  { id: 3, name: "Coffee Maker", price: 49.99, stock: 230, category: "Kitchen" },
  { id: 4, name: "Backpack Pro", price: 89.00, stock: 0, category: "Accessories" },
  { id: 5, name: "Desk Lamp LED", price: 34.50, stock: 89, category: "Office" },
  { id: 6, name: "Yoga Mat", price: 24.99, stock: 310, category: "Sports" },
  { id: 7, name: "Bluetooth Speaker", price: 59.99, stock: 67, category: "Electronics" },
  { id: 8, name: "Water Bottle", price: 19.95, stock: 420, category: "Sports" },
];

export async function GET() {
  return NextResponse.json({
    products: products,
    total: products.length,
  });
}
