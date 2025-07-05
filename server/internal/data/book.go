package data

import (
	"biblioteca/internal/validator"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Book struct {
	IdLibro          int    `json:"idlibro"`
	Titulo           string `json:"titulo"`
	Genero           string `json:"genero"`
	FechaPublicacion string `json:"fechapublicacion"`
	EditorialID      int    `json:"ideditorial"`
	AutoresId        []int  `json:"idautores"`
	Status           string
}

type BookResponseGenAut struct {
	IDLibro int    `json:"id_libro"`
	Titulo  string `json:"titulo"`
	Genero  string `json:"genero"`
	Estado  string `json:"estado"`
	Autor   string `json:"autor"`
}

type BookbyDate struct {
	IDLibro          int    `json:"id_libro"`
	Titulo           string `json:"titulo"`
	Genero           string `json:"genero"`
	FechaPublicacion string `json:"fecha_publicacion"`
	Estado           string `json:"estado"`
	Editorial        string `json:"editorial"`
}

type BookAbailable struct {
	IDLibro int    `json:"id_libro"`
	Titulo  string `json:"titulo"`
	Genero  string `json:"genero"`
	Estado  string `json:"estado"`
	Autor   string `json:"autor"`
	Editorial string `json:"editorial"`
}

type BookUnavailable struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
	Status string `json:"status"`
}

type BookForEdit struct {
	IDLibro          int    `json:"id_libro"`
	Titulo           string `json:"titulo"`
	Genero           string `json:"genero"`
	FechaPublicacion string `json:"fecha_publicacion"`
	EditorialID      int    `json:"editorial_id"`
	Editorial        string `json:"editorial"`
	Estado           string `json:"estado"`
	Autores          []int  `json:"autores"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Titulo != "", "titulo", "debe ser proporcionado")
	v.Check(book.Genero != "", "genero", "debe ser proporcionado")

	v.Check(book.FechaPublicacion != "", "fechapublicacion", "debe ser proporcionada")
	_, err := time.Parse("2006-01-02", book.FechaPublicacion)
	v.Check(err == nil, "fechapublicacion", "debe tener formato YYYY-MM-DD")

	v.Check(book.EditorialID != 0, "ideditorial", "debe de ser proporcionada")
	v.Check(book.AutoresId != nil, "idautores", "debe de ser proporcionada")
}

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) CreateBook(book *Book) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlQuery := `INSERT INTO libro (ideditorial, fechapublicacion,titulo, genero, estado)
        VALUES (?, ?, ?, ?,'disponible')`

	//fecha_publicacion , _ := time.Parse("2006-01-02", book.FechaPublicacion)
	args := []any{book.EditorialID, book.FechaPublicacion, book.Titulo, book.Genero}

	result, err := m.DB.ExecContext(ctx, sqlQuery, args...)

	if err != nil {
		return err
	}

	bookID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, autorID := range book.AutoresId {
		insertAuthorQuery := `
            INSERT INTO libro_autor (idlibro, idautor)
            VALUES (?, ?)`

		_, err = m.DB.Exec(insertAuthorQuery, bookID, autorID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m BookModel) UpdateBook(updatedBook *Book) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var existingBookCount int
	checkBookQuery := `SELECT COUNT(*) FROM libro WHERE idlibro = ?`

	bookId := updatedBook.IdLibro

	if err := m.DB.QueryRow(checkBookQuery, bookId).Scan(&existingBookCount); err != nil || existingBookCount == 0 {
		return err
	}

	updateBookQuery := `
			UPDATE libro
			SET titulo = ?, genero = ?, fechapublicacion = ?, ideditorial = ?
			WHERE idlibro = ?`

	args := []any{updatedBook.Titulo, updatedBook.Genero, updatedBook.FechaPublicacion, updatedBook.EditorialID, updatedBook.IdLibro}
	_, err := m.DB.ExecContext(ctx, updateBookQuery, args...)

	if err != nil {
		return err
	}

	deleteAuthorsQuery := `DELETE FROM libro_autor WHERE idlibro = ?`

	_, err = m.DB.Exec(deleteAuthorsQuery, bookId)

	if err != nil {
		return err
	}

	for _, autorID := range updatedBook.AutoresId {
		insertAuthorQuery := `INSERT INTO libro_autor (idlibro, idautor) VALUES (?, ?)`
		_, err := m.DB.Exec(insertAuthorQuery, bookId, autorID)
		if err != nil {

			return err
		}
	}

	return nil
}

func (m BookModel) GetFilteredBooks(estado string, editorial string) ([]*BookAbailable, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	querySQL := `
	SELECT 
		libro.idlibro, libro.titulo, libro.genero, libro.estado, 
		editorial.nombre AS editorial,
		COALESCE(autor.nombre, 'Sin autor') AS autor
	FROM 
		libro
	INNER JOIN 
		editorial ON libro.ideditorial = editorial.ideditorial
	LEFT JOIN 
		libro_autor ON libro.idlibro = libro_autor.idlibro
	LEFT JOIN 
		autor ON libro_autor.idautor = autor.idautor
	WHERE 
		(libro.estado = ? OR ? = '') 
		AND 
		(editorial.nombre = ? OR ? = '')`

	rows, err := m.DB.QueryContext(ctx, querySQL, estado, estado, editorial, editorial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookAbailable

	for rows.Next() {
		var book BookAbailable
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Editorial, &book.Autor)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m BookModel) GetBooksByGenreAndAuthor(genero string, autor string) ([]*BookResponseGenAut, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT 
		libro.idlibro,
		libro.titulo,
		libro.genero,
		libro.estado,
		autor.nombre AS autor
	FROM 
		libro
	JOIN 
		libro_autor ON libro.idlibro = libro_autor.idlibro
	JOIN 
		autor ON libro_autor.idautor = autor.idautor
	WHERE 
		libro.estado = 'disponible'
		AND libro.genero = ?
		AND autor.nombre LIKE ?
	`

	rows, err := m.DB.QueryContext(ctx, query, genero, "%"+autor+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookResponseGenAut

	for rows.Next() {
		var book BookResponseGenAut
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Autor)
		if err != nil {

			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m BookModel) GetBooksByPublicationDate(startDate, endDate string) ([]*BookbyDate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT 
				libro.idlibro,
				libro.titulo,
				libro.genero,
				libro.fechapublicacion,
				libro.estado,
				editorial.nombre AS editorial
			FROM 
				libro
			JOIN 
				editorial ON libro.ideditorial = editorial.ideditorial
			WHERE 
				libro.fechapublicacion BETWEEN ? AND ?
			`
	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookbyDate

	for rows.Next() {
		var book BookbyDate
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.FechaPublicacion, &book.Estado, &book.Editorial)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (m BookModel) GetBooksAvailableByGenreAndAuthor(genero, autor, titulo string) ([]*BookAbailable, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT 
            libro.idlibro,
            libro.titulo,
            libro.genero,
            libro.estado,
            autor.nombre AS autor
        FROM 
            libro
        JOIN 
            libro_autor ON libro.idlibro = libro_autor.idlibro
        JOIN 
            autor ON libro_autor.idautor = autor.idautor
        WHERE 
            libro.estado = 'disponible' 
    	`

	var params []interface{}

	if genero != "" {
		query += " AND libro.genero LIKE ?"
		params = append(params, "%"+genero+"%")
	}
	if autor != "" {
		query += " AND autor.nombre LIKE ?"
		params = append(params, "%"+autor+"%")
	}
	if titulo != "" {
		query += " AND libro.titulo LIKE ?"
		params = append(params, "%"+titulo+"%")
	}

	rows, err := m.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookAbailable

	for rows.Next() {
		var book BookAbailable
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Autor)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// GetBooksForReservation trae solo libros NO disponibles para reservas
func (m BookModel) GetBooksForReservation(genero, autor, titulo string) ([]*BookAbailable, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT 
            libro.idlibro,
            libro.titulo,
            libro.genero,
            libro.estado,
            autor.nombre AS autor
        FROM 
            libro
        JOIN 
            libro_autor ON libro.idlibro = libro_autor.idlibro
        JOIN 
            autor ON libro_autor.idautor = autor.idautor
        WHERE 
            libro.estado IN ('prestado', 'reservado')
    	`

	var params []interface{}

	if genero != "" {
		query += " AND libro.genero LIKE ?"
		params = append(params, "%"+genero+"%")
	}
	if autor != "" {
		query += " AND autor.nombre LIKE ?"
		params = append(params, "%"+autor+"%")
	}
	if titulo != "" {
		query += " AND libro.titulo LIKE ?"
		params = append(params, "%"+titulo+"%")
	}

	rows, err := m.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*BookAbailable

	for rows.Next() {
		var book BookAbailable
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Autor)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m BookModel) GetUnavailableBooks() ([]*BookUnavailable, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
    SELECT 
        libro.idlibro, libro.titulo, libro.genero, libro.estado
    FROM 
        libro
    WHERE 
        libro.estado IN ('prestado', 'reservado')`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var books []*BookUnavailable

	for rows.Next() {
		var book BookUnavailable
		if err := rows.Scan(&book.ID, &book.Title, &book.Genre, &book.Status); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// UpdateBookStatus actualiza el estado de un libro
func (m BookModel) UpdateBookStatus(bookID int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE libro SET estado = ? WHERE idlibro = ?"
	_, err := m.DB.ExecContext(ctx, query, status, bookID)
	return err
}

func (m BookModel) GetBookForEdit(bookID int) (*BookForEdit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Obtener información básica del libro
	query := `
		SELECT 
			libro.idlibro,
			libro.titulo,
			libro.genero,
			libro.fechapublicacion,
			libro.ideditorial,
			libro.estado,
			editorial.nombre AS editorial
		FROM 
			libro
		JOIN 
			editorial ON libro.ideditorial = editorial.ideditorial
		WHERE 
			libro.idlibro = ?
	`

	var book BookForEdit
	err := m.DB.QueryRowContext(ctx, query, bookID).Scan(
		&book.IDLibro,
		&book.Titulo,
		&book.Genero,
		&book.FechaPublicacion,
		&book.EditorialID,
		&book.Estado,
		&book.Editorial,
	)

	if err != nil {
		return nil, err
	}

	// Obtener los autores del libro
	authorsQuery := `
		SELECT idautor
		FROM libro_autor
		WHERE idlibro = ?
	`

	rows, err := m.DB.QueryContext(ctx, authorsQuery, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var authorID int
		if err := rows.Scan(&authorID); err != nil {
			return nil, err
		}
		book.Autores = append(book.Autores, authorID)
	}

	return &book, nil
}

// DeleteBook elimina un libro y sus relaciones
func (m BookModel) DeleteBook(bookID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Verificar si el libro existe
	var existingBookCount int
	checkBookQuery := `SELECT COUNT(*) FROM libro WHERE idlibro = ?`

	if err := m.DB.QueryRowContext(ctx, checkBookQuery, bookID).Scan(&existingBookCount); err != nil || existingBookCount == 0 {
		return err
	}

	// Verificar si hay préstamos activos
	var activeLoansCount int
	checkLoansQuery := `SELECT COUNT(*) FROM prestamo WHERE idlibro = ? AND fechadevolucion IS NULL`
	
	if err := m.DB.QueryRowContext(ctx, checkLoansQuery, bookID).Scan(&activeLoansCount); err != nil {
		return err
	}
	
	if activeLoansCount > 0 {
		return fmt.Errorf("no se puede eliminar el libro porque tiene préstamos activos")
	}

	// Verificar si hay reservas activas
	var activeReservationsCount int
	checkReservationsQuery := `SELECT COUNT(*) FROM reserva WHERE idlibro = ? AND estado = 'activa'`
	
	if err := m.DB.QueryRowContext(ctx, checkReservationsQuery, bookID).Scan(&activeReservationsCount); err != nil {
		return err
	}
	
	if activeReservationsCount > 0 {
		return fmt.Errorf("no se puede eliminar el libro porque tiene reservas activas")
	}

	// Eliminar multas relacionadas con préstamos de este libro
	deleteFinesQuery := `DELETE m FROM multa m 
		INNER JOIN prestamo p ON m.idprestamo = p.idprestamo 
		WHERE p.idlibro = ?`
	_, err := m.DB.ExecContext(ctx, deleteFinesQuery, bookID)
	if err != nil {
		return err
	}

	// Eliminar préstamos históricos del libro
	deleteLoansQuery := `DELETE FROM prestamo WHERE idlibro = ?`
	_, err = m.DB.ExecContext(ctx, deleteLoansQuery, bookID)
	if err != nil {
		return err
	}

	// Eliminar reservas históricas del libro
	deleteReservationsQuery := `DELETE FROM reserva WHERE idlibro = ?`
	_, err = m.DB.ExecContext(ctx, deleteReservationsQuery, bookID)
	if err != nil {
		return err
	}

	// Eliminar las relaciones con autores
	deleteAuthorsQuery := `DELETE FROM libro_autor WHERE idlibro = ?`
	_, err = m.DB.ExecContext(ctx, deleteAuthorsQuery, bookID)
	if err != nil {
		return err
	}

	// Eliminar el libro
	deleteBookQuery := `DELETE FROM libro WHERE idlibro = ?`
	_, err = m.DB.ExecContext(ctx, deleteBookQuery, bookID)
	if err != nil {
		return err
	}

	return nil
}
