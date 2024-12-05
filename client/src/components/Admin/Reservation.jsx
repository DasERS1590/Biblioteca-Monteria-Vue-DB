import React, { useState } from "react";
import axios from "axios";
import "../../styles/Admin.css";

const Reservation = () => {
  const [userID, setUserID] = useState("");
  const [bookID, setBookID] = useState("");
  const [date, setDate] = useState("");
  const [nameSocio, setNameSocio] = useState("");
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchReservations = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await axios.get("http://localhost:4000/api/admin/reservations", {
        params: {
          usuarioid: userID,
          libro: bookID,
          fecha: date,
          nombre: nameSocio,
        },
      });

      setReservations(response.data);
    } catch (err) {
      setError("Error al obtener las reservas. Por favor, intenta más tarde.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-dashboard">
      <h1>Reservas Activas</h1>
      <p>Consulta las reservas activas según los filtros proporcionados.</p>

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
          type="text"
          placeholder="ID del Libro"
          value={bookID}
          onChange={(e) => setBookID(e.target.value)}
          className="input"
        />
        <input
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          className="input"
        />
        <input
          type="text"
          placeholder="Nombre del Socio"
          value={nameSocio}
          onChange={(e) => setNameSocio(e.target.value)}
          className="input"
        />
        <button onClick={fetchReservations} className="btn">
          Buscar
        </button>
      </div>

      {/* Indicador de carga */}
      {loading && <p className="loading">Cargando reservas...</p>}

      {/* Error */}
      {error && <p className="error">{error}</p>}

      {/* Resultados */}
      {!loading && reservations.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th>ID Reserva</th>
              <th>ID Socio</th>
              <th>Nombre Socio</th>
              <th>ID Libro</th>
              <th>Título Libro</th>
              <th>Fecha Reserva</th>
              <th>Estado Reserva</th>
              <th>Género Libro</th>
              <th>Editorial</th>
              <th>Teléfono Socio</th>
              <th>Correo Socio</th>
              <th>Tipo Socio</th>
              <th>Fecha Nacimiento</th>
              <th>Fecha Registro</th>
            </tr>
          </thead>
          <tbody>
            {reservations.map((reservation) => (
              <tr key={reservation.idreserva}>
                <td>{reservation.idreserva}</td>
                <td>{reservation.idsocio}</td>
                <td>{reservation.nombre_socio}</td>
                <td>{reservation.idlibro}</td>
                <td>{reservation.titulo_libro}</td>
                <td>{reservation.fechareserva}</td>
                <td>{reservation.estado_reserva}</td>
                <td>{reservation.genero_libro}</td>
                <td>{reservation.editorial}</td>
                <td>{reservation.telefono_socio}</td>
                <td>{reservation.correo_socio}</td>
                <td>{reservation.tiposocio}</td>
                <td>{reservation.fechanacimiento}</td>
                <td>{reservation.fecharegistro}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Sin resultados */}
      {!loading && reservations.length === 0 && !error && (
        <p className="no-results">No se encontraron reservas activas.</p>
      )}
    </div>
  );
};

export default Reservation;
