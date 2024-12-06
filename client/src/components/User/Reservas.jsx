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
      maxWidth: "800px",
      margin: "0 auto",
    },
    header: {
      textAlign: "center",
      color: "#333",
      marginBottom: "20px",
    },
    error: {
      color: "red",
      fontWeight: "bold",
      textAlign: "center",
    },
    list: {
      listStyleType: "none",
      padding: 0,
    },
    listItem: {
      backgroundColor: "#f9f9f9",
      margin: "10px 0",
      padding: "15px",
      borderRadius: "8px",
      boxShadow: "0 2px 4px rgba(0, 0, 0, 0.1)",
      borderLeft: "4px solid #4CAF50",
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
        <ul style={styles.list}>
          {reservations.map((reservation) => (
            <li key={reservation.id_reserva} style={styles.listItem}>
              <strong>ID Reserva:</strong> {reservation.id_reserva} |{" "}
              <strong>ID Libro:</strong> {reservation.id_libro} |{" "}
              <strong>Fecha Reserva:</strong> {reservation.fecha_reserva} |{" "}
              <strong>Estado:</strong> {reservation.estado}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default Reserva;