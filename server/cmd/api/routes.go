package main 

import (
    "net/http"
)
 
func (app *application) routes() *http.ServeMux {
    
    mux := http.NewServeMux()
    
    mux.HandleFunc("POST /api/admin/books", app.createBookHandler)
    mux.HandleFunc("PUT /api/admin/books/{id}", app.updateBookHandler)
    mux.HandleFunc("GET /api/admin/books" , app.getBookHandler) 
    mux.HandleFunc("GET /api/editorials", app.checkEditorial ) 
    mux.HandleFunc("POST /api/editorials", app.createEditorial )

    
    mux.HandleFunc("POST /api/admin/users", app.registerUserHandler)
    mux.HandleFunc("PUT /api/admin/users/{id}" , app.updateUserHandler)
    mux.HandleFunc("GET /api/admin/users", app.getUserHandler)

    mux.HandleFunc("POST /api/admin/loans", app.createLoanHandler)
    mux.HandleFunc("PUT /api/admin/loans/{id}", app.updateLoanHandler)
    mux.HandleFunc("GET /api/admin/loans", app.getLoanHandler)

    mux.HandleFunc("POST /api/admin/fines", app.createFineHandler)
    mux.HandleFunc("GET /api/admin/fines",  app.getFineHandler)

    mux.HandleFunc("POST /api/admin/reservations", app.createReservationHandler)
    mux.HandleFunc("PUT /api/admin/reservarions/{id}", app.updateReservationHandler)
    mux.HandleFunc("GET /api/admin/reservations", app.getReservationHandler)
    
    mux.HandleFunc("GET /api/books", app.getBookHandler )
    
    mux.HandleFunc("POST /api/loans", app.createBookHandler)
    mux.HandleFunc("GET /api/loans", app.getLoanHandler)
    
    mux.HandleFunc("GET /api/fines", app.getFineHandler )
    
    mux.HandleFunc("POST /api/reservations", app.createReservationHandler)
    mux.HandleFunc("GET /api/reservations", app.getReservationHandler)
    mux.HandleFunc("PUT /api/reservations/{id}", app. updateReservationHandler)
        
    mux.HandleFunc("POST /api/login", app.loginHandler)

    return mux 
}
