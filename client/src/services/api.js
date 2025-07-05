// Configuración base de la API
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_URL_BACKEND || 'http://localhost:4000/v1/api';

// Crear instancia de axios
export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para agregar token de autenticación
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor para manejar respuestas
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Manejar errores de autenticación
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    
    // Manejar errores de permisos (403)
    if (error.response?.status === 403) {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        // Si es administrador y está en una ruta de admin, no redirigir automáticamente
        if (user.rol === "administrador" && window.location.pathname.startsWith('/admin')) {
          // Solo mostrar el error, no redirigir
          console.warn("Error de permisos en ruta de administrador:", error.response?.data);
        } else {
          // Determinar la ruta correcta basada en el rol del usuario
          let redirectPath = "/";
          if (user.rol === "administrador") {
            redirectPath = "/admin/dashboard";
          } else if (user.rol === "usuario") {
            redirectPath = "/user/dashboard";
          }
          
          // Solo redirigir si no estamos ya en la ruta correcta
          if (window.location.pathname !== redirectPath) {
            alert("No tienes permisos para acceder a esta funcionalidad. Redirigiendo a tu dashboard.");
            window.location.href = redirectPath;
            return Promise.reject(new Error("Redirigiendo por falta de permisos"));
          }
        }
      }
    }
    
    // Crear un error más descriptivo
    const errorMessage = error.response?.data?.error || 
                        error.response?.data?.message || 
                        error.message || 
                        'Error en la petición';
    
    const enhancedError = new Error(errorMessage);
    enhancedError.status = error.response?.status;
    enhancedError.data = error.response?.data;
    
    return Promise.reject(enhancedError);
  }
);

export default api; 