import { api } from './api';

export const bookService = {
  // Usuario - Libros disponibles
  getAvailableBooks: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.genero) {queryParams.append('genero', params.genero);}
    if (params.autor) {queryParams.append('autor', params.autor);}
    if (params.titulo) {queryParams.append('titulo', params.titulo);}
    const response = await api.get(`/books?${queryParams.toString()}`);
    return response.data;
  },

  // Usuario - Libros para reservas (incluye no disponibles)
  getBooksForReservation: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.genero) {queryParams.append('genero', params.genero);}
    if (params.autor) {queryParams.append('autor', params.autor);}
    if (params.titulo) {queryParams.append('titulo', params.titulo);}
    const response = await api.get(`/books/reservation?${queryParams.toString()}`);
    return response.data;
  },

  // Obtener un libro específico
  getBook: async (bookId) => {
    const response = await api.get(`/books/${bookId}`);
    return response.data;
  },

  // Obtener información completa de un libro para editar
  getBookForEdit: async (bookId) => {
    const response = await api.get(`/admin/books/${bookId}/edit`);
    return response.data;
  },

  // Admin - Libros
  getBooks: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.estado) {queryParams.append('estado', params.estado);}
    if (params.editorial) {queryParams.append('editorial', params.editorial);}
    const response = await api.get(`/admin/books?${queryParams.toString()}`);
    return response.data;
  },

  getFilteredBooks: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.estado) {queryParams.append('estado', params.estado);}
    if (params.editorial) {queryParams.append('editorial', params.editorial);}
    const response = await api.get(`/admin/books?${queryParams.toString()}`);
    return response.data;
  },

  getUnavailableBooks: async () => {
    const response = await api.get('/admin/books/unavailable');
    return response.data;
  },
  
  getAvailableBooksAdmin: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.genero) {queryParams.append('genero', params.genero);}
    if (params.autor) {queryParams.append('autor', params.autor);}
    const response = await api.get(`/admin/books/available?${queryParams.toString()}`);
    return response.data;
  },

  getBooksByPublicationDate: async (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.fecha_inicio) {queryParams.append('fecha_inicio', params.fecha_inicio);}
    if (params.fecha_fin) {queryParams.append('fecha_fin', params.fecha_fin);}
    const response = await api.get(`/admin/books/published?${queryParams.toString()}`);
    return response.data;
  },

  createBook: async (bookData) => {
    const response = await api.post('/admin/books', bookData);
    return response.data;
  },
  
  updateBook: async (id, bookData) => {
    const response = await api.post(`/admin/books/${id}`, bookData);
    return response.data;
  },

  deleteBook: async (bookId) => {
    const response = await api.delete(`/admin/books/${bookId}`);
    return response.data;
  },
}; 