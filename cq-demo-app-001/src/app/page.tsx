import Link from "next/link";

export default function Home() {
  return (
    <>
      <header className="header">
        <h1>📦 E-Commerce Dashboard</h1>
        <nav className="nav">
          <Link href="/">Home</Link>
          <Link href="/products">Products</Link>
          <Link href="/orders">Orders</Link>
        </nav>
      </header>
      <main className="container">
        <h2>Dashboard Overview</h2>
        <div className="grid">
          <div className="card stat">
            <div className="value">1,247</div>
            <div className="label">Total Products</div>
          </div>
          <div className="card stat">
            <div className="value">389</div>
            <div className="label">Active Orders</div>
          </div>
          <div className="card stat">
            <div className="value">$52,480</div>
            <div className="label">Revenue (MTD)</div>
          </div>
          <div className="card stat">
            <div className="value">94.2%</div>
            <div className="label">Fulfillment Rate</div>
          </div>
        </div>

        <h2 style={{ marginTop: "2rem" }}>Recent Activity</h2>
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>Order ID</th>
                <th>Customer</th>
                <th>Amount</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>#ORD-1042</td>
                <td>Alice Johnson</td>
                <td>$234.50</td>
                <td><span className="badge badge-success">Shipped</span></td>
              </tr>
              <tr>
                <td>#ORD-1041</td>
                <td>Bob Smith</td>
                <td>$89.99</td>
                <td><span className="badge badge-warning">Processing</span></td>
              </tr>
              <tr>
                <td>#ORD-1040</td>
                <td>Carol White</td>
                <td>$567.00</td>
                <td><span className="badge badge-success">Delivered</span></td>
              </tr>
              <tr>
                <td>#ORD-1039</td>
                <td>Dave Brown</td>
                <td>$45.00</td>
                <td><span className="badge badge-danger">Cancelled</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </>
  );
}
