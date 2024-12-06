import React, { useEffect, useState } from "react";

const Prestamo = () => {
  const [loans, setLoans] = useState([]); // Estado para almacenar los préstamos
  const [loading, setLoading] = useState(true); // Estado de carga
  const [error, setError] = useState(null); // Estado de error

  useEffect(() => {
    // Obtener el usuario desde el localStorage
    const user = JSON.parse(localStorage.getItem("user"));

    if (user) {
      console.log("Usuario logueado:", user);

      // Llamar al endpoint para obtener los préstamos activos
      fetch(`http://localhost:4000/api/loans?usuario_id=${user.id}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("No se pudo obtener los préstamos activos");
          }
          return response.json();
        })
        .then((data) => {
          // Si el endpoint devuelve un mensaje en lugar de datos
          if (data.message) {
            setError(data.message);
          } else {
            setLoans(data); // Guardar los préstamos activos
          }
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

  return (
    <div>
      <h1>Préstamos Activos</h1>

      {/* Mostrar estado de carga */}
      {loading && <p>Cargando...</p>}

      {/* Mostrar errores */}
      {error && !loading && <p style={{ color: "red" }}>{error}</p>}

      {/* Mostrar la tabla si hay préstamos */}
      {!loading && loans.length > 0 && (
        <table
          border="1"
          cellPadding="10"
          style={{
            width: "100%",
            textAlign: "left",
            borderCollapse: "collapse",
          }}
        >
          <thead>
            <tr>
              <th>ID Préstamo</th>
              <th>Fecha de Préstamo</th>
              <th>Fecha de Devolución</th>
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

      {/* Mostrar mensaje si no hay préstamos */}
      {!loading && loans.length === 0 && !error && (
        <p>No hay préstamos activos para este usuario.</p>
      )}
    </div>
  );
};

export default Prestamo;
