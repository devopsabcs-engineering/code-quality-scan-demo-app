import { NextResponse } from "next/server";

export async function GET() {
  return NextResponse.json({
    status: "healthy",
    app: "cq-demo-app-001",
    version: "1.0.0",
    timestamp: new Date().toISOString(),
  });
}
