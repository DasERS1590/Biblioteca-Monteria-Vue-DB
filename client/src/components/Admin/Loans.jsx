import React, { useState } from "react";
import axios from "axios";
import "../../styles/Admin.css";

const Loans = () => {
  const [userID, setUserID] = useState("");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [loans, setLoans] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Función para obtener los préstamos activos
  const fetchLoans = async () => {
    if (!userID || !startDate || !endDate) {
      setError("Por favor, completa todos los campos.");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await axios.get("http://localhost:4000/api/admin/loans", {
        params: { idsocio: userID, startdate: startDate, enddate: endDate },
      });
      setLoans(response.data);
    } catch (err) {
      setError("Error al obtener los préstamos. Por favor, intenta más tarde.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-dashboard">
      <h1>Préstamos Activos</h1>
      <p>Consulta los préstamos activos según el socio y las fechas.</p>

      {/* Formulario de filtros */}
      <div className="filters">
        <input
          type="text"
          placeholder="ID del Socio"
          value={userID}
          onChange={(e) => setUserID(e.target.value)}
          className="input"
        />
        <input
          type="date"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
          className="input"
        />
        <input
          type="date"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
          className="input"
        />
        <button onClick={fetchLoans} className="btn">
          Buscar
        </button>
      </div>

      {/* Indicador de carga */}
      {loading && <p className="loading">Cargando préstamos...</p>}

      {/* Error */}
      {error && <p className="error">{error}</p>}

      {/* Resultados */}
      {!loading && loans.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th>ID Préstamo</th>
              <th>ID Socio</th>
              <th>ID Libro</th>
              <th>Fecha Préstamo</th>
              <th>Fecha Devolución</th>
              <th>Estado</th>
            </tr>
          </thead>
          <tbody>
            {loans.map((loan) => (
              <tr key={loan.idprestamo}>
                <td>{loan.idprestamo}</td>
                <td>{loan.idsocio}</td>
                <td>{loan.idlibro}</td>
                <td>{loan.fechaprestamo}</td>
                <td>{loan.fechadevolucion}</td>
                <td>{loan.estado}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Sin resultados */}
      {!loading && loans.length === 0 && !error && (
        <p className="no-results">No se encontraron préstamos activos.</p>
      )}
    </div>
  );
};

export default Loans;
