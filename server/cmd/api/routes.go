package main

import (
	"net/http"
)

// Definir el manejador de rutas
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Rutas para Administrador

	// 1. Libros filtrados por estado y editorial
	mux.HandleFunc("/api/admin/books", app.getFilteredBooksHandler) // GET

	// 2. Libros no disponibles
	mux.HandleFunc("/api/admin/books/unavailable", app.getUnavailableBooksHandler) // GET

	// 3. Usuarios por tipo de socio
	mux.HandleFunc("/api/admin/users", app.getUsersByTypeHandler) // GET

	// 4. Préstamos activos por usuario y rango de fechas
	mux.HandleFunc("/api/admin/loans", app.getActiveLoansHandler) // GET

	// 5. Multas pendientes de pago
	mux.HandleFunc("/api/admin/fines", app.getPendingFinesHandler) // GET

	// 6. Historial completo de multas por usuario
	mux.HandleFunc("/api/admin/fines", app.getUserFinesHandler) // GET

	// 7. Reservas activas por usuario o libro
	mux.HandleFunc("/api/admin/reservations", app.getActiveReservationsHandler) // GET

	// 8. Historial de préstamos por usuario
	mux.HandleFunc("/api/admin/loans/history", app.getUserLoanHistoryHandler) // GET

	// 9. Libros disponibles por género y autor (Administrador)
	mux.HandleFunc("/api/admin/books/available", app.getBooksByGenreAndAuthorHandler) // GET

	// 10. Libros por fecha de publicación (Administrador)
	mux.HandleFunc("/api/admin/books/published", app.getBooksByPublicationDateHandler) // GET

	// Rutas para Usuario

	// 11. Libros disponibles por género y autor (Usuario)
	mux.HandleFunc("/api/books", app.getBooksAvailableByGenreAndAuthorHandler) // GET

	// 12. Estado de préstamos activos del usuario
	mux.HandleFunc("/api/loans", app.getUserActiveLoanStatusHandler) // GET

	// 13. Historial de préstamos completados del usuario
	mux.HandleFunc("/api/loans/completed", app.getUserCompletedLoanHistoryHandler) // GET

	// 14. Multas pendientes del usuario
	mux.HandleFunc("/api/fines", app.getUserPendingFinesHandler) // GET

	// 15. Reservas activas del usuario
	mux.HandleFunc("/api/reservations", app.getUserActiveReservationsHandler) // GET

	// Rutas Adicionales

	// Rutas de Autenticación
	mux.HandleFunc("/api/login", app.loginHandler) // POST - Iniciar sesión
	mux.HandleFunc("/api/register", app.registerHandler) // POST - Registrar nuevo usuario

	// Rutas de Gestión de Usuarios
	mux.HandleFunc("/api/admin/users/{id}", app.updateUserHandler) // PUT - Actualizar usuario
	mux.HandleFunc("/api/admin/users/{id}", app.deleteUserHandler) // DELETE - Eliminar usuario

	// Rutas de Gestión de Libros
	mux.HandleFunc("/api/admin/books", app.createBookHandler) // POST - Crear nuevo libro
	mux.HandleFunc("/api/admin/books/{id}", app.updateBookHandler) // PUT - Actualizar libro
	mux.HandleFunc("/api/admin/books/{id}", app.deleteBookHandler) // DELETE - Eliminar libro

	// Rutas de Gestión de Reservas
	mux.HandleFunc("/api/reservations/{id}", app.cancelReservationHandler) // DELETE - Cancelar reserva

	// Rutas de Gestión de Préstamos
	mux.HandleFunc("/api/loans/extend/{id}", app.extendLoanHandler) // POST - Extender préstamo

	// Rutas de Notificaciones
	mux.HandleFunc("/api/notifications", app.getNotificationsHandler) // GET - Obtener notificaciones

	return mux
}
