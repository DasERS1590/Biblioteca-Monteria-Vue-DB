import { useApiMutation, useForm } from '../../hooks/useApi';
import { editorialService } from '../../services/editorialService';
import "../../styles/admin/CreateEditorial.css"

function CreateEditorial() {
  const { formData, handleChange, resetForm } = useForm({
    nombre: '',
    direccion: '',
    paginaWeb: ''
  });

  const { execute: createEditorial, loading, error, success, reset } = useApiMutation(
    editorialService.createEditorial
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.nombre || !formData.direccion || !formData.paginaWeb) {
      return;
    }

    try {
      await createEditorial(formData);
      resetForm();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  return (
    <div className='form_container'>
      <h2>Guardar Nueva Editorial</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="text"
            name="nombre"
            value={formData.nombre}
            onChange={handleChange}
            placeholder="Nombre"
            required
          />
        </div>
        <div className="form-group">
          <input
            type="text"
            name="direccion"
            value={formData.direccion}
            onChange={handleChange}
            placeholder="Dirección"
            required
          />
        </div>
        <div className="form-group">
          <input
            type="text"
            name="paginaWeb"
            value={formData.paginaWeb}
            onChange={handleChange}
            placeholder="Página Web"
            required
          />
        </div>

        <button 
          type="submit" 
          className="my-button"
          disabled={loading}
        >
          {loading ? 'Creando...' : 'Crear Editorial'}
        </button>
      </form>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">Editorial guardada exitosamente</div>}
    </div>
  );
}

export default CreateEditorial;