import React, { useState, useEffect } from "react";
import { Navigate } from "react-router-dom";

const ProtectedRoute = ({ children, allowedRoles }) => {
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedUser = JSON.parse(localStorage.getItem("user"));
    if (storedUser) {
      setUser(storedUser);
    }
    setIsLoading(false);
  }, []);

  if (isLoading) {
    return <p>Cargando...</p>;
  }

  if (!user) {
    alert("Acceso denegado. Debes iniciar sesión.");
    return <Navigate to="/login" />;
  }

  if (allowedRoles && !allowedRoles.includes(user.rol)) {
    // Determinar la ruta correcta basada en el rol del usuario
    let redirectPath = "/";
    if (user.rol === "administrador") {
      redirectPath = "/admin/dashboard";
    } else if (user.rol === "usuario") {
      redirectPath = "/user/dashboard";
    }
    
    alert(`No tienes permisos para acceder a esta sección. Redirigiendo a tu dashboard.`);
    return <Navigate to={redirectPath} />;
  }
  
  return children;
};

export default ProtectedRoute;