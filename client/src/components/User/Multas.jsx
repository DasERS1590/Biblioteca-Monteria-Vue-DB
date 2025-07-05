import React, { useState } from "react";
import { useLoansGet, useApiMutation } from "../../hooks/useApi";
import { fineService } from "../../services/fineService";
import DataTable from "../common/DataTable";
import ConfirmDialog from "../common/ConfirmDialog";
import { getIdUser } from "../../auth";
import "../../styles/user/Multas.css";

const Multa = () => {
  const user_id = getIdUser();
  const [showPayDialog, setShowPayDialog] = useState(false);
  const [selectedFine, setSelectedFine] = useState(null);
  
  const { data: finesData, loading, error, refetch } = useLoansGet(
    () => fineService.getUserFines(user_id),
    [user_id]
  );

  const { execute: payFine, loading: payLoading, error: payError, success: paySuccess, reset: resetPay } = useApiMutation(
    (fineId) => fineService.payFine(fineId)
  );

  const handlePayFine = async () => {
    if (!selectedFine) {return;}

    try {
      await payFine(selectedFine.idmulta);
      setShowPayDialog(false);
      setSelectedFine(null);
      refetch();
    } catch (error) {
      // El error se maneja en el hook
    }
  };

  const openPayDialog = (fine) => {
    setSelectedFine(fine);
    setShowPayDialog(true);
    resetPay();
  };

  const closePayDialog = () => {
    setShowPayDialog(false);
    setSelectedFine(null);
    resetPay();
  };

  const fines = finesData?.fines || [];

  const columns = [
    { key: 'idmulta', label: 'ID Multa' },
    { key: 'saldopagar', label: 'Monto', render: (value) => `$${value}` },
    { key: 'fechamulta', label: 'Fecha de Multa' },
    { key: 'estado', label: 'Estado', render: (value) => (
      <span className={`status ${value}`}>{value}</span>
    )},
    { 
      key: 'actions', 
      label: 'Acciones',
      render: (_, fine) => (
        <button 
          className="btn btn-success btn-sm"
          onClick={() => openPayDialog(fine)}
          disabled={fine.estado === 'pagada'}
        >
          {fine.estado === 'pagada' ? 'Pagada' : 'Pagar'}
        </button>
      )
    }
  ];

  return (
    <div className="multa-dashboard">
      <h1>Multas</h1>
      <p className="multa-description">Consulta y paga tus multas pendientes en la biblioteca.</p>
      
      <DataTable
        columns={columns}
        data={fines}
        loading={loading}
        error={error}
        emptyMessage="No hay multas para este usuario."
        className="multa-table"
      />

      {/* Diálogo de confirmación para pagar */}
      <ConfirmDialog
        isOpen={showPayDialog}
        title="Pagar Multa"
        message={
          selectedFine ? (
            <div className="fine-details">
              <p><strong>ID Multa:</strong> {selectedFine.idmulta}</p>
              <p><strong>Monto:</strong> ${selectedFine.saldopagar}</p>
              <p><strong>Fecha:</strong> {selectedFine.fechamulta}</p>
              <p><strong>Estado:</strong> {selectedFine.estado}</p>
              <div className="payment-warning">
                <p>¿Estás seguro de que quieres pagar esta multa?</p>
                <p>Esta acción no se puede deshacer.</p>
              </div>
            </div>
          ) : ""
        }
        onConfirm={handlePayFine}
        onCancel={closePayDialog}
        confirmText={payLoading ? "Pagando..." : "Confirmar Pago"}
        cancelText="Cancelar"
        type="warning"
        disabled={payLoading}
      />

      {payError && (
        <div className="error-message">{payError}</div>
      )}
      {paySuccess && (
        <div className="success-message">Multa pagada exitosamente</div>
      )}
    </div>
  );
};

export default Multa;
