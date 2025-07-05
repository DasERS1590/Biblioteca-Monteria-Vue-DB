import React, { useState } from "react";
import { authService } from '../services/authService';
import "../styles/user/login.css";

const Login = () => {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMessage("");
    

    
    try {
      const response = await authService.login({ correo: email, contrasena: password });

      
      // Verificar que la respuesta tenga la estructura correcta
      if (response && response.data) {
        const { token, id, nombre, rol } = response.data;
        
        if (token && id) {
          // Crear objeto usuario con la estructura esperada
          const user = {
            id: id,
            nombre: nombre,
            rol: rol
          };
          
          // Guardar token y usuario en localStorage
          localStorage.setItem("token", token);
          localStorage.setItem("user", JSON.stringify(user));
          
          window.location.href = "/";
        } else {
          setErrorMessage("Respuesta del servidor incompleta");
        }
      } else {
        setErrorMessage("Respuesta del servidor inválida");
      }
    } catch (error) {
      setError('Error al iniciar sesión. Verifica tus credenciales.');
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <h2 className="login-header">Biblioteca</h2>
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <input
              type="email"
              placeholder="Correo"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="input-group">
            <input
              type="password"
              placeholder="Contraseña"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <button type="submit" className="login-btn">
            Iniciar Sesión
          </button>
          <div className="text">
            <div className="create-acount">
              <a href="/register">Crear Cuenta</a>
            </div>
            <div className="forgot-password">
              <a href="#">¿Olvidaste tu contraseña?</a>
            </div>
          </div>

          {errorMessage && (
            <div className="error-message">
              <h3>{errorMessage}</h3>
            </div>
          )}
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
  
}

*/