import React, { useState, useEffect } from "react";
import axios from "axios";

const Fines = () => {
  const [fines, setFines] = useState([]); // Estado para almacenar las multas
  const [userID, setUserID] = useState(""); // Estado para almacenar el ID del usuario
  const [loading, setLoading] = useState(false); // Estado de carga
  const [error, setError] = useState(null); // Estado de error
  const [filterState, setFilterState] = useState("pendiente"); // Filtro de estado (pendiente/pagada)

  // Función para obtener las multas
  const fetchFines = async () => {
    setLoading(true);
    try {
      let url = "http://localhost:4000/api/admin/fines/to"; // Ruta por defecto
      let params = { estado: filterState }; // Filtros para el estado

      // Si hay un userID, cambiamos la URL y pasamos ese parámetro
      if (userID) {
        url = "http://localhost:4000/api/admin/fines"; // Ruta para obtener las multas de un usuario específico
        params.idsocio = userID; // Añadimos el parámetro del ID de usuario
      }

      const response = await axios.get(url, { params });

      if (response.data && Array.isArray(response.data)) {
        setFines(response.data);
      } else {
        setError("No se encontraron multas.");
      }
    } catch (err) {
      setError("Error al obtener las multas.");
    } finally {
      setLoading(false);
    }
  };

  // Ejecutar la función cada vez que cambie el filtro o el userID
  useEffect(() => {
    fetchFines();
  }, [filterState, userID]);

  return (
    <div style={styles.dashboard}>
      <h1 style={styles.header}>Gestión de Multas</h1>

      {/* Filtro para buscar las multas de un usuario */}
      <div style={styles.filterSection}>
        <input
          type="text"
          placeholder="ID de usuario (opcional)"
          value={userID}
          onChange={(e) => setUserID(e.target.value)}
          style={{ ...styles.input, textAlign: "right" }} // Alineación a la derecha
        />
      </div>

      {/* Mostrando las multas */}
      <div style={styles.finesSection}>
        <h2 style={styles.subheader}>{filterState === "pendiente" ? "Multas Pendientes" : "Multas Pagadas"}</h2>
        {loading && <p>Cargando...</p>}
        {error && <p style={styles.error}>{error}</p>}
        {!loading && !error && (
          <table style={styles.table}>
            <thead>
              <tr>
                <th>ID Multa</th>
                <th>ID Prestamo</th>
                <th>Saldo Pagar</th>
                <th>Fecha Multa</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              {fines.map((fine) => (
                <tr key={fine.idmulta}>
                  <td>{fine.idmulta}</td>
                  <td>{fine.idprestamo}</td>
                  <td>{fine.saldopagar}</td>
                  <td>{fine.fechamulta}</td>
                  <td>{fine.estado}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

const styles = {
  dashboard: {
    padding: "20px",
    fontFamily: "Arial, sans-serif",
  },
  header: {
    textAlign: "center",
    fontSize: "24px",
    marginBottom: "20px",
  },
  filterSection: {
    display: "flex",
    flexDirection: "row",  // Cambiar a fila para alinear los elementos horizontalmente
    justifyContent: "space-between", // Para distribuir los elementos en la fila
    width: "100%",
    marginBottom: "20px",
  },
  input: {
    padding: "8px",
    marginRight: "10px", // Espacio entre el input y el select
    width: "200px",
    borderRadius: "4px",
    border: "1px solid #ccc",
  },
  selectContainer: {
    display: "flex",
    alignItems: "center",
  },
  select: {
    padding: "8px",
    width: "200px",
    borderRadius: "4px",
    border: "1px solid #ccc",
  },
  finesSection: {
    marginBottom: "20px",
  },
  subheader: {
    fontSize: "20px",
    marginBottom: "10px",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
  },
  error: {
    color: "red",
    textAlign: "center",
  },
};

export default Fines;
