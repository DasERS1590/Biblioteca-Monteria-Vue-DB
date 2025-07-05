import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import '../../styles/common/Breadcrumbs.css';

const Breadcrumbs = () => {
  const location = useLocation();
  const pathnames = location.pathname.split('/').filter((x) => x);

  const getBreadcrumbName = (path) => {
    const breadcrumbMap = {
      'admin': 'Administración',
      'user': 'Usuario',
      'dashboard': 'Dashboard',
      'books': 'Libros',
      'users': 'Usuarios',
      'loans': 'Préstamos',
      'fines': 'Multas',
      'reservations': 'Reservas',
      'registrarlibro': 'Registrar Libro',
      'registaredi': 'Registrar Editorial',
      'registrar-autor': 'Registrar Autor',
      'libro': 'Libros',
      'prestamo': 'Préstamos',
      'reserva': 'Reservas',
      'multa': 'Multas',
      'historial': 'Historial'
    };
    return breadcrumbMap[path] || path;
  };

  return (
    <nav className="breadcrumbs">
      <Link to="/" className="breadcrumb-item">
        Inicio
      </Link>
      {pathnames.map((name, index) => {
        const routeTo = `/${pathnames.slice(0, index + 1).join('/')}`;
        const isLast = index === pathnames.length - 1;
        
        return (
          <React.Fragment key={name}>
            <span className="breadcrumb-separator">/</span>
            {isLast ? (
              <span className="breadcrumb-item active">
                {getBreadcrumbName(name)}
              </span>
            ) : (
              <Link to={routeTo} className="breadcrumb-item">
                {getBreadcrumbName(name)}
              </Link>
            )}
          </React.Fragment>
        );
      })}
    </nav>
  );
};

export default Breadcrumbs; 