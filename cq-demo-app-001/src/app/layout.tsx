import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "CQ Demo App 001 — E-Commerce Dashboard",
  description: "Code quality demo app with intentional violations for scanner detection",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
