import React, { useState } from "react";
import { useApiGet, useApiMutation, useForm, useFilters } from "../../hooks/useApi";
import { reservationService } from "../../services/reservationService";
import { bookService } from "../../services/bookService";
import DataTable from "../common/DataTable";
import FilterForm from "../common/FilterForm";
import Modal from "../common/Modal";
import ConfirmDialog from "../common/ConfirmDialog";
import "../../styles/user/Reservas.css";

function Reserva() {
  const user = JSON.parse(localStorage.getItem("user"));
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showCancelDialog, setShowCancelDialog] = useState(false);
  const [selectedReservation, setSelectedReservation] = useState(null);
  const [books, setBooks] = useState([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [error, setError] = useState(null);

  const { filters, updateFilter, clearFilters } = useFilters({
    titulo: "",
    genero: "",
    autor: ""
  });

  const { formData, handleChange, resetForm } = useForm({
    idlibro: "",
    fechareserva: ""
  });

  const { data: reservationsResponse, loading, error: apiError, refetch } = useApiGet(
    () => reservationService.getUserReservations({ usuario_id: user?.id }),
    [user?.id]
  );

  // Extraer las reservas del response, manejando el caso cuando no hay reservas
  const reservations = reservationsResponse?.reservations || [];

  const { execute: createReservation, loading: createLoading, error: createError, success: createSuccess, reset: resetCreate } = useApiMutation(
    (data) => reservationService.createReservation({
      ...data,
      idsocio: user?.id
    })
  );

  const { execute: cancelReservation, loading: cancelLoading, error: cancelError, success: cancelSuccess, reset: resetCancel } = useApiMutation(
    (reservationId) => reservationService.cancelReservation(reservationId)
  );

  const handleSearchBooks = async () => {
    setSearchLoading(true);
    try {
      const response = await bookService.getBooksForReservation(filters);
      // El servicio ya retorna response.data, que contiene { books: [...] }
      const booksArray = response.books || [];
      setBooks(booksArray);
    } catch (err) {
      setError('Error al buscar libros');
    } finally {
      setSearchLoading(false);
    }
  };

  const handleCreateReservation = async (e) => {
    e.preventDefault();
    if (!formData.idlibro || !formData.fechareserva) {
      return;
    }

    try {
      await createReservation(formData);
      resetForm();
      setShowCreateModal(false);
      refetch();
      // Mostrar mensaje de éxito
      alert("¡Reserva creada exitosamente!");
    } catch (error) {
      setError('Error al crear la reserva');
    }
  };

  const handleCancelReservation = async () => {
    if (!selectedReservation) {return;}

    try {
      await cancelReservation(selectedReservation.id_reserva);
      setShowCancelDialog(false);
      setSelectedReservation(null);
      refetch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const openCancelDialog = (reservation) => {
    setSelectedReservation(reservation);
    setShowCancelDialog(true);
    resetCancel();
  };

  const closeCreateModal = () => {
    setShowCreateModal(false);
    resetForm();
    resetCreate();
  };

  const filterConfig = [
    {
      name: "titulo",
      placeholder: "Título del libro",
      value: filters.titulo
    },
    {
      name: "genero",
      placeholder: "Género",
      value: filters.genero
    },
    {
      name: "autor",
      placeholder: "Autor",
      value: filters.autor
    }
  ];

  const columns = [
    { key: 'id_reserva', label: 'ID Reserva' },
    { key: 'id_libro', label: 'ID Libro' },
    { key: 'fecha_reserva', label: 'Fecha Reserva' },
    { key: 'estado', label: 'Estado' },
    { 
      key: 'actions', 
      label: 'Acciones',
      render: (_, reservation) => (
        <button 
          className="btn btn-danger" 
          onClick={() => openCancelDialog(reservation)}
          disabled={reservation.estado !== 'activa'}
        >
          {reservation.estado === 'activa' ? 'Cancelar' : 'Cancelada'}
        </button>
      )
    }
  ];

  return (
    <div className="reservas-dashboard">
      <div className="reservas-header">
        <h1>Reservas</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowCreateModal(true)}
        >
          Crear Nueva Reserva
        </button>
      </div>
      
      <p className="reservas-description">Gestiona tus reservas de libros en la biblioteca.</p>
      
      <DataTable
        columns={columns}
        data={reservations || []}
        loading={loading}
        error={apiError || error}
        emptyMessage="No hay reservas para este usuario."
        className="reservas-table"
      />
      


      {/* Modal para crear reserva */}
      <Modal
        isOpen={showCreateModal}
        onClose={closeCreateModal}
        title="Crear Nueva Reserva"
        className="reserva-modal"
      >
        <div className="modal-content">
          <div className="search-section">
            <h3>Buscar Libros para Reservar</h3>
            <FilterForm
              filters={filterConfig}
              onFilterChange={updateFilter}
              onSearch={handleSearchBooks}
              onClear={clearFilters}
              loading={searchLoading}
            />
            
            {books.length > 0 && (
              <div className="books-list">
                <h4>Libros Encontrados:</h4>
                <div className="books-grid">
                  {books.map((book) => (
                    <div 
                      key={book.id_libro} 
                      className={`book-item ${formData.idlibro === book.id_libro ? 'selected' : ''}`}
                      onClick={() => handleChange({ target: { name: 'idlibro', value: book.id_libro } })}
                    >
                      <h5>{book.titulo}</h5>
                      <p>Género: {book.genero}</p>
                      <p>Autor: {book.autor}</p>
                      <p className={`estado ${book.estado}`}>Estado: {book.estado}</p>
                      <small className="book-hint">
                        {book.estado === 'prestado' ? '📚 Prestado - Puedes reservarlo para cuando se devuelva' : 
                         '🔒 Reservado - Puedes hacer cola de espera'}
                      </small>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          <form onSubmit={handleCreateReservation} className="reservation-form">
            {formData.idlibro && (
              <div className="form-group">
                <label>Libro Seleccionado:</label>
                <div className="selected-book-info">
                  {books.find(book => book.id_libro === formData.idlibro)?.titulo || 'Libro no encontrado'}
                  <span className="book-status">
                    (Estado: {books.find(book => book.id_libro === formData.idlibro)?.estado || 'N/A'})
                  </span>
                </div>
              </div>
            )}
            
            <div className="form-group">
              <label>Fecha de Reserva:</label>
              <input
                type="date"
                name="fechareserva"
                value={formData.fechareserva}
                onChange={handleChange}
                required
                min={new Date().toISOString().split('T')[0]}
              />
            </div>

            <div className="form-actions">
              <button 
                type="submit" 
                className="btn btn-primary"
                disabled={createLoading || !formData.idlibro || !formData.fechareserva}
              >
                {createLoading ? 'Creando...' : 'Crear Reserva'}
              </button>
              <button 
                type="button" 
                className="btn btn-secondary"
                onClick={closeCreateModal}
              >
                Cancelar
              </button>
            </div>

            {createError && (
              <div className="error-message">{createError}</div>
            )}
            {createSuccess && (
              <div className="success-message">Reserva creada exitosamente</div>
            )}
          </form>
        </div>
      </Modal>

      {/* Diálogo de confirmación para cancelar */}
      <ConfirmDialog
        isOpen={showCancelDialog}
        title="Cancelar Reserva"
        message={`¿Estás seguro de que quieres cancelar la reserva #${selectedReservation?.id_reserva}?`}
        onConfirm={handleCancelReservation}
        onCancel={() => {
          setShowCancelDialog(false);
          setSelectedReservation(null);
        }}
        confirmText="Cancelar Reserva"
        cancelText="Mantener"
        type="warning"
      />

      {cancelError && (
        <div className="error-message">{cancelError}</div>
      )}
      {cancelSuccess && (
        <div className="success-message">Reserva cancelada exitosamente</div>
      )}
    </div>
  );
}

export default Reserva;
