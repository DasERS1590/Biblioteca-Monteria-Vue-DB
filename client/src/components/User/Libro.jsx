import { useState } from "react";
import { useApiMutation, useForm, useFilters } from "../../hooks/useApi";
import { bookService } from "../../services/bookService";
import { loanService } from "../../services/loanService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import "../../styles/user/Libro.css";

const Libro = () => {
  const [books, setBooks] = useState([]);
  const [selectedBook, setSelectedBook] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState(null);

  const { filters, updateFilter, clearFilters } = useFilters({
    genero: "",
    autor: "",
    titulo: ""
  });

  const { formData, handleChange, resetForm } = useForm({
    fechaPrestamo: "",
    fechaDevolucion: ""
  });

  const { execute: createLoan, loading: loanLoading, error: loanError, success: loanSuccess, reset: resetLoan } = useApiMutation(loanService.createLoan);

  const handleSearch = async () => {
    setSearchLoading(true);
    setSearchError(null);
    
    // Verificar si el usuario está autenticado
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (!token) {
      setSearchError("Debe iniciar sesión para buscar libros");
      setSearchLoading(false);
      return;
    }
    
    if (!user) {
      setSearchError("Información de usuario no encontrada. Por favor, inicie sesión nuevamente.");
      setSearchLoading(false);
      return;
    }
    

    
    try {
      const response = await bookService.getAvailableBooks(filters);
      
      // El servicio ya retorna response.data, que contiene { books: [...] }
      const booksArray = response.books || [];
      
      setBooks(booksArray);
      
      if (booksArray.length === 0) {

      }
    } catch (err) {
      setError('Error al buscar libros');
    } finally {
      setSearchLoading(false);
  
    }
  };

  const handleBookSelect = (book) => {
    setSelectedBook(book);
    setShowModal(true);
    resetLoan();
  };

  const handleLoanSubmit = async (e) => {
    e.preventDefault();
    const user = JSON.parse(localStorage.getItem("user"));
    
    if (!user) {
      setSearchError("Usuario no encontrado. Asegúrese de estar registrado.");
      return;
    }

    const loanData = {
      usuario_id: user.id,
      libro_id: selectedBook.id_libro,
      fecha_prestamo: formData.fechaPrestamo,
      fecha_devolucion: formData.fechaDevolucion,
    };

    try {
      await createLoan(loanData);
      resetForm();
      setShowModal(false);
      setSelectedBook(null);
      // Recargar la búsqueda para actualizar el estado de los libros
      handleSearch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const closeModal = () => {
    setShowModal(false);
    setSelectedBook(null);
    resetForm();
    resetLoan();
  };

  const filterConfig = [
    {
      name: "genero",
      placeholder: "Género (Ej: Ciencia Ficción)",
      value: filters.genero
    },
    {
      name: "autor", 
      placeholder: "Autor (Ej: Isaac Asimov)",
      value: filters.autor
    },
    {
      name: "titulo",
      placeholder: "Título (Ej: Fundación)", 
      value: filters.titulo
    }
  ];

  const columns = [
    { key: 'id_libro', label: 'ID' },
    { key: 'titulo', label: 'Título' },
    { key: 'genero', label: 'Género' },
    { key: 'estado', label: 'Estado' },
    { key: 'autor', label: 'Autor' },
    { 
      key: 'actions', 
      label: 'Acción',
      render: (_, book) => (
        <button 
          className="btn" 
          onClick={() => handleBookSelect(book)}
          disabled={book.estado !== 'disponible'}
        >
          {book.estado === 'disponible' ? 'Seleccionar' : 'No disponible'}
        </button>
      )
    }
  ];

  return (
    <div className="admin-dashboard">
      <h1>Libros</h1>
      <p className="description">Busca y solicita libros disponibles en la biblioteca.</p>
      
      <FilterForm
        filters={filterConfig}
        onFilterChange={updateFilter}
        onSearch={handleSearch}
        onClear={clearFilters}
        loading={searchLoading}
      />

      {searchError && <p className="error">{searchError}</p>}
      
      <DataTable
        columns={columns}
        data={books}
        loading={searchLoading}
        error={searchError}
        emptyMessage="No se encontraron libros."
        className="libro-table"
      />

      {showModal && selectedBook && (
        <div className="libro-modal">
          <div className="libro-modal-content">
            {/* Header con título y botón de cerrar */}
            <div className="libro-modal-header">
              <h2>Formulario de Préstamo</h2>
              <button 
                className="libro-modal-close"
                onClick={closeModal}
                type="button"
              >
                ×
              </button>
            </div>
            
            {/* Cuerpo del modal */}
            <div className="libro-modal-body">
              {/* Información del libro seleccionado */}
              <div className="libro-selected-book">
                <p><strong>Libro Seleccionado:</strong> {selectedBook.titulo}</p>
              </div>
              
              {/* Formulario */}
              <form onSubmit={handleLoanSubmit} className="libro-modal-form">
                <div className="libro-form-group">
                  <label htmlFor="fechaPrestamo">Fecha de Préstamo</label>
                  <input
                    type="date"
                    id="fechaPrestamo"
                    name="fechaPrestamo"
                    value={formData.fechaPrestamo}
                    onChange={handleChange}
                    required
                    className="libro-input"
                  />
                </div>
                
                <div className="libro-form-group">
                  <label htmlFor="fechaDevolucion">Fecha de Devolución</label>
                  <input
                    type="date"
                    id="fechaDevolucion"
                    name="fechaDevolucion"
                    value={formData.fechaDevolucion}
                    onChange={handleChange}
                    required
                    className="libro-input"
                  />
                </div>
                
                <button 
                  type="submit" 
                  className="libro-btn"
                  disabled={loanLoading}
                >
                  {loanLoading ? 'Procesando...' : 'Realizar Préstamo'}
                </button>
                
                {loanError && (
                  <p className="libro-error">{loanError}</p>
                )}
                {loanSuccess && (
                  <p className="success-message">Préstamo realizado con éxito</p>
                )}
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Libro;