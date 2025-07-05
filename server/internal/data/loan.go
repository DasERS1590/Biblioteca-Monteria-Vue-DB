package data

import (
	"context"
	"database/sql"
	"time"
)

type Loan struct {
	UsuarioID       int    `json:"usuario_id"`
	LibroID         int    `json:"libro_id"`
	FechaPrestamo   string `json:"fecha_prestamo"`
	FechaDevolucion string `json:"fecha_devolucion"`
}

type Loans struct {
	IDPrestamo      int    `json:"idprestamo"`
	IDSocio         int    `json:"idsocio"`
	IDLibro         int    `json:"idlibro"`
	FechaPrestamo   string `json:"fechaprestamo"`
	FechaDevolucion string `json:"fechadevolucion"`
	Estado          string `json:"estado"`
}

type LoanHystory struct {
	IDPrestamo      int    `json:"id_prestamo"`
	IDUsuario       int    `json:"id_usuario"`
	IDLibro         int    `json:"id_libro"`
	TituloLibro     string `json:"titulo_libro"`
	FechaPrestamo   string `json:"fecha_prestamo"`
	FechaDevolucion string `json:"fecha_devolucion"`
	Estado          string `json:"estado"`
}

type LoanActive struct {
	IDPrestamo      int    `json:"id_prestamo"`
	FechaPrestamo   string `json:"fecha_prestamo"`
	FechaDevolucion string `json:"fecha_devolucion"`
	Estado          string `json:"estado"`
	TituloLibro     string `json:"titulo_libro"`
}

type LoanCompleted struct {
	IDPrestamo      int    `json:"id_prestamo"`
	FechaPrestamo   string `json:"fecha_prestamo"`
	FechaDevolucion string `json:"fecha_devolucion"`
	Estado          string `json:"estado"`
	TituloLibro     string `json:"titulo_libro"`
}

type LoanAdmin struct {
	IDPrestamo      int    `json:"idprestamo"`
	IDSocio         int    `json:"idsocio"`
	NombreSocio     string `json:"nombre_socio"`
	IDLibro         int    `json:"idlibro"`
	TituloLibro     string `json:"titulo_libro"`
	FechaPrestamo   string `json:"fechaprestamo"`
	FechaDevolucion string `json:"fechadevolucion"`
	Estado          string `json:"estado"`
}

type LoanModel struct {
	DB *sql.DB
}

func (m LoanModel) CreateLoan(input *Loan) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(
		ctx,
		"CALL realizarPrestamo(?, ?, ?, ?)",
		input.UsuarioID,
		input.LibroID,
		input.FechaPrestamo,
		input.FechaDevolucion,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m LoanModel) ExtendLoand(reservationid string, nuevaFecha time.Time) error {

	var estado string
	var fechaDevolucion time.Time

	err := m.DB.QueryRow("SELECT estado, fechadevolucion FROM prestamo WHERE idprestamo = ?", reservationid).Scan(&estado, &fechaDevolucion)
	if err != nil {
		return err
	}

	if estado == "completado" {
		return err
	}

	updateQuery := "UPDATE prestamo SET fechadevolucion = ? WHERE idprestamo = ?"
	_, err = m.DB.Exec(updateQuery, nuevaFecha, reservationid)
	if err != nil {
		return err
	}
	return nil

}

func (m LoanModel) GetActiveLoans(idsocio, startdate, enddate string) ([]*Loans, error) {

	query := `
			SELECT idprestamo, idsocio, idlibro, fechaprestamo, fechadevolucion, estado
			FROM prestamo
			WHERE idsocio = ? AND estado = 'activo' AND fechaprestamo BETWEEN ? AND ?
		`
	rows, err := m.DB.Query(query, idsocio, startdate, enddate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*Loans

	for rows.Next() {
		var loan Loans

		err := rows.Scan(&loan.IDPrestamo, &loan.IDSocio, &loan.IDLibro, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado)
		if err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	return loans, nil
}

func (m LoanModel) GetUserLoanHistory(idsocio string) ([]*LoanHystory, error) {
	query := `
	SELECT 
		prestamo.idprestamo,
		prestamo.idsocio,
		prestamo.idlibro,
		prestamo.fechaprestamo,
		prestamo.fechadevolucion,
		prestamo.estado,
		libro.titulo AS titulo_libro
	FROM 
		prestamo
	INNER JOIN 
		libro ON prestamo.idlibro = libro.idlibro
	WHERE 
		prestamo.idsocio = ?
	ORDER BY 
		prestamo.fechaprestamo DESC
`

	rows, err := m.DB.Query(query, idsocio)

	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var loans []*LoanHystory

	for rows.Next() {
		var loan LoanHystory
		err := rows.Scan(
			&loan.IDPrestamo,
			&loan.IDUsuario,
			&loan.IDLibro,
			&loan.FechaPrestamo,
			&loan.FechaDevolucion,
			&loan.Estado,
			&loan.TituloLibro,
		)
		if err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (m LoanModel) GetUserActiveLoanStatus(usuarioID string) ([]*LoanActive, error) {

	query := `
	SELECT 
		prestamo.idprestamo,
		prestamo.fechaprestamo,
		prestamo.fechadevolucion,
		prestamo.estado,
		libro.titulo AS titulo_libro
	FROM 
		prestamo
	JOIN 
		libro ON prestamo.idlibro = libro.idlibro
	WHERE 
		prestamo.idsocio = ? AND prestamo.estado = 'activo'
	`

	rows, err := m.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*LoanActive

	for rows.Next() {
		var loan LoanActive
		err := rows.Scan(&loan.IDPrestamo, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado, &loan.TituloLibro)
		if err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (m LoanModel) GetUserCompletedLoanHistory(usuarioID string) ([]*LoanCompleted, error) {

	query := `
	SELECT 
		prestamo.idprestamo,
		prestamo.fechaprestamo,
		prestamo.fechadevolucion,
		prestamo.estado,
		libro.titulo AS titulo_libro
	FROM 
		prestamo
	JOIN 
		libro ON prestamo.idlibro = libro.idlibro
	WHERE 
		prestamo.idsocio = ? AND prestamo.estado = 'completado'
	`

	rows, err := m.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*LoanCompleted

	for rows.Next() {
		var loan LoanCompleted
		err := rows.Scan(&loan.IDPrestamo, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado, &loan.TituloLibro)
		if err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

// GetLoanByID obtiene un préstamo específico por su ID
func (m LoanModel) GetLoanByID(loanID string) (*Loans, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT idprestamo, idsocio, idlibro, fechaprestamo, fechadevolucion, estado
		FROM prestamo
		WHERE idprestamo = ?
	`

	var loan Loans
	err := m.DB.QueryRowContext(ctx, query, loanID).Scan(
		&loan.IDPrestamo,
		&loan.IDSocio,
		&loan.IDLibro,
		&loan.FechaPrestamo,
		&loan.FechaDevolucion,
		&loan.Estado,
	)

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// CompleteLoan marca un préstamo como completado
func (m LoanModel) CompleteLoan(loanID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE prestamo SET estado = 'completado' WHERE idprestamo = ?"
	_, err := m.DB.ExecContext(ctx, query, loanID)
	return err
}

// GetActiveLoansByDateRange obtiene todos los préstamos activos en un rango de fechas con información completa
func (m LoanModel) GetActiveLoansByDateRange(startdate, enddate string) ([]*LoanAdmin, error) {
	query := `
		SELECT 
			prestamo.idprestamo,
			prestamo.idsocio,
			socio.nombre AS nombre_socio,
			prestamo.idlibro,
			libro.titulo AS titulo_libro,
			prestamo.fechaprestamo,
			prestamo.fechadevolucion,
			prestamo.estado
		FROM prestamo
		JOIN socio ON prestamo.idsocio = socio.idsocio
		JOIN libro ON prestamo.idlibro = libro.idlibro
		WHERE prestamo.estado = 'activo' AND prestamo.fechaprestamo BETWEEN ? AND ?
		ORDER BY prestamo.fechaprestamo DESC
	`
	
	rows, err := m.DB.Query(query, startdate, enddate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*LoanAdmin

	for rows.Next() {
		var loan LoanAdmin
		err := rows.Scan(
			&loan.IDPrestamo,
			&loan.IDSocio,
			&loan.NombreSocio,
			&loan.IDLibro,
			&loan.TituloLibro,
			&loan.FechaPrestamo,
			&loan.FechaDevolucion,
			&loan.Estado,
		)
		if err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}


