import React, { useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

function App() {
  const [health, setHealth] = useState("checking");
  const [products, setProducts] = useState([]);
  const [productsStatus, setProductsStatus] = useState("loading");
  const [category, setCategory] = useState("");
  const [search, setSearch] = useState("");
  const [authMode, setAuthMode] = useState("login");
  const [authForm, setAuthForm] = useState({
    name: "Backend Learner",
    email: "learner@example.com",
    password: "password123",
  });
  const [authStatus, setAuthStatus] = useState("signed out");
  const [currentUser, setCurrentUser] = useState(null);
  const [cart, setCart] = useState({ items: [], total_cents: 0 });
  const [orders, setOrders] = useState([]);
  const [actionStatus, setActionStatus] = useState("Ready");

  const token = localStorage.getItem("auth_token");
  const categories = useMemo(
    () => Array.from(new Set(products.map((product) => product.category))).sort(),
    [products],
  );

  useEffect(() => {
    checkHealth();
    fetchProducts();

    if (token) {
      fetchCurrentUser(token);
      fetchCart(token);
      fetchOrders(token);
    }
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [category, search]);

  function authHeaders(authToken = token) {
    return authToken ? { Authorization: `Bearer ${authToken}` } : {};
  }

  async function apiRequest(path, options = {}) {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(options.headers ?? {}),
      },
    });

    const hasBody = response.status !== 204;
    const body = hasBody ? await response.json() : null;

    if (!response.ok) {
      const details = body?.fields
        ? ` (${body.fields.map((field) => `${field.field}: ${field.rule}`).join(", ")})`
        : "";
      const requestId = body?.request_id ? ` Request ID: ${body.request_id}` : "";
      throw new Error(`${body?.error ?? "request failed"}${details}.${requestId}`);
    }

    return body;
  }

  async function checkHealth() {
    try {
      const body = await apiRequest("/health");
      setHealth(body.status);
    } catch {
      setHealth("offline");
    }
  }

  async function fetchProducts() {
    setProductsStatus("loading");

    const params = new URLSearchParams();
    if (category) {
      params.set("category", category);
    }
    if (search) {
      params.set("search", search);
    }

    try {
      const query = params.toString();
      const body = await apiRequest(`/api/v1/products${query ? `?${query}` : ""}`);
      setProducts(body.data ?? []);
      setProductsStatus("ready");
    } catch (error) {
      setProductsStatus("offline");
      setActionStatus(error.message);
    }
  }

  function updateAuthForm(event) {
    setAuthForm({
      ...authForm,
      [event.target.name]: event.target.value,
    });
  }

  async function submitAuth(event) {
    event.preventDefault();

    const endpoint = authMode === "register" ? "/api/v1/auth/register" : "/api/v1/auth/login";
    const payload =
      authMode === "register"
        ? authForm
        : { email: authForm.email, password: authForm.password };

    setAuthStatus("checking");

    try {
      const body = await apiRequest(endpoint, {
        method: "POST",
        body: JSON.stringify(payload),
      });

      localStorage.setItem("auth_token", body.data.token);
      setCurrentUser(body.data.user);
      setAuthStatus("signed in");
      setActionStatus(`Signed in as ${body.data.user.email}`);
      await Promise.all([fetchCart(body.data.token), fetchOrders(body.data.token)]);
    } catch (error) {
      setAuthStatus(error.message);
    }
  }

  async function fetchCurrentUser(authToken) {
    try {
      const body = await apiRequest("/api/v1/me", {
        headers: authHeaders(authToken),
      });

      setCurrentUser(body.data);
      setAuthStatus("signed in");
    } catch {
      localStorage.removeItem("auth_token");
      setCurrentUser(null);
      setAuthStatus("signed out");
      setCart({ items: [], total_cents: 0 });
      setOrders([]);
    }
  }

  async function fetchCart(authToken = token) {
    if (!authToken) {
      return;
    }

    const body = await apiRequest("/api/v1/cart", {
      headers: authHeaders(authToken),
    });
    setCart(body.data);
  }

  async function fetchOrders(authToken = token) {
    if (!authToken) {
      return;
    }

    const body = await apiRequest("/api/v1/orders", {
      headers: authHeaders(authToken),
    });
    setOrders(body.data ?? []);
  }

  async function addToCart(productID) {
    if (!token) {
      setActionStatus("Login first to add products to your cart.");
      return;
    }

    try {
      const body = await apiRequest("/api/v1/cart/items", {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify({ product_id: productID, quantity: 1 }),
      });
      setCart(body.data);
      setActionStatus("Product added to cart.");
    } catch (error) {
      setActionStatus(error.message);
    }
  }

  async function updateCartItem(productID, quantity) {
    if (quantity < 1) {
      return;
    }

    try {
      const body = await apiRequest(`/api/v1/cart/items/${productID}`, {
        method: "PATCH",
        headers: authHeaders(),
        body: JSON.stringify({ quantity }),
      });
      setCart(body.data);
      setActionStatus("Cart quantity updated.");
    } catch (error) {
      setActionStatus(error.message);
    }
  }

  async function removeCartItem(productID) {
    try {
      const body = await apiRequest(`/api/v1/cart/items/${productID}`, {
        method: "DELETE",
        headers: authHeaders(),
      });
      setCart(body.data);
      setActionStatus("Cart item removed.");
    } catch (error) {
      setActionStatus(error.message);
    }
  }

  async function checkout() {
    if (!cart.items.length) {
      setActionStatus("Cart is empty.");
      return;
    }

    try {
      const body = await apiRequest("/api/v1/orders", {
        method: "POST",
        headers: authHeaders(),
      });
      setCart({ items: [], total_cents: 0 });
      await fetchOrders();
      setActionStatus(`Order ${body.data.id.slice(0, 8)} created.`);
    } catch (error) {
      setActionStatus(error.message);
    }
  }

  function signOut() {
    localStorage.removeItem("auth_token");
    setCurrentUser(null);
    setAuthStatus("signed out");
    setCart({ items: [], total_cents: 0 });
    setOrders([]);
    setActionStatus("Signed out.");
  }

  return (
    <main className="app-shell">
      <header className="page-header">
        <div>
          <p className="eyebrow">Electronics Commerce</p>
          <h1>Backend Learning Storefront</h1>
          <p>Phase 8 connects React to the production-style API features.</p>
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
          <div className="status-row">
            <span>Auth</span>
            <strong>{currentUser ? currentUser.role : authStatus}</strong>
          </div>
        </div>
      </header>

      <section className="auth-panel" aria-label="Authentication">
        <div>
          <p className="eyebrow">Account</p>
          <h2>{currentUser ? currentUser.name : "Register or login"}</h2>
          <p>
            {currentUser
              ? `${currentUser.email} is authenticated as ${currentUser.role}.`
              : "JWT is stored in localStorage and sent as a bearer token for protected requests."}
          </p>
        </div>

        {currentUser ? (
          <button type="button" onClick={signOut}>
            Sign out
          </button>
        ) : (
          <form className="auth-form" onSubmit={submitAuth}>
            <div className="mode-row">
              <button
                type="button"
                className={authMode === "login" ? "active" : ""}
                onClick={() => setAuthMode("login")}
              >
                Login
              </button>
              <button
                type="button"
                className={authMode === "register" ? "active" : ""}
                onClick={() => setAuthMode("register")}
              >
                Register
              </button>
            </div>

            {authMode === "register" && (
              <label>
                Name
                <input name="name" value={authForm.name} onChange={updateAuthForm} />
              </label>
            )}

            <label>
              Email
              <input name="email" value={authForm.email} onChange={updateAuthForm} />
            </label>

            <label>
              Password
              <input
                name="password"
                type="password"
                value={authForm.password}
                onChange={updateAuthForm}
              />
            </label>

            <button type="submit">{authMode === "register" ? "Create account" : "Login"}</button>
            <p className="auth-status">{authStatus}</p>
          </form>
        )}
      </section>

      <section className="workspace-grid">
        <section className="catalog-panel" aria-label="Products">
          <div className="section-heading">
            <div>
              <p className="eyebrow">Catalog</p>
              <h2>Products</h2>
            </div>
            <div className="filters">
              <input
                aria-label="Search products"
                placeholder="Search"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
              />
              <select
                aria-label="Filter by category"
                value={category}
                onChange={(event) => setCategory(event.target.value)}
              >
                <option value="">All categories</option>
                {categories.map((item) => (
                  <option key={item} value={item}>
                    {item}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="product-grid">
            {products.map((product) => (
              <article className="product-card" key={product.id}>
                <img src={product.image_url} alt={product.name} />
                <div className="product-body">
                  <p className="category">{product.category}</p>
                  <h3>{product.name}</h3>
                  <p>{product.description}</p>
                  <div className="product-meta">
                    <span>{product.brand}</span>
                    <strong>{formatPrice(product.price_cents)}</strong>
                  </div>
                  <div className="product-actions">
                    <span>{product.stock_quantity} in stock</span>
                    <button type="button" onClick={() => addToCart(product.id)}>
                      Add
                    </button>
                  </div>
                </div>
              </article>
            ))}
          </div>
        </section>

        <aside className="commerce-panel" aria-label="Cart and orders">
          <section className="side-section">
            <div className="section-heading compact">
              <div>
                <p className="eyebrow">Cart</p>
                <h2>{formatPrice(cart.total_cents)}</h2>
              </div>
              <button type="button" onClick={checkout} disabled={!currentUser || cart.items.length === 0}>
                Checkout
              </button>
            </div>

            <div className="line-items">
              {cart.items.length === 0 ? (
                <p className="empty-state">No cart items yet.</p>
              ) : (
                cart.items.map((item) => (
                  <div className="line-item" key={item.product_id}>
                    <div>
                      <strong>{item.name}</strong>
                      <span>{formatPrice(item.line_total_cents)}</span>
                    </div>
                    <div className="quantity-row">
                      <button
                        type="button"
                        onClick={() => updateCartItem(item.product_id, item.quantity - 1)}
                      >
                        -
                      </button>
                      <span>{item.quantity}</span>
                      <button
                        type="button"
                        onClick={() => updateCartItem(item.product_id, item.quantity + 1)}
                      >
                        +
                      </button>
                      <button type="button" onClick={() => removeCartItem(item.product_id)}>
                        Remove
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </section>

          <section className="side-section">
            <div className="section-heading compact">
              <div>
                <p className="eyebrow">Orders</p>
                <h2>{orders.length}</h2>
              </div>
              <button type="button" onClick={() => fetchOrders()} disabled={!currentUser}>
                Refresh
              </button>
            </div>

            <div className="line-items">
              {orders.length === 0 ? (
                <p className="empty-state">No orders yet.</p>
              ) : (
                orders.slice(0, 5).map((order) => (
                  <div className="order-row" key={order.id}>
                    <div>
                      <strong>{order.id.slice(0, 8)}</strong>
                      <span>{order.status}</span>
                    </div>
                    <span>{formatPrice(order.total_cents)}</span>
                  </div>
                ))
              )}
            </div>
          </section>

          <p className="action-status">{actionStatus}</p>
        </aside>
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
