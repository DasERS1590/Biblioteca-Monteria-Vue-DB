package data

import "database/sql"

type Reservation struct {
	SocioID      int    `json:"idsocio"`
	LibroID      int    `json:"idlibro"`
	FechaReserva string `json:"fechareserva"`
}

type ReservationData struct {
	IdReserva       int    `json:"idreserva"`
	IdSocio         int    `json:"idsocio"`
	NombreSocio     string `json:"nombre_socio"`
	IdLibro         int    `json:"idlibro"`
	TituloLibro     string `json:"titulo_libro"`
	FechaReserva    string `json:"fechareserva"`
	EstadoReserva   string `json:"estado_reserva"`
	GeneroLibro     string `json:"genero_libro"`
	Editorial       string `json:"editorial"`
	TelefonoSocio   string `json:"telefono_socio"`
	CorreoSocio     string `json:"correo_socio"`
	TipoSocio       string `json:"tiposocio"`
	FechaNacimiento string `json:"fechanacimiento"`
	FechaRegistro   string `json:"fecharegistro"`
}

type ReservationUser struct {
	IDReserva    int    `json:"id_reserva"`
	IDSocio      int    `json:"id_socio"`
	IDLibro      int    `json:"id_libro"`
	FechaReserva string `json:"fecha_reserva"`
	Estado       string `json:"estado"`
}

type ReservationModel struct {
	DB *sql.DB
}

func (m ReservationModel) CreateReservation(reservation *Reservation) error {

	insertReservationQuery := `
        INSERT INTO reserva (idsocio, idlibro, fechareserva, estado)
        VALUES (?, ?, ?, 'activa')`

	_, err := m.DB.Exec(insertReservationQuery, reservation.SocioID, reservation.LibroID, reservation.FechaReserva)
	if err != nil {
		return err
	}

	return nil
}

func (m ReservationModel) CancelReservation(reservationid string) error {

	var estado string
	err := m.DB.QueryRow("SELECT estado FROM reserva WHERE idreserva = ?", reservationid).Scan(&estado)
	if err != nil {
		return err
	}

	if estado != "activa" {
		return err
	}

	updateQuery := "UPDATE reserva SET estado = 'cancelada' WHERE idreserva = ?"
	_, err = m.DB.Exec(updateQuery, reservationid)
	if err != nil {
		return err
	}

	return nil
}

func (m ReservationModel) GetActiveReservations(usuarioid string, libro string, fecha string, nombreSocio string) ([]*ReservationData, error) {
	query := `
        SELECT 
            r.idreserva,
            r.idsocio,
            s.nombre AS nombre_socio,
            r.idlibro,
            l.titulo AS titulo_libro,
            r.fechareserva,
            r.estado AS estado_reserva,
            l.genero AS genero_libro,
            e.nombre AS editorial,
            s.telefono AS telefono_socio,
            s.correo AS correo_socio,
            s.tiposocio,
            s.fechanacimiento,
            s.fecharegistro
        FROM 
            reserva r
        JOIN 
            socio s ON r.idsocio = s.idsocio
        JOIN 
            libro l ON r.idlibro = l.idlibro
        JOIN 
            editorial e ON l.ideditorial = e.ideditorial
        WHERE 
            r.estado = 'activa'
    `

	var params []interface{}

	if usuarioid != "" {
		query += " AND r.idsocio = ?"
		params = append(params, usuarioid)
	}
	if libro != "" {
		query += " AND r.idlibro = ?"
		params = append(params, libro)
	}
	if fecha != "" {
		query += " AND r.fechareserva = ?"
		params = append(params, fecha)
	}
	if nombreSocio != "" {
		query += " AND s.nombre LIKE ?"
		params = append(params, "%"+nombreSocio+"%")
	}

	query += " ORDER BY r.fechareserva DESC"

	rows, err := m.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*ReservationData
	for rows.Next() {
		var reservation ReservationData
		if err := rows.Scan(&reservation.IdReserva, &reservation.IdSocio, &reservation.NombreSocio,
			&reservation.IdLibro, &reservation.TituloLibro, &reservation.FechaReserva, &reservation.EstadoReserva,
			&reservation.GeneroLibro, &reservation.Editorial, &reservation.TelefonoSocio, &reservation.CorreoSocio,
			&reservation.TipoSocio, &reservation.FechaNacimiento, &reservation.FechaRegistro); err != nil {

			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, err
}

func (m ReservationModel) GetUserActiveReservations(userid string) ([]*ReservationUser, error) {

	query := `
	SELECT 
		reserva.idreserva,
		reserva.idsocio,
		reserva.idlibro,
		reserva.fechareserva,
		reserva.estado
	FROM 
		reserva
	WHERE 
		reserva.idsocio = ? AND reserva.estado = 'activa'
`

	rows, err := m.DB.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*ReservationUser

	for rows.Next() {
		var reservation ReservationUser
		err := rows.Scan(&reservation.IDReserva, &reservation.IDSocio, &reservation.IDLibro, &reservation.FechaReserva, &reservation.Estado)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reservations, err
}
