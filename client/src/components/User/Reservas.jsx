import React, { useEffect, useState } from "react";
import axios from "axios";

function Reserva() {
  const [reservations, setReservations] = useState([]);
  const [error, setError] = useState("");

  // Definir estilos en una variable
  const styles = {
    container: {
      padding: "20px",
      fontFamily: "Arial, sans-serif",
      maxWidth: "1000px",
      margin: "0 auto",
      backgroundColor: "#fff",
      borderRadius: "8px",
      boxShadow: "0 2px 10px rgba(0, 0, 0, 0.1)",
    },
    header: {
      textAlign: "center",
      color: "#333",
      marginBottom: "20px",
      fontSize: "24px",
    },
    error: {
      color: "red",
      fontWeight: "bold",
      textAlign: "center",
    },
    tableContainer: {
      overflowX: "auto",
    },
    table: {
      width: "100%",
      borderCollapse: "collapse",
      marginTop: "20px",
    },
    th: {
      backgroundColor: "#4CAF50",
      color: "white",
      padding: "12px",
      textAlign: "center",
    },
    td: {
      padding: "12px",
      textAlign: "center",
      borderBottom: "1px solid #ddd",
    },
    noReservations: {
      textAlign: "center",
      color: "#666",
      marginTop: "20px",
    },
  };

  useEffect(() => {
    const fetchReservations = async () => {
      try {
        setError(""); // Limpiar errores previos

        // Obtener el ID del usuario desde el localStorage
        const user = JSON.parse(localStorage.getItem("user"));

        if (!user) {
          setError("No se encontró el ID del usuario en el localStorage.");
          return;
        }

        // Realizar la solicitud al backend
        const response = await axios.get(
          `http://localhost:4000/api/reservations?usuario_id=${user.id}`
        );

        // Guardar las reservas obtenidas en el estado
        setReservations(response.data);
      } catch (err) {
        setReservations([]);
        setError(err.response?.data?.message || "Error al obtener reservas.");
      }
    };

    fetchReservations();
  }, []); // Ejecutar solo al cargar el componente

  return (
    <div style={styles.container}>
      <h1 style={styles.header}>Reservas Activas</h1>
      {error && <p style={styles.error}>{error}</p>}

      {reservations.length === 0 && !error ? (
        <p style={styles.noReservations}>No hay reservas activas para este usuario.</p>
      ) : (
        <div style={styles.tableContainer}>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>ID Reserva</th>
                <th style={styles.th}>ID Libro</th>
                <th style={styles.th}>Fecha Reserva</th>
                <th style={styles.th}>Estado</th>
              </tr>
            </thead>
            <tbody>
              {reservations.map((reservation) => (
                <tr key={reservation.id_reserva}>
                  <td style={styles.td}>{reservation.id_reserva}</td>
                  <td style={styles.td}>{reservation.id_libro}</td>
                  <td style={styles.td}>{reservation.fecha_reserva}</td>
                  <td style={styles.td}>{reservation.estado}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default Reserva;
