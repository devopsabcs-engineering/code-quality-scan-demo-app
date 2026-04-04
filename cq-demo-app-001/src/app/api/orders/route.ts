import { NextResponse } from "next/server";

// VIOLATION: unused variable (no-unused-vars)
let orderCache: any[] = [];

// VIOLATION: unused variable
let maxRetries = 3;

const orders = [
  { id: "ORD-1042", customer: "Alice Johnson", amount: 234.50, status: "shipped", date: "2024-12-15", items: 3 },
  { id: "ORD-1041", customer: "Bob Smith", amount: 89.99, status: "processing", date: "2024-12-14", items: 1 },
  { id: "ORD-1040", customer: "Carol White", amount: 567.00, status: "delivered", date: "2024-12-13", items: 5 },
  { id: "ORD-1039", customer: "Dave Brown", amount: 45.00, status: "cancelled", date: "2024-12-12", items: 1 },
  { id: "ORD-1038", customer: "Eve Davis", amount: 199.95, status: "shipped", date: "2024-12-11", items: 2 },
  { id: "ORD-1037", customer: "Frank Wilson", amount: 312.00, status: "delivered", date: "2024-12-10", items: 4 },
  { id: "ORD-1036", customer: "Grace Lee", amount: 78.50, status: "processing", date: "2024-12-09", items: 1 },
  { id: "ORD-1035", customer: "Hank Moore", amount: 425.00, status: "delivered", date: "2024-12-08", items: 6 },
];

export async function GET() {
  return NextResponse.json({
    orders: orders,
    total: orders.length,
    summary: {
      totalRevenue: orders.reduce((sum, o) => sum + o.amount, 0),
      avgOrderValue: orders.reduce((sum, o) => sum + o.amount, 0) / orders.length,
    },
  });
}
