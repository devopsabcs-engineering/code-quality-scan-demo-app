import Link from "next/link";

// VIOLATION: unused interface (no-unused-vars)
interface ProductFilter {
  category: string;
  minPrice: number;
  maxPrice: number;
}

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

export default function ProductsPage() {
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
        <h2>Products Catalog</h2>
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Price</th>
                <th>Stock</th>
                <th>Category</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id}>
                  <td>{product.id}</td>
                  <td>{product.name}</td>
                  <td>${product.price.toFixed(2)}</td>
                  <td>{product.stock}</td>
                  <td>{product.category}</td>
                  <td>
                    <span
                      className={`badge ${
                        product.stock > 50
                          ? "badge-success"
                          : product.stock > 0
                          ? "badge-warning"
                          : "badge-danger"
                      }`}
                    >
                      {product.stock > 50
                        ? "In Stock"
                        : product.stock > 0
                        ? "Low Stock"
                        : "Out of Stock"}
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
