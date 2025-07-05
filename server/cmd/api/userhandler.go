package main

import "net/http"

// @Summary Get users by type
// @Description Get users filtered by type (normal, estudiante, profesor)
// @Tags Usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tiposocio query string true "User type (normal, estudiante, profesor)"
// @Success 200 {array} data.User_
// @Router /api/admin/users [post]
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

	usertype, err := app.models.User_.GetUserByType(userType)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"usertype": usertype}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get all users
// @Description Get all users in the system
// @Tags Usuarios
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} data.User_
// @Router /api/admin/users/all [get]
func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.User_.GetAllUsers()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
