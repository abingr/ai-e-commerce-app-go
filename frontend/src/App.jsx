import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

function App() {
  const [health, setHealth] = useState("checking");
  const [products, setProducts] = useState([]);
  const [productsStatus, setProductsStatus] = useState("loading");
  const [authMode, setAuthMode] = useState("register");
  const [authForm, setAuthForm] = useState({
    name: "Backend Learner",
    email: "learner@example.com",
    password: "password123",
  });
  const [authStatus, setAuthStatus] = useState("signed out");
  const [currentUser, setCurrentUser] = useState(null);

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

    const token = localStorage.getItem("auth_token");
    if (token) {
      fetchCurrentUser(token);
    }
  }, []);

  function updateAuthForm(event) {
    setAuthForm({
      ...authForm,
      [event.target.name]: event.target.value,
    });
  }

  function submitAuth(event) {
    event.preventDefault();

    const endpoint =
      authMode === "register"
        ? "http://localhost:8080/api/v1/auth/register"
        : "http://localhost:8080/api/v1/auth/login";

    const payload =
      authMode === "register"
        ? authForm
        : { email: authForm.email, password: authForm.password };

    setAuthStatus("checking");

    fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })
      .then(async (response) => {
        const body = await response.json();
        if (!response.ok) {
          throw new Error(body.error ?? "authentication failed");
        }

        localStorage.setItem("auth_token", body.data.token);
        setCurrentUser(body.data.user);
        setAuthStatus("signed in");
      })
      .catch((error) => setAuthStatus(error.message));
  }

  function fetchCurrentUser(token) {
    fetch("http://localhost:8080/api/v1/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(async (response) => {
        const body = await response.json();
        if (!response.ok) {
          throw new Error(body.error ?? "session expired");
        }

        setCurrentUser(body.data);
        setAuthStatus("signed in");
      })
      .catch(() => {
        localStorage.removeItem("auth_token");
        setCurrentUser(null);
        setAuthStatus("signed out");
      });
  }

  function signOut() {
    localStorage.removeItem("auth_token");
    setCurrentUser(null);
    setAuthStatus("signed out");
  }

  return (
    <main className="app-shell">
      <header className="page-header">
        <div>
          <p className="eyebrow">Electronics Commerce</p>
          <h1>Backend Learning Storefront</h1>
          <p>
            Phase 3 combines the product catalog with JWT authentication.
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
              : "Phase 3 adds JWT authentication to the backend."}
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
                className={authMode === "register" ? "active" : ""}
                onClick={() => setAuthMode("register")}
              >
                Register
              </button>
              <button
                type="button"
                className={authMode === "login" ? "active" : ""}
                onClick={() => setAuthMode("login")}
              >
                Login
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
