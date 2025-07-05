import { useApiMutation, useForm } from "../../hooks/useApi";
import { authorService } from "../../services/authorService";
import "../../styles/admin/Author.css";

function CreateAuthor() {
  const { formData, handleChange, resetForm } = useForm({
    name: '',
    nationality: ''
  });

  const { execute: createAuthor, loading, error, success, reset } = useApiMutation(
    authorService.createAuthor
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name || !formData.nationality) {
      return;
    }

    try {
      await createAuthor(formData);
      resetForm();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  return (
    <div className='author-form-container'>
      <h2>Guardar Nuevo Autor</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Nombre"
            required
          />
        </div>
        <div className="form-group">
          <input
            type="text"
            name="nationality"
            value={formData.nationality}
            onChange={handleChange}
            placeholder="Nacionalidad"
            required
          />
        </div>

        <button 
          type="submit" 
          className="my-button"
          disabled={loading}
        >
          {loading ? 'Creando...' : 'Crear Autor'}
        </button>
      </form>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">Autor guardado exitosamente</div>}
    </div>
  );
}

export default CreateAuthor; 