package main

import (
	"encoding/json"
	"net/http"
)

// @Summary Get pending fines
// @Description Get all pending fines for admin
// @Tags Multas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} data.Fine
// @Router /api/admin/fines/to [get]
func (app *application) getPendingFinesHandler(w http.ResponseWriter, r *http.Request) {

	fine, err := app.models.Fine.GetPendingFines()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fines": fine}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get user fines
// @Description Get fines for a specific user
// @Tags Multas
// @Accept json
// @Produce json
// @Param idsocio query string true "User ID"
// @Success 200 {array} data.Fine
// @Router /api/admin/fines [get]
func (app *application) getUserFinesHandler(w http.ResponseWriter, r *http.Request) {

	idsocio := r.URL.Query().Get("idsocio")
	if idsocio == "" {
		http.Error(w, "El parámetro 'idsocio' es requerido", http.StatusBadRequest)
		return
	}

	userfines, err := app.models.Fine.GetUserFines(idsocio)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fines": userfines}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get user pending fines
// @Description Get pending fines for a specific user
// @Tags Multas
// @Accept json
// @Produce json
// @Param usuario_id query string true "User ID"
// @Success 200 {array} data.Fine
// @Router /api/fines [get]
func (app *application) getUserPendingFinesHandler(w http.ResponseWriter, r *http.Request) {

	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	fines, err := app.models.Fine.GetUserPendingFines(usuarioID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(fines) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay multas pendientes para este usuario"})
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fines": fines}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Pay fine
// @Description Mark a fine as paid
// @Tags Multas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Fine ID"
// @Success 200 {object} object "Payment confirmation"
// @Router /api/fines/{id}/pay [put]
func (app *application) payFineHandler(w http.ResponseWriter, r *http.Request) {
	fineID := r.PathValue("id")
	if fineID == "" {
		http.Error(w, "ID de la multa no proporcionado", http.StatusBadRequest)
		return
	}

	// Obtener información de la multa
	fine, err := app.models.Fine.GetFineByID(fineID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Verificar si la multa ya está pagada
	if fine.Estado == "pagada" {
		http.Error(w, "La multa ya está pagada", http.StatusBadRequest)
		return
	}

	// Marcar multa como pagada
	err = app.models.Fine.PayFine(fineID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"message": "Multa pagada exitosamente",
		"fine_id": fineID,
		"amount_paid": fine.SaldoPagar,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payment": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Search fines by user
// @Description Search fines by user name or email
// @Tags Multas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param busqueda query string true "Search term (name or email)"
// @Success 200 {array} data.Fine
// @Router /api/admin/fines/search [get]
func (app *application) searchFinesByUserHandler(w http.ResponseWriter, r *http.Request) {
	busqueda := r.URL.Query().Get("busqueda")
	if busqueda == "" {
		http.Error(w, "El parámetro 'busqueda' es requerido", http.StatusBadRequest)
		return
	}

	fines, err := app.models.Fine.SearchFinesByUser(busqueda)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"fines": fines}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
