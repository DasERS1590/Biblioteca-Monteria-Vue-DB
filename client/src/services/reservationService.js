import { api } from './api';

export const reservationService = {
  // Obtener reservas de un usuario
  getUserReservations: async (params) => {
    const queryParams = new URLSearchParams();
    if (params.usuario_id) {
      queryParams.append('usuario_id', params.usuario_id);
    }
    try {
      const response = await api.get(`/reservations?${queryParams.toString()}`);
      return response.data;
    } catch (error) {
      // Si es un 404 (no hay reservas), devolver array vacío
      if (error.response?.status === 404) {
        return { reservations: [] };
      }
      throw error;
    }
  },

  // Obtener todas las reservas (admin)
  getAllReservations: async (params) => {
    const response = await api.get('/admin/reservations', { params });
    return response.data;
  },

  // Crear una nueva reserva
  createReservation: async (data) => {
    const response = await api.post('/reservation', data);
    return response.data;
  },

  // Cancelar una reserva
  cancelReservation: async (reservationId) => {
    const response = await api.delete(`/reservations/${reservationId}`);
    return response.data;
  },

  // Obtener una reserva específica
  getReservation: async (reservationId) => {
    const response = await api.get(`/reservations/${reservationId}`);
    return response.data;
  },

  // Actualizar una reserva
  updateReservation: async (reservationId, data) => {
    const response = await api.put(`/reservations/${reservationId}`, data);
    return response.data;
  },

  // Eliminar una reserva
  deleteReservation: async (reservationId) => {
    const response = await api.delete(`/reservations/${reservationId}`);
    return response.data;
  },

  // Admin - Obtener reservas activas
  getActiveReservations: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.idsocio) {
      queryParams.append('idsocio', params.idsocio);
    }
    if (params.idlibro) {
      queryParams.append('idlibro', params.idlibro);
    }
    if (params.fechareserva) {
      queryParams.append('fechareserva', params.fechareserva);
    }
    if (params.nombre_socio) {
      queryParams.append('nombre_socio', params.nombre_socio);
    }
    const response = await api.get(`/admin/reservations?${queryParams.toString()}`);
    return response.data;
  }
}; 