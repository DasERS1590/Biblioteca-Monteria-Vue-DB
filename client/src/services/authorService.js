import { api } from './api';

export const authorService = {
  getAuthors: async () => {
    const response = await api.get('/autores');
    return response.data;
  },
  createAuthor: async (authorData) => {
    const response = await api.post('/admin/autores', authorData);
    return response.data;
  },
}; 