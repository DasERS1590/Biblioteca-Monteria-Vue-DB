import React, { useState, useEffect } from "react";
import axios from "axios";

const User = () => {
  const [userType, setUserType] = useState(""); // Estado para el tipo de socio
  const [users, setUsers] = useState([]); // Estado para la lista de usuarios
  const [loading, setLoading] = useState(false); // Estado para el indicador de carga
  const [error, setError] = useState(null); // Estado para errores

  // Función para obtener los usuarios según el tipo de socio
  const fetchUsers = async () => {
    if (!userType) {
      setError("Por favor, selecciona un tipo de socio.");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const response = await axios.get("http://localhost:4000/api/admin/users", {
        params: { tiposocio: userType },
      });
      setUsers(response.data);
    } catch (err) {
      setError("Error al obtener los usuarios. Por favor, intenta más tarde.");
    } finally {
      setLoading(false);
    }
  };

  // Maneja el cambio en el selector de tipo de socio
  const handleUserTypeChange = (e) => {
    setUserType(e.target.value);
  };

  return (
    <div style={styles.container}>
      <h1 style={styles.title}>Lista de Usuarios</h1>
      <p style={styles.description}>
        Selecciona un tipo de socio para ver los usuarios correspondientes.
      </p>

      {/* Filtros */}
      <div style={styles.filterContainer}>
        <select value={userType} onChange={handleUserTypeChange} style={styles.select}>
          <option value="">Seleccione un tipo de socio</option>
          <option value="normal">Normal</option>
          <option value="estudiante">Estudiante</option>
          <option value="profesor">Profesor</option>
        </select>
        <button onClick={fetchUsers} style={styles.button}>
          Buscar
        </button>
      </div>

      {/* Indicador de carga */}
      {loading && <p style={styles.loading}>Cargando usuarios...</p>}

      {/* Error */}
      {error && <p style={styles.error}>{error}</p>}

      {/* Tabla de resultados */}
      {!loading && users.length > 0 && (
        <table style={styles.table}>
          <thead style={styles.tableHead}>
            <tr>
              <th style={styles.tableCell}>ID</th>
              <th style={styles.tableCell}>Nombre</th>
              <th style={styles.tableCell}>Dirección</th>
              <th style={styles.tableCell}>Teléfono</th>
              <th style={styles.tableCell}>Correo</th>
              <th style={styles.tableCell}>Tipo de Socio</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} style={styles.tableRow}>
                <td style={styles.tableCell}>{user.id}</td>
                <td style={styles.tableCell}>{user.nombre}</td>
                <td style={styles.tableCell}>{user.direccion}</td>
                <td style={styles.tableCell}>{user.telefono}</td>
                <td style={styles.tableCell}>{user.correo}</td>
                <td style={styles.tableCell}>{user.tiposocio}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Sin resultados */}
      {!loading && users.length === 0 && !error && <p style={styles.noResults}>No se encontraron usuarios.</p>}
    </div>
  );
};

// Estilos en línea
const styles = {
  container: {
    fontFamily: "Arial, sans-serif",
    padding: "20px",
    maxWidth: "800px",
    margin: "0 auto",
    backgroundColor: "#f9f9f9",
    borderRadius: "8px",
    boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
  },
  title: {
    fontSize: "24px",
    color: "#333",
    marginBottom: "10px",
  },
  description: {
    fontSize: "16px",
    color: "#555",
    marginBottom: "20px",
  },
  filterContainer: {
    display: "flex",
    gap: "10px",
    marginBottom: "20px",
  },
  select: {
    padding: "10px",
    fontSize: "14px",
    borderRadius: "4px",
    border: "1px solid #ccc",
    flex: "1",
  },
  button: {
    padding: "10px 20px",
    fontSize: "14px",
    backgroundColor: "#007bff",
    color: "#fff",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer",
  },
  loading: {
    fontSize: "16px",
    color: "#007bff",
  },
  error: {
    fontSize: "14px",
    color: "red",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
    marginTop: "20px",
  },
  tableHead: {
    backgroundColor: "#007bff",
    color: "#fff",
  },
  tableRow: {
    borderBottom: "1px solid #ddd",
  },
  tableCell: {
    padding: "10px",
    textAlign: "left",
  },
  noResults: {
    fontSize: "14px",
    color: "#555",
    textAlign: "center",
  },
};

export default User;
