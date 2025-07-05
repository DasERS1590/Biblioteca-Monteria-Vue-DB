import { api } from './api';

export const editorialService = {
  getEditorials: async () => {
    const response = await api.get('/editoriales');
    return response.data;
  },
  createEditorial: async (editorialData) => {
    const response = await api.post('/editoriales', editorialData);
    return response.data;
  },
}; 