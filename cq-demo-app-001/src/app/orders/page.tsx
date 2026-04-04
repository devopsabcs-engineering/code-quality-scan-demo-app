import Link from "next/link";

const orders = [
  { id: "ORD-1042", customer: "Alice Johnson", amount: 234.50, status: "shipped", date: "2024-12-15" },
  { id: "ORD-1041", customer: "Bob Smith", amount: 89.99, status: "processing", date: "2024-12-14" },
  { id: "ORD-1040", customer: "Carol White", amount: 567.00, status: "delivered", date: "2024-12-13" },
  { id: "ORD-1039", customer: "Dave Brown", amount: 45.00, status: "cancelled", date: "2024-12-12" },
  { id: "ORD-1038", customer: "Eve Davis", amount: 199.95, status: "shipped", date: "2024-12-11" },
  { id: "ORD-1037", customer: "Frank Wilson", amount: 312.00, status: "delivered", date: "2024-12-10" },
  { id: "ORD-1036", customer: "Grace Lee", amount: 78.50, status: "processing", date: "2024-12-09" },
  { id: "ORD-1035", customer: "Hank Moore", amount: 425.00, status: "delivered", date: "2024-12-08" },
];

function getStatusBadge(status: string) {
  switch (status) {
    case "shipped":
      return "badge-success";
    case "processing":
      return "badge-warning";
    case "delivered":
      return "badge-success";
    case "cancelled":
      return "badge-danger";
    default:
      return "";
  }
}

export default function OrdersPage() {
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
        <h2>Orders Management</h2>
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>Order ID</th>
                <th>Customer</th>
                <th>Amount</th>
                <th>Date</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td>#{order.id}</td>
                  <td>{order.customer}</td>
                  <td>${order.amount.toFixed(2)}</td>
                  <td>{order.date}</td>
                  <td>
                    <span className={`badge ${getStatusBadge(order.status)}`}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </>
  );
}
