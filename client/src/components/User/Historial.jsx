import { useLoansGet } from '../../hooks/useApi';
import { loanService } from '../../services/loanService';
import DataTable from '../common/DataTable';
import { getIdUser } from '../../auth';
import "../../styles/user/Historial.css";

const Historial = () => {
  const user = getIdUser();
  
  const { data: loansData, loading, error } = useLoansGet(
    () => loanService.getUserCompletedLoanHistory(user),
    [user]
  );

  const loans = loansData?.loanscomplete || [];

  const columns = [
    { key: 'id_prestamo', label: 'ID Préstamo' },
    { key: 'fecha_prestamo', label: 'Fecha Préstamo' },
    { key: 'fecha_devolucion', label: 'Fecha Devolución' },
    { key: 'estado', label: 'Estado' },
    { key: 'titulo_libro', label: 'Título del Libro' }
  ];

  return (
    <div className="historial-container">
      <h1>Historial de Préstamos Completados</h1>
      
      <DataTable
        columns={columns}
        data={loans}
        loading={loading}
        error={error}
        emptyMessage="No hay préstamos completados para este usuario"
        className="historial-table"
      />
    </div>
  );
};

export default Historial;
