package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("GET /api/admin/books", app.getFilteredBooksHandler) 
	mux.HandleFunc("GET /api/admin/books/unavailable", app.getUnavailableBooksHandler)
	mux.HandleFunc("GET /api/admin/users", app.getUsersByTypeHandler) 
	mux.HandleFunc("GET /api/admin/loans", app.getActiveLoansHandler) // GET

	// 5. Multas pendientes de pago
	mux.HandleFunc("GET /api/admin/fines/to", app.getPendingFinesHandler) // GET

	// 6. Historial completo de multas por usuario
	mux.HandleFunc("GET /api/admin/fines", app.getUserFinesHandler) // GET

	// 7. Reservas activas por usuario o libro
	mux.HandleFunc("GET /api/admin/reservations", app.getActiveReservationsHandler) // GET

	// 8. Historial de préstamos por usuario
	mux.HandleFunc("GET /api/admin/loans/history", app.getUserLoanHistoryHandler) // GET

	// 9. Libros disponibles por género y autor (Administrador)
	mux.HandleFunc("GET /api/admin/books/available", app.getBooksByGenreAndAuthorHandler) // GET

	// 10. Libros por fecha de publicación (Administrador)
	mux.HandleFunc("GET /api/admin/books/published", app.getBooksByPublicationDateHandler) // GET

	// Rutas para Usuario

	// 11. Libros disponibles por género y autor (Usuario)
	mux.HandleFunc("GET /api/books", app.getBooksAvailableByGenreAndAuthorHandler) // GET

	// 12. Estado de préstamos activos del usuario
	mux.HandleFunc("GET /api/loans", app.getUserActiveLoanStatusHandler) // GET

	// 13. Historial de préstamos completados del usuario
	mux.HandleFunc("GET /api/loans/completed", app.getUserCompletedLoanHistoryHandler) // GET

	// 14. Multas pendientes del usuario
	mux.HandleFunc("GET /api/fines", app.getUserPendingFinesHandler) // GET

	// 15. Reservas activas del usuario
	mux.HandleFunc("GET /api/reservations", app.getUserActiveReservationsHandler) // GET

	// Rutas Adicionales

	// Rutas de Autenticación
	mux.HandleFunc("POST /api/login", app.loginHandler)       // POST - Iniciar sesión
	mux.HandleFunc("POST /api/register", app.registerHandler) // POST - Registrar nuevo usuario

	// Rutas de Gestión de Usuarios
	mux.HandleFunc("PUT /api/admin/users/{id}", app.updateUserHandler)    // PUT - Actualizar usuario
	mux.HandleFunc("DELETE /api/admin/users/{id}", app.deleteUserHandler) // DELETE - Eliminar usuario

	// Rutas de Gestión de Libros
	mux.HandleFunc("POST /api/admin/books", app.createBookHandler)       // POST - Crear nuevo libro
	mux.HandleFunc("PUT /api/admin/books/{id}", app.updateBookHandler)   // PUT - Actualizar libro
	mux.HandleFunc("DELTE /api/admin/books/{id}", app.deleteBookHandler) // DELETE - Eliminar libro

	// Rutas de Gestión de Reservas
	mux.HandleFunc("DELETE /api/reservations/{id}", app.cancelReservationHandler) // DELETE - Cancelar reserva

	// Rutas de Gestión de Préstamos
	mux.HandleFunc("POST /api/loans/extend/{id}", app.extendLoanHandler) // POST - Extender préstamo

	// Rutas de Notificaciones
	mux.HandleFunc("GET /api/notifications", app.getNotificationsHandler) // GET - Obtener notificaciones

	return mux
}
