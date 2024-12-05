package main

import (
	"net/http"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("GET /api/admin/books", app.getFilteredBooksHandler) 
	mux.HandleFunc("GET /api/admin/books/unavailable", app.getUnavailableBooksHandler)
	mux.HandleFunc("GET /api/admin/users", app.getUsersByTypeHandler) 
	mux.HandleFunc("GET /api/admin/loans", app.getActiveLoansHandler) 
	mux.HandleFunc("GET /api/admin/fines/to", app.getPendingFinesHandler) 
	mux.HandleFunc("GET /api/admin/fines", app.getUserFinesHandler)
	mux.HandleFunc("GET /api/admin/reservations", app.getActiveReservationsHandler) 
	mux.HandleFunc("GET /api/admin/loans/history", app.getUserLoanHistoryHandler) 
	mux.HandleFunc("GET /api/admin/books/available", app.getBooksByGenreAndAuthorHandler) 
	mux.HandleFunc("GET /api/admin/books/published", app.getBooksByPublicationDateHandler) 

	// Rutas para Usuario
	mux.HandleFunc("GET /api/books", app.getBooksAvailableByGenreAndAuthorHandler) 
	mux.HandleFunc("GET /api/loans", app.getUserActiveLoanStatusHandler) 
	mux.HandleFunc("GET /api/loans/completed", app.getUserCompletedLoanHistoryHandler) 
	mux.HandleFunc("GET /api/fines", app.getUserPendingFinesHandler)
	mux.HandleFunc("GET /api/reservations", app.getUserActiveReservationsHandler) // GET


	mux.HandleFunc("POST /api/login", app.loginHandler)     
	mux.HandleFunc("POST /api/register", app.registerHandler) 

	// Rutas de Gestión de Usuarios
	//mux.HandleFunc("PUT /api/admin/users/{id}", app.updateUserHandler)    // PUT - Actualizar usuario
	//mux.HandleFunc("DELETE /api/admin/users/{id}", app.deleteUserHandler) // DELETE - Eliminar usuario

	// Rutas de Gestión de Libros
	mux.HandleFunc("POST /api/admin/books", app.createBookHandler)       
	mux.HandleFunc("PUT /api/admin/books/{id}", app.updateBookHandler)   // PUT - Actualizar libro
	mux.HandleFunc("DELTE /api/admin/books/{id}", app.deleteBookHandler) // DELETE - Eliminar libro

	// Rutas de Gestión de Reservas
	mux.HandleFunc("DELETE /api/reservations/{id}", app.cancelReservationHandler) // DELETE - Cancelar reserva

	// Rutas de Gestión de Préstamos
	mux.HandleFunc("POST /api/loans/extend/{id}", app.extendLoanHandler) // POST - Extender préstamo

	// Rutas de Notificaciones
	//mux.HandleFunc("GET /api/notifications", app.getNotificationsHandler) // GET - Obtener notificaciones


	c := cors.New(cors.Options{
        //AllowedOrigins:   []string{"http://127.0.0.1:3000"}, // Dominios permitidos
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

	
	return c.Handler(mux)
}
