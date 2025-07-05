package data

import (
	"context"
	"database/sql"
	"time"
)

type Fine struct {
	IDMulta    int     `json:"idmulta"`
	IDPrestamo int     `json:"idprestamo"`
	SaldoPagar float64 `json:"saldopagar"`
	FechaMulta string  `json:"fechamulta"`
	Estado     string  `json:"estado"`
}

type FineWithUser struct {
	IDMulta      int     `json:"idmulta"`
	IDPrestamo   int     `json:"idprestamo"`
	IDSocio      int     `json:"idsocio"`
	NombreSocio  string  `json:"nombre_socio"`
	SaldoPagar   float64 `json:"saldopagar"`
	FechaMulta   string  `json:"fechamulta"`
	Estado       string  `json:"estado"`
}

type FineModel struct {
	DB *sql.DB
}

func (m FineModel) GetPendingFines() ([]*FineWithUser, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			m.idmulta, 
			m.idprestamo, 
			s.idsocio,
			s.nombre AS nombre_socio,
			m.saldopagar, 
			m.fechamulta, 
			m.estado
		FROM multa m
		JOIN prestamo p ON m.idprestamo = p.idprestamo
		JOIN socio s ON p.idsocio = s.idsocio
		WHERE m.estado = 'pendiente'
		ORDER BY m.fechamulta DESC
	`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pendingFines []*FineWithUser

	for rows.Next() {
		var fine FineWithUser

		err := rows.Scan(&fine.IDMulta, &fine.IDPrestamo, &fine.IDSocio, &fine.NombreSocio, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			return nil, err
		}
		pendingFines = append(pendingFines, &fine)
	}
	return pendingFines, nil
}

func (m FineModel) GetUserFines(idsocio string) ([]*Fine, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT m.idmulta, m.idprestamo, m.saldopagar, m.fechamulta, m.estado
		FROM multa m
		INNER JOIN 
			prestamo p ON m.idprestamo = p.idprestamo
		WHERE p.idsocio = ?
		ORDER BY m.fechamulta DESC
	`
	rows, err := m.DB.QueryContext(ctx, query, idsocio)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fines []*Fine

	for rows.Next() {
		var fine Fine
		err := rows.Scan(&fine.IDMulta, &fine.IDPrestamo, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			return nil, err
		}
		fines = append(fines, &fine)
	}

	return fines, nil

}

func (m FineModel) GetUserPendingFines(usuario_id string) ([]*Fine, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			multa.idmulta,
			multa.saldopagar,
			multa.fechamulta,
			multa.estado
		FROM 
			multa
		JOIN 
			prestamo ON multa.idprestamo = prestamo.idprestamo
		WHERE 
			prestamo.idsocio = ? AND multa.estado = 'pendiente'
	`

	rows, err := m.DB.QueryContext(ctx, query, usuario_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fines []*Fine

	for rows.Next() {
		var fine Fine
		err := rows.Scan(&fine.IDMulta, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			return nil, err
		}
		fines = append(fines, &fine)
	}

	if err := rows.Err(); err != nil {
		return nil , err
	}

	return fines , nil 
}

// CreateFine crea una nueva multa
func (m FineModel) CreateFine(fine *Fine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO multa (idprestamo, saldopagar, fechamulta, estado)
		VALUES (?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, query, fine.IDPrestamo, fine.SaldoPagar, fine.FechaMulta, fine.Estado)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fine.IDMulta = int(id)
	return nil
}

// GetFineByID obtiene una multa específica por su ID
func (m FineModel) GetFineByID(fineID string) (*Fine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT idmulta, idprestamo, saldopagar, fechamulta, estado
		FROM multa
		WHERE idmulta = ?
	`

	var fine Fine
	err := m.DB.QueryRowContext(ctx, query, fineID).Scan(
		&fine.IDMulta,
		&fine.IDPrestamo,
		&fine.SaldoPagar,
		&fine.FechaMulta,
		&fine.Estado,
	)

	if err != nil {
		return nil, err
	}

	return &fine, nil
}

// PayFine marca una multa como pagada
func (m FineModel) PayFine(fineID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE multa SET estado = 'pagada' WHERE idmulta = ?"
	_, err := m.DB.ExecContext(ctx, query, fineID)
	return err
}

// SearchFinesByUser busca multas por nombre o email del usuario
func (m FineModel) SearchFinesByUser(busqueda string) ([]*FineWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			m.idmulta, 
			m.idprestamo, 
			s.idsocio,
			s.nombre AS nombre_socio,
			m.saldopagar, 
			m.fechamulta, 
			m.estado
		FROM multa m
		JOIN prestamo p ON m.idprestamo = p.idprestamo
		JOIN socio s ON p.idsocio = s.idsocio
		WHERE (s.nombre LIKE ? OR s.correo LIKE ?)
		ORDER BY m.fechamulta DESC
	`
	
	searchTerm := "%" + busqueda + "%"
	rows, err := m.DB.QueryContext(ctx, query, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fines []*FineWithUser

	for rows.Next() {
		var fine FineWithUser
		err := rows.Scan(&fine.IDMulta, &fine.IDPrestamo, &fine.IDSocio, &fine.NombreSocio, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			return nil, err
		}
		fines = append(fines, &fine)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fines, nil
}
