package main

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/docs/", httpSwagger.WrapHandler)

	// Rutas de estado
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	// Autenticación
	mux.HandleFunc("POST /v1/api/login", app.loginHandler)
	mux.HandleFunc("POST /v1/api/register", app.registerHandler)

	// Usuario - Libros disponibles
	mux.HandleFunc("GET  /v1/api/books", app.requirePermission(PermissionBooksRead, app.getBooksAvailableByGenreAndAuthorHandler))
	mux.HandleFunc("GET  /v1/api/books/reservation", app.requirePermission(PermissionBooksRead, app.getBooksForReservationHandler))

	// Admin - Libros
	mux.HandleFunc("GET  /v1/api/admin/books", app.requirePermission(PermissionBooksRead, app.getFilteredBooksHandler))
	mux.HandleFunc("GET  /v1/api/admin/books/{id}/edit", app.requirePermission(PermissionBooksRead, app.getBookForEditHandler))
	mux.HandleFunc("GET  /v1/api/admin/books/unavailable", app.requirePermission(PermissionBooksRead, app.getUnavailableBooksHandler))
	mux.HandleFunc("GET  /v1/api/admin/books/available", app.requirePermission(PermissionBooksRead, app.getBooksByGenreAndAuthorHandler))
	mux.HandleFunc("GET  /v1/api/admin/books/published", app.requirePermission(PermissionBooksRead, app.getBooksByPublicationDateHandler))
	mux.HandleFunc("POST /v1/api/admin/books", app.requirePermission(PermissionBooksWrite, app.createBookHandler))
	mux.HandleFunc("POST /v1/api/admin/books/{id}", app.requirePermission(PermissionBooksWrite, app.updateBookHandler))
	mux.HandleFunc("DELETE /v1/api/admin/books/{id}", app.requirePermission(PermissionBooksWrite, app.deleteBookHandler))

	// Admin - Usuarios, préstamos, multas, reservas
	mux.HandleFunc("POST /v1/api/admin/users", app.requirePermission(PermissionUsersManage, app.getUsersByTypeHandler))
	mux.HandleFunc("GET /v1/api/admin/users/all", app.requirePermission(PermissionUsersManage, app.getAllUsersHandler))
	mux.HandleFunc("POST /v1/api/admin/loans", app.requirePermission(PermissionLoansManage, app.getActiveLoansHandler))
	mux.HandleFunc("GET  /v1/api/admin/fines/to", app.requirePermission(PermissionFinesRead, app.getPendingFinesHandler))
	mux.HandleFunc("GET  /v1/api/admin/fines", app.requirePermission(PermissionFinesRead, app.getUserFinesHandler))
	mux.HandleFunc("GET  /v1/api/admin/fines/search", app.requirePermission(PermissionFinesRead, app.searchFinesByUserHandler))
	mux.HandleFunc("GET  /v1/api/admin/reservations", app.requirePermission(PermissionReservationsView, app.getActiveReservationsHandler))
	mux.HandleFunc("GET  /v1/api/admin/loans/history", app.requirePermission(PermissionLoansView, app.getUserLoanHistoryHandler))

	// Usuario - Préstamos y multas
	mux.HandleFunc("POST /v1/api/loans", app.requirePermission(PermissionLoansCreate, app.createLoanHandler))
	mux.HandleFunc("GET  /v1/api/loans", app.requirePermission(PermissionLoansView, app.getUserActiveLoanStatusHandler))
	mux.HandleFunc("GET  /v1/api/loans/completed", app.requirePermission(PermissionLoansView, app.getUserCompletedLoanHistoryHandler))
	mux.HandleFunc("POST   /v1/api/loans/return/{id}", app.requirePermission(PermissionLoansCreate, app.returnLoanHandler))
	mux.HandleFunc("GET  /v1/api/fines", app.requirePermission(PermissionFinesRead, app.getUserPendingFinesHandler))
	mux.HandleFunc("PUT  /v1/api/fines/{id}/pay", app.requirePermission(PermissionFinesCreate, app.payFineHandler))
	mux.HandleFunc("GET  /v1/api/reservations", app.requirePermission(PermissionReservationsView, app.getUserActiveReservationsHandler))

	// Editoriales y autores
	mux.HandleFunc("POST /v1/api/editoriales", app.requirePermission(PublishersCreate, app.createEditorialHandler))
	mux.HandleFunc("GET  /v1/api/editoriales", app.requirePermission(PublishersRead, app.getEditorialsHandler))
	mux.HandleFunc("POST /v1/api/admin/autores", app.requirePermission(AuthorsCreate, app.createAutorHander))
	mux.HandleFunc("GET  /v1/api/autores", app.requirePermission(AauthorsRead, app.getAutoresHandler))

	// Reservas y extensiones
	mux.HandleFunc("POST   /v1/api/reservation", app.requirePermission(PermissionReservationsCreate, app.createReservation))
	mux.HandleFunc("DELETE /v1/api/reservations/{id}", app.requirePermission(PermissionReservationsCreate, app.cancelReservationHandler))
	mux.HandleFunc("POST   /v1/api/loans/extend/{id}", app.requirePermission(PermissionReservationsCreate, app.extendLoanHandler))

	// Encadenar middlewares
	return app.recoverPanic(
		app.enableCORS(
			app.rateLimit(
				app.authenticate(mux),
			),
		),
	)
}
