import React, { useEffect } from 'react';
import { useApiGet, useApiMutation, useForm } from '../../hooks/useApi';
import { bookService } from '../../services/bookService';
import { editorialService } from '../../services/editorialService';
import { authorService } from '../../services/authorService';
import '../../styles/admin/CreateBook.css';

function CreateBook() {
  const { formData, handleChange, resetForm } = useForm({
    titulo: '',
    genero: '',
    fechapublicacion: '',
    ideditorial: '',
    idautores: '',
  });

  const { data: editorialsData, loading: editorialsLoading, error: editorialsError } = useApiGet(
    editorialService.getEditorials,
    []
  );

  const { data: authorsData, loading: authorsLoading, error: authorsError } = useApiGet(
    authorService.getAuthors,
    []
  );

  const { execute: createBook, loading: createLoading, error: createError, success: createSuccess, reset: resetCreate } = useApiMutation(
    (data) => bookService.createBook({
      ...data,
      ideditorial: Number(data.ideditorial),
      idautores: [Number(data.idautores)],
    })
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const { titulo, genero, fechapublicacion, ideditorial, idautores } = formData;

    if (!titulo || !genero || !fechapublicacion || !ideditorial || !idautores) {
      return;
    }

    try {
      await createBook(formData);
      resetForm();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const editorials = editorialsData?.editorials || [];
  const authors = authorsData?.authors || authorsData?.author || [];

  return (
    <div className="form-container">
      <h2>Crear Nuevo Libro</h2>
      
      {(editorialsError || authorsError) && (
        <div className="error-message">
          Error al cargar los datos necesarios
        </div>
      )}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="text"
            name="titulo"
            value={formData.titulo}
            onChange={handleChange}
            required
            placeholder="Título"
          />
        </div>

        <div className="form-group">
          <input
            type="text"
            name="genero"
            value={formData.genero}
            onChange={handleChange}
            required
            placeholder="Género"
          />
        </div>

        <div className="form-group">
          <input
            type="date"
            name="fechapublicacion"
            value={formData.fechapublicacion}
            onChange={handleChange}
            required
            placeholder="Fecha de Publicación"
          />
        </div>

        <div className="form-group">
          <select
            name="ideditorial"
            value={formData.ideditorial}
            onChange={handleChange}
            required
            disabled={editorialsLoading}
          >
            <option value="">Seleccione una editorial</option>
            {editorials.length > 0 ? (
              editorials.map(editorial => (
                <option key={editorial.ideditorial} value={editorial.ideditorial}>
                  {editorial.nombre}
                </option>
              ))
            ) : (
              <option disabled>
                {editorialsLoading ? 'Cargando editoriales...' : 'No hay editoriales disponibles'}
              </option>
            )}
          </select>
        </div>

        <div className="form-group">
          <select
            name="idautores"
            value={formData.idautores}
            onChange={handleChange}
            required
            disabled={authorsLoading}
          >
            <option value="">Seleccione un autor</option>
            {authors.length > 0 ? (
              authors.map(author => (
                <option key={author.idautor} value={author.idautor}>
                  {author.nombre}
                </option>
              ))
            ) : (
              <option disabled>
                {authorsLoading ? 'Cargando autores...' : 'No hay autores disponibles'}
              </option>
            )}
          </select>
        </div>

        <button 
          type="submit" 
          className="submit-btn"
          disabled={createLoading || editorialsLoading || authorsLoading}
        >
          {createLoading ? 'Creando...' : 'Crear Libro'}
        </button>
      </form>

      {createError && <div className="error-message">{createError}</div>}
      {createSuccess && <div className="success-message">Libro creado exitosamente</div>}
    </div>
  );
}

export default CreateBook;