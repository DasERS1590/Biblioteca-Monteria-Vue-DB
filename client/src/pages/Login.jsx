// src/pages/Login.jsx
import React, { useState } from "react";
import "../styles/login.css"; // Importamos los estilos
import { useNavigate } from "react-router-dom";

const Login = () => {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      // Hacer la solicitud al backend
      const response = await fetch("http://localhost:4000/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ correo: email, contrasena: password }),
      });

      if (!response.ok) {
        // Si la respuesta no es OK, mostrar un error
        const errorData = await response.json();
        setErrorMessage(errorData.message || "Login failed.");
        return;
      }

      const data = await response.json();

      // para guardar en local store 
      localStorage.setItem("user", JSON.stringify(data));
      console.log("Usuario autenticado:", data);

      if (data.rol === "administrador") {
        navigate("/admin/dashboard");
      } else {
        navigate("/user/dashboard");
      }
    } catch (error) {
      console.error("Error during login:", error );
      setErrorMessage("An error occurred. Please try again.");
    }
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




/*

para recuperar el localestorage
const user = JSON.parse(localStorage.getItem("user"));

if (user) {
  console.log("Usuario logueado:", user);
  console.log("Rol del usuario:", user.rol);
}

*/