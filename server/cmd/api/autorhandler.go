package main

import (
	"biblioteca/internal/data"
	"net/http"
)

// getAutoresHandler godoc
// @Summary Get all authors
// @Description Returns all authors
// @Tags Authors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} data.Author
// @Router /api/autores [get]
func (app *application) getAutoresHandler(w http.ResponseWriter, r *http.Request) {

	authors , err  := app.models.Autor.GetAutores()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"authors": authors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getAutoresHandler godoc
// @Summary Create author
// @Description Create a new author
// @Tags Authors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body data.NewAuthor true "Create an author"
// @Success 201 {object} data.Author
// @Router /api/admin/autores [post]
func (app *application) createAutorHander(w http.ResponseWriter, r *http.Request) {
	var input data.NewAuthor
	
	err :=  app.readJSON(w , r , &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	author := &data.Author{
		Nombre: input.Name,
		Nacionalidad: input.Nationality,
	}

	err = app.models.Autor.CreateAuthor(author)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"author": author}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
