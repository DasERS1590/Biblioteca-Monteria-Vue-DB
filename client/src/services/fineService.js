import { api } from './api';

export const fineService = {
  // Obtener multas pendientes de un usuario
  getUserFines: async (usuarioId) => {
    try {
      const response = await api.get(`/fines?usuario_id=${usuarioId}`);
      return response.data;
    } catch (error) {
      // Si es un 404 (no hay multas), devolver array vacío
      if (error.response?.status === 404) {
        return { fines: [] };
      }
      throw error;
    }
  },

  // Obtener todas las multas (admin)
  getAllFines: async (params) => {
    const response = await api.get('/fines', { params });
    return response.data;
  },

  // Crear una nueva multa
  createFine: async (data) => {
    const response = await api.post('/fines', data);
    return response.data;
  },

  // Pagar una multa
  payFine: async (fineId) => {
    const response = await api.put(`/fines/${fineId}/pay`);
    return response.data;
  },

  // Obtener una multa específica
  getFine: async (fineId) => {
    const response = await api.get(`/fines/${fineId}`);
    return response.data;
  },

  // Actualizar una multa
  updateFine: async (fineId, data) => {
    const response = await api.put(`/fines/${fineId}`, data);
    return response.data;
  },

  // Eliminar una multa
  deleteFine: async (fineId) => {
    const response = await api.delete(`/fines/${fineId}`);
    return response.data;
  },

  // Admin
  getPendingFines: async () => {
    const response = await api.get('/admin/fines/to');
    return response.data;
  },
  getUserFinesAdmin: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.usuario_id) {queryParams.append('idsocio', params.usuario_id);}
    const response = await api.get(`/admin/fines?${queryParams.toString()}`);
    return response.data;
  },
  
  // Buscar multas por nombre o email del usuario
  searchFinesByUser: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.busqueda) {queryParams.append('busqueda', params.busqueda);}
    const response = await api.get(`/admin/fines/search?${queryParams.toString()}`);
    return response.data;
  },
}; 