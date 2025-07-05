import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getUser } from '../auth';

const DashboardRedirect = () => {
  const navigate = useNavigate();
  const user = getUser();

  useEffect(() => {
    if (user) {
      if (user.rol === 'administrador') {
        navigate('/admin/dashboard');
      } else {
        navigate('/user/dashboard');
      }
    } else {
      navigate('/login');
    }
  }, [user, navigate]);

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '100vh',
      fontSize: '18px',
      color: '#666'
    }}>
      Redirigiendo al dashboard...
    </div>
  );
};

export default DashboardRedirect;
