import React, { useEffect, useState } from 'react';

const Historial = () => {
  const [loans, setLoans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Recuperar el usuario desde localStorage
    const user = JSON.parse(localStorage.getItem("user"));
    
    if (user) {
      console.log("Usuario logueado:", user);
      console.log("Rol del usuario:", user.rol);

      // Realizar la solicitud GET para obtener el historial de préstamos completados
      fetch(`http://localhost:4000/api/loans/completed?usuario_id=${user.id}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("No se pudo obtener el historial de préstamos");
          }
          return response.json();
        })
        .then((data) => {
          setLoans(data);
          setLoading(false);
        })
        .catch((err) => {
          setError(err.message);
          setLoading(false);
        });
    } else {
      setError("No hay usuario logueado");
      setLoading(false);
    }
  }, []);

  // Renderizar la tabla con el historial de préstamos
  return (
    <div>
      <h1>Historial de Préstamos Completados</h1>
      
      {loading && <p>Cargando...</p>}
      {error && <p>{error}</p>}
      
      {!loading && loans.length > 0 && (
        <table border="1" cellPadding="10" style={{ width: "100%", textAlign: "left", borderCollapse: "collapse" }}>
          <thead>
            <tr>
              <th>ID Préstamo</th>
              <th>Fecha Préstamo</th>
              <th>Fecha Devolución</th>
              <th>Estado</th>
              <th>Título del Libro</th>
            </tr>
          </thead>
          <tbody>
            {loans.map((loan) => (
              <tr key={loan.id_prestamo}>
                <td>{loan.id_prestamo}</td>
                <td>{loan.fecha_prestamo}</td>
                <td>{loan.fecha_devolucion}</td>
                <td>{loan.estado}</td>
                <td>{loan.titulo_libro}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {!loading && loans.length === 0 && <p>No hay préstamos completados para este usuario</p>}
    </div>
  );
};

export default Historial;
