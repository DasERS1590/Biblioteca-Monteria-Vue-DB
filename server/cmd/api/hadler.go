package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	type RootResponse struct {
		Message   string            `json:"message"`
		Version   string            `json:"version"`
		Endpoints map[string]string `json:"endpoints"`
	}

	response := RootResponse{
		Message: "Welcome to My API!",
		Version: "1.0.0",
		Endpoints: map[string]string{
			"/api/admin/books": "look a book by state",
			"/products":        "Operations related to products",
			"/orders":          "Operations related to orders",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (app *application) getFilteredBooksHandler(w http.ResponseWriter, r *http.Request) {

	estado := r.URL.Query().Get("estado")
	editorial := r.URL.Query().Get("editorial")

	query := `
		SELECT 
			libro.idlibro, libro.titulo, libro.genero, libro.estado, editorial.nombre AS editorial
		FROM 
			libro
		INNER JOIN 
			editorial ON libro.ideditorial = editorial.ideditorial
		WHERE 
			(estado = ? OR ? = '') AND 
			(editorial.nombre = ? OR ? = '')`

	rows, err := app.db.Query(query, estado, estado, editorial, editorial)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}

	type Book struct {
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Genre     string `json:"genre"`
		Status    string `json:"status"`
		Editorial string `json:"editorial"`
	}

	books := make([]Book, 0)

	for rows.Next() {
		var book Book

		if rows.Scan(&book.ID, &book.Title, &book.Genre, &book.Status, &book.Editorial) == nil {
			books = append(books, book)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (app *application) getUnavailableBooksHandler(w http.ResponseWriter, r *http.Request) {
	query := `
    SELECT 
        libro.idlibro, libro.titulo, libro.genero, libro.estado
    FROM 
        libro
    WHERE 
        libro.estado IN ('prestado', 'reservado')`

	rows, err := app.db.Query(query)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}

	type Book struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Genre  string `json:"genre"`
		Status string `json:"status"`
	}

	books := make([]Book, 0)

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Genre, &book.Status); err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(books)
}

func (app *application) getUsersByTypeHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.URL.Query().Get("tiposocio")

	if userType == "" {
		http.Error(w, "tipo de socio es requerido", http.StatusBadRequest)
		return
	}

	validUserTypes := map[string]bool{
		"normal":     true,
		"estudiante": true,
		"profesor":   true,
	}
	if !validUserTypes[userType] {
		http.Error(w, "tipo de socio no válido", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, imagenperfil, rol 
		FROM socio 
		WHERE tiposocio = ?`

	rows, err := app.db.Query(query, userType)

	if err != nil {
		http.Error(w, "error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type User struct {
		ID              int    `json:"id"`
		Nombre          string `json:"nombre"`
		Direccion       string `json:"direccion"`
		Telefono        string `json:"telefono"`
		Correo          string `json:"correo"`
		FechaNacimiento string `json:"fechanacimiento"`
		TipoSocio       string `json:"tiposocio"`
		FechaRegistro   string `json:"fecharegistro"`
		ImagenPerfil    string `json:"imagenperfil"`
		Rol             string `json:"rol"`
	}

	users := make([]User, 0)

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Nombre, &user.Direccion, &user.Telefono, &user.Correo, &user.FechaNacimiento, &user.TipoSocio, &user.FechaRegistro, &user.ImagenPerfil, &user.Rol)
		if err != nil {
			http.Error(w, "error al leer los resultados", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Responder con los usuarios encontrados
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getActiveLoansHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("idsocio")
	startDate := r.URL.Query().Get("startdate")
	endDate := r.URL.Query().Get("enddate")

	if userID == "" || startDate == "" || endDate == "" {
		http.Error(w, "idsocio, startdate, y enddate son requeridos", http.StatusBadRequest)
		return
	}

	query := `
		SELECT idprestamo, idsocio, idlibro, fechaprestamo, fechadevolucion, estado
		FROM prestamo
		WHERE idsocio = ? AND estado = 'activo' AND fechaprestamo BETWEEN ? AND ?
	`
	rows, err := app.db.Query(query, userID, startDate, endDate)
	if err != nil {
		http.Error(w, "error al obtener préstamos activos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Loans struct {
		IDPrestamo      int    `json:"idprestamo"`
		IDSocio         int    `json:"idsocio"`
		IDLibro         int    `json:"idlibro"`
		FechaPrestamo   string `json:"fechaprestamo"`
		FechaDevolucion string `json:"fechadevolucion"`
		Estado          string `json:"estado"`
	}

	loans := make([]Loans, 0)

	for rows.Next() {
		var loan Loans

		err := rows.Scan(&loan.IDPrestamo, &loan.IDSocio, &loan.IDLibro, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado)
		if err != nil {
			http.Error(w, "error al leer resultados de la base de datos", http.StatusInternalServerError)
			return
		}
		loans = append(loans, loan)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		http.Error(w, "error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getPendingFinesHandler(w http.ResponseWriter, r *http.Request) {

	query := `
		SELECT idmulta, idprestamo, saldopagar, fechamulta, estado
		FROM multa
		WHERE estado = 'pendiente'
	`
	rows, err := app.db.Query(query)
	if err != nil {
		http.Error(w, "Error al obtener multas pendientes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PendingFines struct {
		IDMulta    int     `json:"idmulta"`
		IDPrestamo int     `json:"idprestamo"`
		SaldoPagar float64 `json:"saldopagar"`
		FechaMulta string  `json:"fechamulta"`
		Estado     string  `json:"estado"`
	}

	pendingFines := make([]PendingFines, 0)

	for rows.Next() {
		var fine PendingFines
		err := rows.Scan(&fine.IDMulta, &fine.IDPrestamo, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			http.Error(w, "Error al leer los resultados de la base de datos", http.StatusInternalServerError)
			return
		}
		pendingFines = append(pendingFines, fine)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pendingFines)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getUserFinesHandler(w http.ResponseWriter, r *http.Request) {

	idsocio := r.URL.Query().Get("idsocio")
	if idsocio == "" {
		http.Error(w, "El parámetro 'idsocio' es requerido", http.StatusBadRequest)
		return
	}

	query := `
		SELECT m.idmulta, m.idprestamo, m.saldopagar, m.fechamulta, m.estado
		FROM multa m
		INNER JOIN 
			prestamo p ON m.idprestamo = p.idprestamo
		WHERE p.idsocio = ?
		ORDER BY m.fechamulta DESC
	`
	rows, err := app.db.Query(query, idsocio)
	if err != nil {
		http.Error(w, "Error al obtener el historial de multas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Fines struct {
		IDMulta    int     `json:"idmulta"`
		IDPrestamo int     `json:"idprestamo"`
		SaldoPagar float64 `json:"saldopagar"`
		FechaMulta string  `json:"fechamulta"`
		Estado     string  `json:"estado"`
	}

	fines := make([]Fines, 0)

	for rows.Next() {
		var fine Fines
		err := rows.Scan(&fine.IDMulta, &fine.IDPrestamo, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			http.Error(w, "Error al leer los resultados de la base de datos", http.StatusInternalServerError)
			return
		}
		fines = append(fines, fine)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fines)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}


func (app *application) getActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {
    type Reservation struct {
        IdReserva      int    `json:"idreserva"`
        IdSocio        int    `json:"idsocio"`
        NombreSocio    string `json:"nombre_socio"`
        IdLibro        int    `json:"idlibro"`
        TituloLibro    string `json:"titulo_libro"`
        FechaReserva   string `json:"fechareserva"`
        EstadoReserva  string `json:"estado_reserva"`
        GeneroLibro    string `json:"genero_libro"`
        Editorial      string `json:"editorial"`
        TelefonoSocio  string `json:"telefono_socio"`
        CorreoSocio    string `json:"correo_socio"`
        TipoSocio      string `json:"tiposocio"`
        FechaNacimiento string `json:"fechanacimiento"`
        FechaRegistro  string `json:"fecharegistro"`
    }

    // Obtener parámetros de búsqueda
    usuario := r.URL.Query().Get("usuarioid")
    libro := r.URL.Query().Get("libro")
    fecha := r.URL.Query().Get("fecha")
    nombreSocio := r.URL.Query().Get("nombre")

    // Crear la consulta SQL base
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

    // Lista para los parámetros
    var params []interface{}

    // Agregar filtros condicionalmente
    if usuario != "" {
        query += " AND r.idsocio = ?"
        params = append(params, usuario)
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

    // Ordenar por fecha de reserva de manera descendente
    query += " ORDER BY r.fechareserva DESC"

    // Ejecutar la consulta con los parámetros
    rows, err := app.db.Query(query, params...)
    if err != nil {
        http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Leer los resultados
    var reservations []Reservation
    for rows.Next() {
        var reservation Reservation
        if err := rows.Scan(&reservation.IdReserva, &reservation.IdSocio, &reservation.NombreSocio,
            &reservation.IdLibro, &reservation.TituloLibro, &reservation.FechaReserva, &reservation.EstadoReserva,
            &reservation.GeneroLibro, &reservation.Editorial, &reservation.TelefonoSocio, &reservation.CorreoSocio,
            &reservation.TipoSocio, &reservation.FechaNacimiento, &reservation.FechaRegistro); err != nil {
            http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
            return
        }
        reservations = append(reservations, reservation)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
        return
    }

    // Establecer el encabezado de la respuesta como JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(reservations); err != nil {
        http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
    }
}





func (app *application) getUserLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario := r.URL.Query().Get("idsocio")

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

	rows, err := app.db.Query(query, idUsuario)

	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Loan struct {
		IDPrestamo      int    `json:"id_prestamo"`
		IDUsuario       int    `json:"id_usuario"`
		IDLibro         int    `json:"id_libro"`
		TituloLibro     string `json:"titulo_libro"`
		FechaPrestamo   string `json:"fecha_prestamo"`
		FechaDevolucion string `json:"fecha_devolucion"`
		Estado          string `json:"estado"`
	}

	loans := make([]Loan, 0)

	for rows.Next() {
		var loan Loan
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
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}
	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos registrados para este usuario"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app application) getBooksByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	genero := r.URL.Query().Get("genero")
	autor := r.URL.Query().Get("autor")

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

	rows, err := app.db.Query(query, genero, "%"+autor+"%")
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Book struct {
		IDLibro int    `json:"id_libro"`
		Titulo  string `json:"titulo"`
		Genero  string `json:"genero"`
		Estado  string `json:"estado"`
		Autor   string `json:"autor"`
	}

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Autor)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros disponibles para el criterio especificado"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getBooksByPublicationDateHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener los parámetros de la URL
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Validar que ambos parámetros estén presentes
	if startDate == "" || endDate == "" {
		http.Error(w, "Los parámetros 'start_date' y 'end_date' son obligatorios", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL
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

	// Ejecutar la consulta
	rows, err := app.db.Query(query, startDate, endDate)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para los libros
	type Book struct {
		IDLibro          int    `json:"id_libro"`
		Titulo           string `json:"titulo"`
		Genero           string `json:"genero"`
		FechaPublicacion string `json:"fecha_publicacion"`
		Estado           string `json:"estado"`
		Editorial        string `json:"editorial"`
	}

	var books []Book

	// Procesar resultados
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.FechaPublicacion, &book.Estado, &book.Editorial)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar mensaje claro
	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros publicados en el rango de fechas especificado"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}

}

func (app *application) getBooksAvailableByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener los parámetros de la URL
	genero := r.URL.Query().Get("genero")
	autor := r.URL.Query().Get("autor")

	// Validar que los parámetros sean proporcionados
	if genero == "" || autor == "" {
		http.Error(w, "Los parámetros 'genero' y 'autor' son obligatorios", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL
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
			libro.estado = 'disponible' AND libro.genero = ? AND autor.nombre LIKE ?
	`

	// Ejecutar la consulta
	rows, err := app.db.Query(query, genero, "%"+autor+"%")
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para los libros
	type Book struct {
		IDLibro int    `json:"id_libro"`
		Titulo  string `json:"titulo"`
		Genero  string `json:"genero"`
		Estado  string `json:"estado"`
		Autor   string `json:"autor"`
	}

	var books []Book

	// Procesar resultados
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.IDLibro, &book.Titulo, &book.Genero, &book.Estado, &book.Autor)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar mensaje claro
	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros disponibles para el criterio especificado"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getUserActiveLoanStatusHandler(w http.ResponseWriter, r *http.Request) {

	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro del usuario desde la URL
	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL
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

	// Ejecutar la consulta
	rows, err := app.db.Query(query, usuarioID)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para los préstamos
	type Loan struct {
		IDPrestamo      int    `json:"id_prestamo"`
		FechaPrestamo   string `json:"fecha_prestamo"`
		FechaDevolucion string `json:"fecha_devolucion"`
		Estado          string `json:"estado"`
		TituloLibro     string `json:"titulo_libro"`
	}

	var loans []Loan

	// Procesar resultados
	for rows.Next() {
		var loan Loan
		err := rows.Scan(&loan.IDPrestamo, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado, &loan.TituloLibro)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar un mensaje claro
	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos activos para este usuario"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}

}

func (app *application) getUserCompletedLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro del usuario desde la URL
	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL
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

	// Ejecutar la consulta
	rows, err := app.db.Query(query, usuarioID)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para los préstamos
	type Loan struct {
		IDPrestamo      int    `json:"id_prestamo"`
		FechaPrestamo   string `json:"fecha_prestamo"`
		FechaDevolucion string `json:"fecha_devolucion"`
		Estado          string `json:"estado"`
		TituloLibro     string `json:"titulo_libro"`
	}

	var loans []Loan

	// Procesar resultados
	for rows.Next() {
		var loan Loan
		err := rows.Scan(&loan.IDPrestamo, &loan.FechaPrestamo, &loan.FechaDevolucion, &loan.Estado, &loan.TituloLibro)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar un mensaje claro
	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos completados para este usuario"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getUserPendingFinesHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro del usuario desde la URL
	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL para obtener las multas pendientes
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

	// Ejecutar la consulta
	rows, err := app.db.Query(query, usuarioID)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para las multas
	type Fine struct {
		IDMulta    int     `json:"id_multa"`
		SaldoPagar float64 `json:"saldo_pagar"`
		FechaMulta string  `json:"fecha_multa"`
		Estado     string  `json:"estado"`
	}

	var fines []Fine

	// Procesar resultados
	for rows.Next() {
		var fine Fine
		err := rows.Scan(&fine.IDMulta, &fine.SaldoPagar, &fine.FechaMulta, &fine.Estado)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		fines = append(fines, fine)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar mensaje claro
	if len(fines) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay multas pendientes para este usuario"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fines)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) getUserActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {
	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el parámetro del usuario desde la URL
	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	// Construir consulta SQL para obtener las reservas activas del usuario
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

	// Ejecutar la consulta
	rows, err := app.db.Query(query, usuarioID)
	if err != nil {
		http.Error(w, "Error ejecutando consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Definir la estructura para las reservas
	type Reservation struct {
		IDReserva    int    `json:"id_reserva"`
		IDSocio      int    `json:"id_socio"`
		IDLibro      int    `json:"id_libro"`
		FechaReserva string `json:"fecha_reserva"`
		Estado       string `json:"estado"`
	}

	var reservations []Reservation

	// Procesar resultados
	for rows.Next() {
		var reservation Reservation
		err := rows.Scan(&reservation.IDReserva, &reservation.IDSocio, &reservation.IDLibro, &reservation.FechaReserva, &reservation.Estado)
		if err != nil {
			http.Error(w, "Error al leer los resultados", http.StatusInternalServerError)
			return
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error durante la iteración de filas", http.StatusInternalServerError)
		return
	}

	// Si no hay resultados, enviar mensaje claro
	if len(reservations) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay reservas activas para este usuario"})
		return
	}

	// Responder con los resultados en formato JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reservations)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	// Estructura para la solicitud de inicio de sesión
	type Credentials struct {
		Correo     string `json:"correo"`
		Contrasena string `json:"contrasena"`
	}

	// Estructura para el usuario
	type Usuario struct {
		ID       int    `json:"id"`
		Nombre   string `json:"nombre"`
		Rol      string `json:"rol"`
		Password string `json:"password"`
	}

	// Asegurarse de que el método sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Consulta para obtener la información del usuario
	query := `
		SELECT socio.idsocio, socio.nombre, socio.rol, usuariopassword.hash_contrasena
		FROM socio
		JOIN usuariopassword ON socio.idsocio = usuariopassword.idusuario
		WHERE socio.correo = ?
	`

	var usuario Usuario
	err = app.db.QueryRow(query, credentials.Correo).Scan(&usuario.ID, &usuario.Nombre, &usuario.Rol, &usuario.Password)
	fmt.Println(usuario)
	if err != nil {
		if err == sql.ErrNoRows {

			http.Error(w, "Correo o contraseña incorrectos", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error en la consulta de usuario", http.StatusInternalServerError)
		}
		return
	}

	// Función para encriptar la contraseña usando SHA-256
	hashPassword := func(password string) string {
		// Generar el hash de la contraseña utilizando SHA-256
		hash := sha256.New()
		hash.Write([]byte(password))
		return fmt.Sprintf("%x", hash.Sum(nil)) // Retornar el hash como cadena hex
	}

	// Verificar si la contraseña es correcta usando SHA-256
	hashedPassword := hashPassword(credentials.Contrasena)
	if hashedPassword != usuario.Password {
		//http.Error(w, "Correo o contraseña incorrectos", http.StatusUnauthorized)
		//return
	}

	// Responder con el usuario autenticado y su rol
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     usuario.ID,
		"nombre": usuario.Nombre,
		"rol":    usuario.Rol,
	})
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Validar el método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Estructura para el registro de usuario
	type RegisterRequest struct {
		Nombre          string `json:"nombre"`
		Direccion       string `json:"direccion"`
		Telefono        string `json:"telefono"`
		Correo          string `json:"correo"`
		FechaNacimiento string `json:"fecha_nacimiento"` // Formato: YYYY-MM-DD
		TipoSocio       string `json:"tipo_socio"`
		Contrasena      string `json:"contrasena"`
		Rol             string `json:"rol"` // Usuario o administrador
	}

	// Decodificar el cuerpo JSON
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if req.Nombre == "" || req.Direccion == "" || req.Telefono == "" || req.Correo == "" ||
		req.FechaNacimiento == "" || req.TipoSocio == "" || req.Contrasena == "" || req.Rol == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Validar el rol
	if req.Rol != "usuario" && req.Rol != "administrador" {
		http.Error(w, "El campo 'rol' debe ser 'usuario' o 'administrador'", http.StatusBadRequest)
		return
	}

	// Validar formato de la fecha de nacimiento
	_, err = time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		http.Error(w, "El formato de la fecha de nacimiento debe ser YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Verificar que el correo no esté registrado
	var existingID int
	queryCheck := "SELECT idsocio FROM socio WHERE correo = ?"
	err = app.db.QueryRow(queryCheck, req.Correo).Scan(&existingID)
	if err == nil {
		http.Error(w, "El correo ya está registrado", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Error al verificar el correo", http.StatusInternalServerError)
		return
	}

	var lastID int
	queryGetLastID := "SELECT MAX(idsocio) FROM socio"
	err = app.db.QueryRow(queryGetLastID).Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error al obtener el último ID", http.StatusInternalServerError)
		return
	}

	// Calcular el nuevo ID, incrementando el último ID por 1
	newID := lastID + 1

	// Insertar el nuevo usuario en la tabla `socio` con el ID manual
	queryInsertSocio := `
		INSERT INTO socio (idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, imagenperfil, rol)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?)
	`
	_, err = app.db.Exec(queryInsertSocio, newID, req.Nombre, req.Direccion, req.Telefono, req.Correo, req.FechaNacimiento, req.TipoSocio, "NULL", req.Rol)

	if err != nil {
		http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// El nuevo ID ya está asignado
	userID := newID

	// Función para encriptar la contraseña usando SHA-256
	hashPassword := func(password string) string {
		hash := sha256.New()
		hash.Write([]byte(password))
		return fmt.Sprintf("%x", hash.Sum(nil))
	}

	// Insertar la contraseña en la tabla `usuariopassword`
	hashedPassword := hashPassword(req.Contrasena)
	queryInsertPassword := `
	INSERT INTO usuariopassword (idusuario, hash_contrasena)
	VALUES (?, ?)
`
	_, err = app.db.Exec(queryInsertPassword, userID, hashedPassword)
	if err != nil {
		http.Error(w, "Error al guardar la contraseña del usuario", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuario registrado exitosamente",
		"id":      userID,
	})

}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook struct {
		IdLibro          int    `json:"idlibro"`
		Titulo           string `json:"titulo"`
		Genero           string `json:"genero"`
		FechaPublicacion string `json:"fechapublicacion"`
		EditorialID      int    `json:"ideditorial"`
		Autores          []int  `json:"autores"` // Lista de IDs de autores
	}

	// Decodificar el cuerpo del request JSON
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	// Genera un nuevo ID único para el libro
	newBook.IdLibro = generateNewId(app, "libro", "idlibro")

	// Inserción del nuevo libro
	insertBookQuery := `
        INSERT INTO libro (idlibro, titulo, genero, fechapublicacion, estado, ideditorial)
        VALUES (?, ?, ?, ?, 'disponible', ?)`

	_, err = app.db.Exec(insertBookQuery, newBook.IdLibro, newBook.Titulo, newBook.Genero, newBook.FechaPublicacion, newBook.EditorialID)
	if err != nil {
		http.Error(w, "Error al insertar el libro", http.StatusInternalServerError)
		return
	}

	// Insertar la relación de los autores con el libro
	for _, autorID := range newBook.Autores {
		insertAuthorQuery := `
            INSERT INTO libro_autor (idlibro, idautor)
            VALUES (?, ?)`

		_, err = app.db.Exec(insertAuthorQuery, newBook.IdLibro, autorID)
		if err != nil {
			http.Error(w, "Error al asociar el autor al libro", http.StatusInternalServerError)
			return
		}
	}

	// Retornar respuesta
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Libro creado exitosamente"))
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del libro de la URL
	bookID := strings.TrimPrefix(r.URL.Path, "/api/admin/books/")

	if bookID == "" {
		http.Error(w, "ID del libro no proporcionado", http.StatusBadRequest)
		return
	}

	var updatedBook struct {
		Titulo           string `json:"titulo"`
		Genero           string `json:"genero"`
		FechaPublicacion string `json:"fechapublicacion"`
		EditorialID      int    `json:"ideditorial"`
		Autores          []int  `json:"autores"` // Lista de IDs de autores
	}

	// Decodificar el cuerpo del request JSON
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Asegurarse de cerrar el cuerpo del request

	// Verificar si el libro existe
	var existingBookCount int
	checkBookQuery := `SELECT COUNT(*) FROM libro WHERE idlibro = ?`
	if err := app.db.QueryRow(checkBookQuery, bookID).Scan(&existingBookCount); err != nil || existingBookCount == 0 {
		http.Error(w, "Libro no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar los datos del libro
	updateBookQuery := `
			UPDATE libro
			SET titulo = ?, genero = ?, fechapublicacion = ?, ideditorial = ?
			WHERE idlibro = ?`

	if _, err := app.db.Exec(updateBookQuery, updatedBook.Titulo, updatedBook.Genero, updatedBook.FechaPublicacion, updatedBook.EditorialID, bookID); err != nil {
		http.Error(w, "Error al actualizar el libro", http.StatusInternalServerError)
		return
	}

	// Eliminar asociaciones anteriores de autores
	deleteAuthorsQuery := `DELETE FROM libro_autor WHERE idlibro = ?`
	if _, err := app.db.Exec(deleteAuthorsQuery, bookID); err != nil {
		http.Error(w, "Error al eliminar las asociaciones de autores", http.StatusInternalServerError)
		return
	}

	// Insertar las nuevas asociaciones de autores
	for _, autorID := range updatedBook.Autores {
		insertAuthorQuery := `
				INSERT INTO libro_autor (idlibro, idautor)
				VALUES (?, ?)`

		if _, err := app.db.Exec(insertAuthorQuery, bookID, autorID); err != nil {
			http.Error(w, "Error al asociar el autor al libro", http.StatusInternalServerError)
			return
		}
	}

	// Retornar respuesta exitosa
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Libro actualizado exitosamente"))
}

func (app *application) createReservation(w http.ResponseWriter, r *http.Request) {
	var newReservation struct {
		SocioID      int    `json:"idsocio"`
		LibroID      int    `json:"idlibro"`
		FechaReserva string `json:"fechareserva"`
	}

	err := json.NewDecoder(r.Body).Decode(&newReservation)
	if err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	newReservationID := generateNewId(app, "reserva", "idreserva")

	insertReservationQuery := `
        INSERT INTO reserva (idreserva, idsocio, idlibro, fechareserva, estado)
        VALUES (?, ?, ?, ?, 'activa')`

	_, err = app.db.Exec(insertReservationQuery, newReservationID, newReservation.SocioID, newReservation.LibroID, newReservation.FechaReserva)
	if err != nil {
		http.Error(w, "Error al crear la reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Reserva creada exitosamente"))
}

func (app *application) cancelReservationHandler(w http.ResponseWriter, r *http.Request) {

	reservationID := strings.TrimPrefix(r.URL.Path, "/api/reservations/")

	if reservationID == "" {
		http.Error(w, "ID del reserva no proporcionado", http.StatusBadRequest)
		return
	}

	// Consulta para verificar si la reserva existe y está activa
	var estado string
	err := app.db.QueryRow("SELECT estado FROM reserva WHERE idreserva = ?", reservationID).Scan(&estado)
	if err != nil {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		fmt.Println(err)
		return
	}

	if estado != "activa" {
		http.Error(w, "No se puede cancelar una reserva que no está activa", http.StatusBadRequest)
		return
	}

	// Eliminar la reserva (cambiar su estado a 'cancelada' o eliminarla completamente)
	updateQuery := "UPDATE reserva SET estado = 'cancelada' WHERE idreserva = ?"
	_, err = app.db.Exec(updateQuery, reservationID)
	if err != nil {
		http.Error(w, "Error al cancelar la reserva", http.StatusInternalServerError)
		return
	}

	// Retornar respuesta exitosa
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reserva cancelada exitosamente"))
}

// Función para extender préstamo
func (app *application) extendLoanHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del préstamo de la URL
	reservationID := strings.TrimPrefix(r.URL.Path, "/api/loans/extend/")
	if reservationID == "" {
		http.Error(w, "ID del reserva no proporcionado", http.StatusBadRequest)
		return
	}

	// Decodificar la nueva fecha de devolución desde el cuerpo de la solicitud
	var request struct {
		NuevaFechaDevolucion string `json:"nuevafechadevolucion"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	// Verificar si la fecha proporcionada es válida
	nuevaFecha, err := time.Parse("2006-01-02", request.NuevaFechaDevolucion)
	if err != nil {
		http.Error(w, "Fecha inválida, el formato debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	// Consulta para verificar el préstamo
	var estado string
	var fechaDevolucion time.Time
	err = app.db.QueryRow("SELECT estado, fechadevolucion FROM prestamo WHERE idprestamo = ?", reservationID).Scan(&estado, &fechaDevolucion)
	if err != nil {
		http.Error(w, "Préstamo no encontrado", http.StatusNotFound)

		fmt.Println(err)
		return
	}

	// Verificar si el préstamo ya ha sido completado
	if estado == "completado" {
		http.Error(w, "El préstamo ya está completado, no se puede extender", http.StatusBadRequest)
		return
	}

	// Actualizar la fecha de devolución en la base de datos
	updateQuery := "UPDATE prestamo SET fechadevolucion = ? WHERE idprestamo = ?"
	_, err = app.db.Exec(updateQuery, nuevaFecha, reservationID)
	if err != nil {
		http.Error(w, "Error al extender el préstamo", http.StatusInternalServerError)
		return
	}

	// Retornar respuesta exitosa
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Préstamo extendido exitosamente"))
}

func generateNewId(app *application, tableName string, idColumn string) int {
	var maxId int

	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", idColumn, tableName)

	row := app.db.QueryRow(query)
	err := row.Scan(&maxId)
	if err != nil {
		maxId = 0
	}

	return maxId + 1
}
