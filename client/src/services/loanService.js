import { api } from './api';

export const loanService = {
  // Obtener préstamos activos de un usuario
  getUserActiveLoans: async (usuarioId) => {
    try {
      const response = await api.get(`/loans?usuario_id=${usuarioId}`);
      return response.data;
    } catch (error) {
      // Si es un 404 (no hay préstamos activos), devolver array vacío
      if (error.response?.status === 404) {
        return { loansactive: [] };
      }
      throw error;
    }
  },

  // Obtener historial de préstamos completados de un usuario
  getUserCompletedLoanHistory: async (usuarioId) => {
    try {
      const response = await api.get(`/loans/completed?usuario_id=${usuarioId}`);
      return response.data;
    } catch (error) {
      // Si es un 404 (no hay préstamos completados), devolver array vacío
      if (error.response?.status === 404) {
        return { loanscomplete: [] };
      }
      throw error;
    }
  },

  // Obtener todos los préstamos (admin)
  getAllLoans: async (params) => {
    const response = await api.get('/loans', { params });
    return response.data;
  },

  // Crear un nuevo préstamo
  createLoan: async (data) => {
    const response = await api.post('/loans', data);
    return response.data;
  },

  // Extender un préstamo
  extendLoan: async (loanId, data) => {
    const response = await api.post(`/loans/extend/${loanId}`, data);
    return response.data;
  },

  // Devolver un préstamo
  returnLoan: async (loanId) => {
    const response = await api.post(`/loans/return/${loanId}`);
    return response.data;
  },

  // Obtener un préstamo específico
  getLoan: async (loanId) => {
    const response = await api.get(`/loans/${loanId}`);
    return response.data;
  },

  // Actualizar un préstamo
  updateLoan: async (loanId, data) => {
    const response = await api.put(`/loans/${loanId}`, data);
    return response.data;
  },

  // Eliminar un préstamo
  deleteLoan: async (loanId) => {
    const response = await api.delete(`/loans/${loanId}`);
    return response.data;
  },

  // Admin - Obtener préstamos activos por rango de fechas
  getActiveLoans: async (params = {}) => {
    try {
      const queryParams = new URLSearchParams();
      if (params.startdate) {
        queryParams.append('startdate', params.startdate);
      }
      if (params.enddate) {
        queryParams.append('enddate', params.enddate);
      }
      const response = await api.post(`/admin/loans?${queryParams.toString()}`, {});
      return response.data;
    } catch (error) {
      // Si es un error de permisos (403), devolver un error más descriptivo
      if (error.response?.status === 403) {
        throw new Error("No tienes permisos para acceder a esta funcionalidad. Solo los administradores pueden ver todos los préstamos.");
      }
      throw error;
    }
  },
  getUserLoanHistoryAdmin: async () => {
    const response = await api.get('/admin/loans/history');
    return response.data;
  },
}; 