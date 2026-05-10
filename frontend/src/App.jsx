import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

function App() {
  const [health, setHealth] = useState("checking");
  const [products, setProducts] = useState([]);
  const [productsStatus, setProductsStatus] = useState("loading");

  useEffect(() => {
    fetch("http://localhost:8080/health")
      .then((response) => response.json())
      .then((data) => setHealth(data.status))
      .catch(() => setHealth("offline"));

    fetch("http://localhost:8080/api/v1/products")
      .then((response) => response.json())
      .then((data) => {
        setProducts(data.data ?? []);
        setProductsStatus("ready");
      })
      .catch(() => setProductsStatus("offline"));
  }, []);

  return (
    <main className="app-shell">
      <header className="page-header">
        <div>
          <p className="eyebrow">Electronics Commerce</p>
          <h1>Backend Learning Storefront</h1>
          <p>
            Phase 2 reads product catalog data from a Go Gin API backed by
            PostgreSQL.
          </p>
        </div>
        <div className="status-panel">
          <div className="status-row">
            <span>API</span>
            <strong>{health}</strong>
          </div>
          <div className="status-row">
            <span>Products</span>
            <strong>{productsStatus}</strong>
          </div>
        </div>
      </header>

      <section className="product-grid" aria-label="Products">
        {products.map((product) => (
          <article className="product-card" key={product.id}>
            <img src={product.image_url} alt={product.name} />
            <div className="product-body">
              <p className="category">{product.category}</p>
              <h2>{product.name}</h2>
              <p>{product.description}</p>
              <div className="product-meta">
                <span>{product.brand}</span>
                <strong>{formatPrice(product.price_cents)}</strong>
              </div>
              <div className="stock-row">{product.stock_quantity} in stock</div>
            </div>
          </article>
        ))}
      </section>
    </main>
  );
}

function formatPrice(priceCents) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(priceCents / 100);
}

createRoot(document.getElementById("root")).render(<App />);
