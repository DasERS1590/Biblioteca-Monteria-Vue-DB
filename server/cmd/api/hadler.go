package main

import (
	"encoding/json"
	"net/http"
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

	defer rows.Close();

	type User struct {
		ID            int    `json:"id"`
		Nombre        string `json:"nombre"`
		Direccion     string `json:"direccion"`
		Telefono      string `json:"telefono"`
		Correo        string `json:"correo"`
		FechaNacimiento string `json:"fechanacimiento"`
		TipoSocio     string `json:"tiposocio"`
		FechaRegistro string `json:"fecharegistro"`
		ImagenPerfil  string `json:"imagenperfil"`
		Rol           string `json:"rol"`
	}

	users := make([]User ,0)

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
	// Lógica para obtener préstamos activos por usuario y rango de fechas
}

func (app *application) getPendingFinesHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener multas pendientes de pago
}

func (app *application) getUserFinesHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener el historial completo de multas por usuario
}

func (app *application) getActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener reservas activas por usuario o libro
}

func (app *application) getUserLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener historial completo de préstamos por usuario
}

func (app application) getBooksByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener libros disponibles por género y autor (Administrador)
}

func (app *application) getBooksByPublicationDateHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener libros por fecha de publicación (Administrador)
}

func (app *application) getBooksAvailableByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener libros disponibles por género y autor (Usuario)
}

func (app *application) getUserActiveLoanStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener estado de préstamos activos del usuario
}

func (app *application) getUserCompletedLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener historial de préstamos completados del usuario
}

func (app *application) getUserPendingFinesHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener multas pendientes del usuario
}

func (app *application) getUserActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener reservas activas del usuario
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para login
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para registrar usuario
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para actualizar usuario
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para eliminar usuario
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para crear un nuevo libro
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para actualizar libro
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para eliminar libro
}

func (app *application) cancelReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para cancelar reserva
}

func (app *application) extendLoanHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para extender préstamo
}

func (app *application) getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para obtener notificaciones
}
