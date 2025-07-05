import React from 'react';
import { useApiGet } from '../../hooks/useApi';
import { bookService } from '../../services/bookService';
import { userService } from '../../services/userService';
import { loanService } from '../../services/loanService';
import { fineService } from '../../services/fineService';
import '../../styles/admin/Dashboard.css';

const Dashboard = () => {

  const { data: booksData, loading: booksLoading, error: booksError } = useApiGet(
    () => bookService.getFilteredBooks({}),
    []
  );

  // Obtener todos los usuarios
  const { data: allUsers, loading: usersLoading, error: usersError } = useApiGet(
    () => userService.getAllUsers(),
    []
  );

  // Para préstamos activos, usamos un rango amplio de fechas
  const currentYear = new Date().getFullYear();
  const { data: loansData, loading: loansLoading, error: loansError } = useApiGet(
    () => loanService.getActiveLoans({ 
      startdate: `${currentYear}-01-01`, 
      enddate: `${currentYear}-12-31` 
    }),
    []
  );

  const { data: finesData, loading: finesLoading, error: finesError } = useApiGet(
    () => fineService.getPendingFines(),
    []
  );



  const stats = [
    {
      title: 'Total de Libros',
      value: booksData?.books?.length || 0,
      color: '#10b981',
      icon: '📚',
      description: 'Libros en el catálogo'
    },
    {
      title: 'Usuarios Registrados',
      value: allUsers?.users?.length || 0,
      color: '#3b82f6',
      icon: '👥',
      description: 'Usuarios activos'
    },
    {
      title: 'Préstamos Activos',
      value: loansData?.loans?.length || 0,
      color: '#f59e0b',
      icon: '📖',
      description: 'Préstamos en curso'
    },
    {
      title: 'Multas Pendientes',
      value: finesData?.fines?.length || 0,
      color: '#ef4444',
      icon: '⚠️',
      description: 'Multas por cobrar'
    }
  ];

  const recentActivities = [
    {
      title: 'Libros Más Prestados',
      items: booksData?.books?.slice(0, 5) || [],
      type: 'books'
    },
    {
      title: 'Préstamos Recientes',
      items: loansData?.loans?.slice(0, 5) || [],
      type: 'loans'
    }
  ];

  const testBackendConnection = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/healthcheck`);
      const data = await response.json();
      setBackendStatus('connected');
    } catch (error) {
      setBackendStatus('disconnected');
    }
  };

  return (
    <div className="admin-dashboard-container">
      <div className="dashboard-header">
        <h1>Panel de Administración</h1>
        <p>Gestión integral de la biblioteca</p>
      </div>

      <div className="stats-grid">
        {stats.map((stat, index) => (
          <div key={index} className="stat-card" style={{ borderLeftColor: stat.color }}>
            <div className="stat-icon">{stat.icon}</div>
            <div className="stat-content">
              <h3>{stat.title}</h3>
              <p className="stat-value">
                {booksLoading || usersLoading || loansLoading || finesLoading ? 
                  'Cargando...' : stat.value}
              </p>
              <p className="stat-description">{stat.description}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="dashboard-sections">
        {recentActivities.map((section, index) => (
          <div key={index} className="section">
            <h2>{section.title}</h2>
            <div className="recent-items">
              {booksLoading || loansLoading ? (
                <p>Cargando...</p>
              ) : section.items.length > 0 ? (
                section.items.map((item, itemIndex) => (
                  <div key={itemIndex} className="item-card">
                    {section.type === 'books' ? (
                      <>
                        <h4>{item.titulo}</h4>
                        <p>Género: {item.genero}</p>
                        <p>Estado: {item.estado}</p>
                      </>
                    ) : (
                      <>
                        <h4>Préstamo #{item.idprestamo}</h4>
                        <p>Usuario: {item.idsocio}</p>
                        <p>Libro: {item.idlibro}</p>
                        <p>Fecha: {item.fechaprestamo}</p>
                      </>
                    )}
                  </div>
                ))
              ) : (
                <p className="no-items">
                  {section.type === 'books' ? 'No hay libros registrados' : 'No hay préstamos activos'}
                </p>
              )}
            </div>
          </div>
        ))}
      </div>



      <div className="quick-actions">
        <h2>Acciones Rápidas</h2>
        <div className="actions-grid">
          <a href="/admin/registrarlibro" className="action-card">
            <div className="action-icon">➕</div>
            <h3>Registrar Libro</h3>
            <p>Agregar nuevo libro al catálogo</p>
          </a>
          <a href="/admin/users" className="action-card">
            <div className="action-icon">👤</div>
            <h3>Gestionar Usuarios</h3>
            <p>Ver y administrar usuarios</p>
          </a>
          <a href="/admin/loans" className="action-card">
            <div className="action-icon">📋</div>
            <h3>Préstamos</h3>
            <p>Gestionar préstamos activos</p>
          </a>
          <a href="/admin/fines" className="action-card">
            <div className="action-icon">💰</div>
            <h3>Multas</h3>
            <p>Administrar multas pendientes</p>
          </a>
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 