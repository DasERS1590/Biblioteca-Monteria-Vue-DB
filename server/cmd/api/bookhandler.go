package main

import (
	"biblioteca/internal/data"
	"biblioteca/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"
)

// @Summary Create a new book
// @Description Creates a new book and stores it in the database
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book body data.Book true "Book input"
// @Success 201 {object} data.Book
// @Router /api/admin/books [post]
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input data.Book

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateBook(v , &input)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	book := &data.Book{
		Titulo:           input.Titulo,
		Genero:           input.Genero,
		FechaPublicacion: input.FechaPublicacion,
		EditorialID:      input.EditorialID,
		AutoresId:        input.AutoresId,
	}

	err = app.models.Book.CreateBook(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"book": book}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Update a book
// @Description Updates an existing book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param book body data.Book true "Book input"
// @Success 200 {object} data.Book
// @Router /api/admin/books/{id} [post]
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {

	bookID := r.PathValue("id")

	if bookID == "" {
		app.notProvidedID(w, r)
		return
	}

	// Convertir bookID a int
	var id int
	_, err := fmt.Sscanf(bookID, "%d", &id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input data.Book

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	book := &data.Book{
		IdLibro:          id, // Usar el ID de la URL
		Titulo:           input.Titulo,
		Genero:           input.Genero,
		FechaPublicacion: input.FechaPublicacion,
		EditorialID:      input.EditorialID,
		AutoresId:        input.AutoresId,
	}

	err = app.models.Book.UpdateBook(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"book": book}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Get filtered books
// @Description Get books filtered by estado and editorial
// @Tags Books
// @Accept json
// @Produce json
// @Param estado query string false "Estado filter"
// @Param editorial query string false "Editorial filter"
// @Success 200 {array} data.Book
// @Router /api/admin/books [get]
func (app *application) getFilteredBooksHandler(w http.ResponseWriter, r *http.Request) {

	estado := r.URL.Query().Get("estado")
	editorial := r.URL.Query().Get("editorial")

	books, err := app.models.Book.GetFilteredBooks(estado, editorial)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get books by genre and author
// @Description Get books filtered by genre and author
// @Tags Books
// @Accept json
// @Produce json
// @Param genero query string false "Genre filter"
// @Param autor query string false "Author filter"
// @Success 200 {array} data.Book
// @Router /api/books [get]
func (app *application) getBooksByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {

	genero := r.URL.Query().Get("genero")
	autor := r.URL.Query().Get("autor")

	books, err := app.models.Book.GetBooksByGenreAndAuthor(genero, autor)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros disponibles para el criterio especificado"})
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get books by publication date
// @Description Get books published within a date range
// @Tags Books
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} data.Book
// @Router /api/books/publication-date [get]
func (app *application) getBooksByPublicationDateHandler(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		http.Error(w, "Los parámetros 'start_date' y 'end_date' son obligatorios", http.StatusBadRequest)
		return
	}

	
	books , err := app.models.Book.GetBooksByPublicationDate(startDate , endDate)
	if err != nil{
		app.serverErrorResponse(w , r , err )
	}

	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros publicados en el rango de fechas especificado"})
		return
	}

	
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Get available books by criteria
// @Description Get available books filtered by genre, author and title
// @Tags Books
// @Accept json
// @Produce json
// @Param genero query string false "Genre filter"
// @Param autor query string false "Author filter"
// @Param titulo query string false "Title filter"
// @Success 200 {array} data.Book
// @Router /api/books/available [get]
func (app *application) getBooksAvailableByGenreAndAuthorHandler(w http.ResponseWriter, r *http.Request) {


	genero := r.URL.Query().Get("genero")
	autor := r.URL.Query().Get("autor")
	titulo := r.URL.Query().Get("titulo")

	
	books , err := app.models.Book.GetBooksAvailableByGenreAndAuthor( genero , autor , titulo)
	if err != nil{
		app.serverErrorResponse(w , r , err )
	}

	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros disponibles para el criterio especificado"})
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Get books for reservation
// @Description Get books that can be reserved
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {array} data.Book
// @Router /api/books/reservation [get]
func (app *application) getBooksForReservationHandler(w http.ResponseWriter, r *http.Request) {

	genero := r.URL.Query().Get("genero")
	autor := r.URL.Query().Get("autor")
	titulo := r.URL.Query().Get("titulo")

	
	books , err := app.models.Book.GetBooksForReservation( genero , autor , titulo)
	if err != nil{
		app.serverErrorResponse(w , r , err )
	}

	if len(books) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay libros para reservar con el criterio especificado"})
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Get unavailable books
// @Description Get books that are currently unavailable
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {array} data.Book
// @Router /api/admin/books/unavailable [get]
func (app *application) getUnavailableBooksHandler(w http.ResponseWriter, r *http.Request) {

	books , err := app.models.Book.GetUnavailableBooks()
	
	w.Header().Set("Content-Type", "application/json")
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Get book for editing
// @Description Get book details for editing by ID
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} data.Book
// @Router /api/admin/books/{id}/edit [get]
func (app *application) getBookForEditHandler(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")

	if bookID == "" {
		app.notProvidedID(w, r)
		return
	}

	// Convertir bookID a int
	var id int
	_, err := fmt.Sscanf(bookID, "%d", &id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book, err := app.models.Book.GetBookForEdit(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Delete a book
// @Description Delete a book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book deleted successfully"
// @Router /api/admin/books/{id} [delete]
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("id")

	if bookID == "" {
		app.notProvidedID(w, r)
		return
	}

	// Convertir bookID a int
	var id int
	_, err := fmt.Sscanf(bookID, "%d", &id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Book.DeleteBook(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Libro eliminado exitosamente"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
