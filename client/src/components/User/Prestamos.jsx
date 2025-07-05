import React, { useState } from "react";
import { useLoansGet, useApiMutation, useForm } from "../../hooks/useApi";
import { loanService } from "../../services/loanService";
import DataTable from "../common/DataTable";
import Modal from "../common/Modal";
import { getIdUser } from "../../auth";
import "../../styles/user/Prestamos.css";

const Prestamo = () => {
  const user_id = getIdUser();
  const [showExtendModal, setShowExtendModal] = useState(false);
  const [selectedLoan, setSelectedLoan] = useState(null);
  const [error, setError] = useState(null);
  
  const { data: loansData, loading, error: loansError, refetch } = useLoansGet(
    () => loanService.getUserActiveLoans(user_id),
    [user_id]
  );

  const { formData, handleChange, resetForm } = useForm({
    nuevafechadevolucion: ""
  });

  const { execute: extendLoan, loading: extendLoading, error: extendError, success: extendSuccess, reset: resetExtend } = useApiMutation(
    (data) => loanService.extendLoan(selectedLoan?.id_prestamo, data)
  );

  const { execute: returnLoan, loading: returnLoading, error: returnError, success: returnSuccess, reset: resetReturn } = useApiMutation(
    (loanId) => loanService.returnLoan(loanId)
  );

  const handleExtendLoan = async (e) => {
    e.preventDefault();
    if (!formData.nuevafechadevolucion) {
      return;
    }

    try {
      await extendLoan(formData);
      resetForm();
      setShowExtendModal(false);
      setSelectedLoan(null);
      refetch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const openExtendModal = (loan) => {
    setSelectedLoan(loan);
    setShowExtendModal(true);
    resetExtend();
  };

  const handleReturnLoan = async (loan) => {
    if (window.confirm(`¿Estás seguro de que quieres devolver el libro "${loan.titulo_libro}"?`)) {
      try {
        const result = await returnLoan(loan.id_prestamo);
  
        
        // Mostrar mensaje de éxito o multa generada
        if (result.return.fine_generated) {
          alert(`Libro devuelto exitosamente.\n\nSe generó una multa por retraso:\n- Días de retraso: ${result.return.days_late}\n- Monto: $${result.return.fine_amount}`);
        } else {
          alert("Libro devuelto exitosamente.");
        }
        
        refetch();
      } catch (error) {
        setError('Error al devolver el libro');
      }
    }
  };

  const closeExtendModal = () => {
    setShowExtendModal(false);
    setSelectedLoan(null);
    resetForm();
    resetExtend();
  };

  const loans = loansData?.loansactive || [];

  const columns = [
    { key: 'id_prestamo', label: 'ID Préstamo' },
    { key: 'fecha_prestamo', label: 'Fecha de Préstamo' },
    { key: 'fecha_devolucion', label: 'Fecha de Devolución' },
    { key: 'estado', label: 'Estado' },
    { key: 'titulo_libro', label: 'Título del Libro' },
    { 
      key: 'actions', 
      label: 'Acciones',
      render: (_, loan) => (
        <div className="loan-actions">
          <button 
            className="btn btn-primary btn-sm"
            onClick={() => openExtendModal(loan)}
            disabled={loan.estado !== 'activo'}
          >
            Extender
          </button>
          <button 
            className="btn btn-secondary btn-sm"
            onClick={() => handleReturnLoan(loan)}
            disabled={loan.estado !== 'activo'}
          >
            Devolver
          </button>
        </div>
      )
    }
  ];

  return (
    <div className="prestamo-dashboard">
      <h1>Préstamos Activos</h1>
      <p className="prestamo-description">Gestiona tus préstamos activos en la biblioteca.</p>
      
      <DataTable
        columns={columns}
        data={loans}
        loading={loading}
        error={loansError}
        emptyMessage="No hay préstamos activos para este usuario."
        className="prestamo-table"
      />

      {/* Modal para extender préstamo */}
      <Modal
        isOpen={showExtendModal}
        onClose={closeExtendModal}
        title="Extender Préstamo"
        className="extend-loan-modal"
      >
        {selectedLoan && (
          <div className="modal-content">
            <div className="loan-info">
              <h3>Información del Préstamo</h3>
              <div className="info-grid">
                <div className="info-item">
                  <label>Libro:</label>
                  <span>{selectedLoan.titulo_libro}</span>
                </div>
                <div className="info-item">
                  <label>Fecha de Préstamo:</label>
                  <span>{selectedLoan.fecha_prestamo}</span>
                </div>
                <div className="info-item">
                  <label>Fecha de Devolución Actual:</label>
                  <span>{selectedLoan.fecha_devolucion}</span>
                </div>
                <div className="info-item">
                  <label>Estado:</label>
                  <span className={`status ${selectedLoan.estado}`}>{selectedLoan.estado}</span>
                </div>
              </div>
            </div>

            <form onSubmit={handleExtendLoan} className="extend-form">
              <div className="form-group">
                <label>Nueva Fecha de Devolución:</label>
                <input
                  type="date"
                  name="nuevafechadevolucion"
                  value={formData.nuevafechadevolucion}
                  onChange={handleChange}
                  required
                  min={selectedLoan.fecha_devolucion}
                />
                <small className="form-help">
                  La nueva fecha debe ser posterior a la fecha de devolución actual.
                </small>
              </div>

              <div className="form-actions">
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={extendLoading || !formData.nuevafechadevolucion}
                >
                  {extendLoading ? 'Extendiendo...' : 'Extender Préstamo'}
                </button>
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={closeExtendModal}
                >
                  Cancelar
                </button>
              </div>

              {extendError && (
                <div className="error-message">{extendError}</div>
              )}
              {extendSuccess && (
                <div className="success-message">Préstamo extendido exitosamente</div>
              )}
            </form>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Prestamo;
