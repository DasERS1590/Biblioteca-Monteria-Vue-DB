import { Navigate } from "react-router-dom";
import { useState } from "react";
import { useEffect } from "react";

const ProtectedRoute = ({ children, allowedRoles }) => {
  const [isLoading, setIsLoading] = useState(true); 
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedUser = JSON.parse(localStorage.getItem("user"));
    if (storedUser) {
      setUser(storedUser);  // Establece el usuario desde localStorage
    }
    setIsLoading(false); // Marca como cargado después de obtener los datos
  }, []);

  if (isLoading) {
    return <p>Cargando...</p>; // Muestra un mensaje de "Cargando..." mientras se obtienen los datos
  }

  if (!user) {
    alert("Acceso denegado");
    return <Navigate to="/login" />;
  }

  if (allowedRoles && !allowedRoles.includes(user.rol)) {
    alert("No tiene Permiso");
    return <Navigate to="/login" />;
  }

  return children; // Si el usuario es válido, muestra los componentes hijos
};

export default ProtectedRoute;