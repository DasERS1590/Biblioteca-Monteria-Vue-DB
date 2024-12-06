import React, { useState, useEffect } from 'react';
import axios from 'axios';

function CreateBook() {
  const [formData, setFormData] = useState({
    titulo: '',
    genero: '',
    fechaPublicacion: '',
    editorialID: '',
    autores: [], // Lista de IDs de autores
  });

  const [editorials, setEditorials] = useState([]);
  const [authors, setAuthors] = useState([]);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Fetch editoriales y autores
  useEffect(() => {
    // Obtener editoriales
    axios.get('http://localhost:4000/api/editoriales')
      .then(response => setEditorials(response.data)) // Asegúrate que la respuesta sea un array de editoriales
      .catch(error => setError('Error al cargar editoriales'));

    // Obtener autores
    axios.get('http://localhost:4000/api/autores')
      .then(response => setAuthors(response.data)) // Asegúrate que la respuesta sea un array de autores
      .catch(error => setError('Error al cargar autores'));
  }, []);

  // Maneja los cambios de los campos del formulario
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  // Maneja la selección de autores
  const handleSelectAutores = (e) => {
    const { options } = e.target;
    const selectedAutores = Array.from(options).filter(option => option.selected).map(option => option.value);
    setFormData((prevData) => ({
      ...prevData,
      autores: selectedAutores,
    }));
  };

  // Maneja el envío del formulario
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    // Validación de campos
    if (
      !formData.titulo ||
      !formData.genero ||
      !formData.fechaPublicacion ||
      !formData.editorialID ||
      formData.autores.length === 0
    ) {
      setError('Todos los campos son obligatorios');
      return;
    }

    try {
      // Realizar solicitud POST al backend
      const response = await axios.post('http://localhost:4000/api/libros', formData, {
        headers: {
          'Content-Type': 'application/json',
        },
      });

      // Respuesta exitosa
      setSuccess('Libro creado exitosamente');
    } catch (error) {
      // Manejo de errores
      if (error.response) {
        setError(error.response.data.message || 'Error al crear el libro');
      } else {
        setError('Error en la conexión con el servidor');
      }
    }
  };

  const formStyle = {
    display: 'flex',
    flexDirection: 'column',
    maxWidth: '400px',
    margin: '0 auto',
    padding: '20px',
    borderRadius: '8px',
    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
  };

  const inputStyle = {
    padding: '8px',
    margin: '10px 0',
    border: '1px solid #ddd',
    borderRadius: '4px',
  };

  const selectStyle = {
    padding: '8px',
    margin: '10px 0',
    border: '1px solid #ddd',
    borderRadius: '4px',
  };

  const buttonStyle = {
    padding: '10px 20px',
    backgroundColor: '#007BFF',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    marginTop: '10px',
  };

  const errorStyle = {
    color: 'red',
    marginTop: '10px',
  };

  const successStyle = {
    color: 'green',
    marginTop: '10px',
  };

  return (
    <div>
      <h2>Crear Nuevo Libro</h2>
      <form onSubmit={handleSubmit} style={formStyle}>
        <div>
          <label>Título:</label>
          <input
            type="text"
            name="titulo"
            value={formData.titulo}
            onChange={handleChange}
            style={inputStyle}
          />
        </div>
        <div>
          <label>Género:</label>
          <input
            type="text"
            name="genero"
            value={formData.genero}
            onChange={handleChange}
            style={inputStyle}
          />
        </div>
        <div>
          <label>Fecha de Publicación:</label>
          <input
            type="date"
            name="fechaPublicacion"
            value={formData.fechaPublicacion}
            onChange={handleChange}
            style={inputStyle}
          />
        </div>
        <div>
          <label>Editorial:</label>
          <select
            name="editorialID"
            value={formData.editorialID}
            onChange={handleChange}
            style={selectStyle}
          >
            <option value="">Seleccione una editorial</option>
            {editorials && editorials.length > 0 ? (
              editorials.map(editorial => (
                <option key={editorial.id} value={editorial.id}>
                  {editorial.nombre}
                </option>
              ))
            ) : (
              <option value="">No hay editoriales disponibles</option>
            )}
          </select>
        </div>
        <div>
          <label>Autores:</label>
          <select
            name="autores"
            multiple
            value={formData.autores}
            onChange={handleSelectAutores}
            style={selectStyle}
          >
            {authors && authors.length > 0 ? (
              authors.map(author => (
                <option key={author.id} value={author.id}>
                  {author.nombre}
                </option>
              ))
            ) : (
              <option value="">No hay autores disponibles</option>
            )}
          </select>
        </div>
        <button type="submit" style={buttonStyle}>Crear Libro</button>
      </form>

      {error && <div style={errorStyle}>{error}</div>}
      {success && <div style={successStyle}>{success}</div>}
    </div>
  );
}

export default CreateBook;
