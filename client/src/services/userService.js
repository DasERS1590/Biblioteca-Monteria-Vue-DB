import { api } from './api';

export const userService = {
  // Admin
  getUsersByType: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.tiposocio) {
      queryParams.append('tiposocio', params.tiposocio);
    }
    const response = await api.post(`/admin/users?${queryParams.toString()}`, {});
    return response.data;
  },
  
  // Obtener todos los usuarios
  getAllUsers: async () => {
    const response = await api.get('/admin/users/all');
    return response.data;
  },
}; 