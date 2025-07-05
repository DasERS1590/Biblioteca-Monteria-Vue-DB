import { useApiGet } from '../../hooks/useApi';
import { loanService } from '../../services/loanService';
import { fineService } from '../../services/fineService';
import { reservationService } from '../../services/reservationService';
import { getUser, getIdUser } from '../../auth';
import '../../styles/user/Dashboard.css';

const Dashboard = () => {
  const user = getUser();
  const userId = getIdUser();

  const { data: activeLoans, loading: loansLoading } = useApiGet(
    () => loanService.getUserActiveLoans(userId),
    [userId]
  );

  const { data: fines, loading: finesLoading } = useApiGet(
    () => fineService.getUserFines(userId),
    [userId]
  );

  const { data: reservations, loading: reservationsLoading } = useApiGet(
    () => reservationService.getUserReservations({ usuario_id: userId }),
    [userId]
  );

  const stats = [
    {
      title: 'Préstamos Activos',
      value: activeLoans?.loansactive?.length || 0,
      color: '#10b981',
      icon: '📚'
    },
    {
      title: 'Multas Pendientes',
      value: fines?.fines?.length || 0,
      color: '#ef4444',
      icon: '⚠️'
    },
    {
      title: 'Reservas Activas',
      value: reservations?.reservations?.length || 0,
      color: '#3b82f6',
      icon: '📅'
    }
  ];

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Bienvenido, {user?.nombre}</h1>
        <p>Panel de control de tu cuenta de biblioteca</p>
      </div>

      <div className="stats-grid">
        {stats.map((stat, index) => (
          <div key={index} className="stat-card" style={{ borderLeftColor: stat.color }}>
            <div className="stat-icon">{stat.icon}</div>
            <div className="stat-content">
              <h3>{stat.title}</h3>
              <p className="stat-value">{stat.value}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="dashboard-sections">
        <div className="section">
          <h2>Préstamos Recientes</h2>
          {loansLoading ? (
            <p>Cargando préstamos...</p>
          ) : (
            <div className="recent-items">
              {activeLoans?.loansactive?.slice(0, 3).map((loan, index) => (
                <div key={index} className="item-card">
                  <h4>{loan.titulo_libro}</h4>
                  <p>Fecha de devolución: {loan.fecha_devolucion}</p>
                  <span className={`status ${loan.estado}`}>{loan.estado}</span>
                </div>
              ))}
              {(!activeLoans?.loansactive || activeLoans.loansactive.length === 0) && (
                <p className="no-items">No tienes préstamos activos</p>
              )}
            </div>
          )}
        </div>

        <div className="section">
          <h2>Multas Pendientes</h2>
          {finesLoading ? (
            <p>Cargando multas...</p>
          ) : (
            <div className="recent-items">
              {fines?.fines?.slice(0, 3).map((fine, index) => (
                <div key={index} className="item-card fine">
                  <h4>Multa #{fine.id_multa}</h4>
                  <p>Monto: ${fine.monto}</p>
                  <p>Motivo: {fine.motivo}</p>
                </div>
              ))}
              {(!fines?.fines || fines.fines.length === 0) && (
                <p className="no-items">No tienes multas pendientes</p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 