// src/pages/Login.jsx
import React, { useState } from "react";
import "../styles/login.css"; // Importamos los estilos

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Form submitted with:", { email, password });

    role === "admin" ? navigate("/admin/dashboard") : navigate("/user/dashboard");
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h2 className="login-header">Biblioteca</h2> {/* Título ajustado */}
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="input-group">
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <div className="forgot-password">
            <a href="#">Forgot password?</a>
          </div>
          <button type="submit" className="login-btn">
            Log in
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
