import React, { useEffect, useState } from 'react';

const Multa = () => {
  const [fines, setFines] = useState([]);  // Estado para las multas
  const [loading, setLoading] = useState(true);  // Estado de carga
  const [error, setError] = useState(null);  // Estado de error

  useEffect(() => {
    // Recuperar el usuario desde localStorage
    const user = JSON.parse(localStorage.getItem("user"));
    
    if (user) {
      console.log("Usuario logueado:", user);
      console.log("Rol del usuario:", user.rol);

      // Realizar la solicitud GET para obtener las multas pendientes
      fetch(`http://localhost:4000/api/fines?usuario_id=${user.id}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("No existen multas pendientes");
          }
          return response.json();
        })
        .then((data) => {
          // Si no hay multas pendientes, mostrar un mensaje
          if (data.message) {
            setError(data.message);  // Mostrar mensaje en caso de que no haya multas
          } else {
            setFines(data);  // Guardar multas pendientes en el estado
          }
          setLoading(false);  // Terminar el estado de carga
        })
        .catch((err) => {
          setError(err.message);  // Capturar errores y mostrar el mensaje
          setLoading(false);  // Terminar el estado de carga
        });
    } else {
      setError("No hay usuario logueado");  // Si no hay usuario logueado
      setLoading(false);  // Terminar el estado de carga
    }
  }, []);

  return (
    <div>
      <h1>Multas Pendientes</h1>
      
      {/* Mostrar mensaje de carga */}
      {loading && <p>Cargando...</p>}

      {/* Mostrar mensaje de error */}
      {error && !loading && <p style={{ color: 'red' }}>{error}</p>}

      {/* Mostrar tabla si hay multas */}
      {!loading && fines.length > 0 && (
        <table border="1" cellPadding="10" style={{ width: "100%", textAlign: "left", borderCollapse: "collapse" }}>
          <thead>
            <tr>
              <th>ID Multa</th>
              <th>Saldo a Pagar</th>
              <th>Fecha de Multa</th>
              <th>Estado</th>
            </tr>
          </thead>
          <tbody>
            {fines.map((fine) => (
              <tr key={fine.id_multa}>
                <td>{fine.id_multa}</td>
                <td>{fine.saldo_pagar}</td>
                <td>{fine.fecha_multa}</td>
                <td>{fine.estado}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Si no hay multas, mostrar un mensaje */}
      {!loading && fines.length === 0 && !error && (
        <p>No hay multas pendientes para este usuario.</p>
      )}
    </div>
  );
};

export default Multa;
