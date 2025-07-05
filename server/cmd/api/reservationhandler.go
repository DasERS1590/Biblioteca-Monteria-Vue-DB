package main

import (
	"biblioteca/internal/data"
	"encoding/json"
	"net/http"
)

// @Summary Create reservation
// @Description Create a new book reservation
// @Tags Reservas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reservation body data.Reservation true "Reservation data"
// @Success 201 {object} data.Reservation
// @Router /api/reservation [post]
func (app *application) createReservation(w http.ResponseWriter, r *http.Request) {

	var input data.Reservation

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	reservation := &data.Reservation{
		SocioID:      input.SocioID,
		LibroID:      input.LibroID,
		FechaReserva: input.FechaReserva,
	}

	err = app.models.Reservation.CreateReservation(reservation)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"reservation": reservation}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Cancel reservation
// @Description Cancel an existing reservation
// @Tags Reservas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {string} string "Reservation cancelled successfully"
// @Router /api/reservations/{id} [delete]
func (app *application) cancelReservationHandler(w http.ResponseWriter, r *http.Request) {

	reservationID := r.PathValue("id")

	if reservationID == "" {
		http.Error(w, "ID del reserva no proporcionado", http.StatusBadRequest)
		return
	}

	err := app.models.Reservation.CancelReservation(reservationID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reserva cancelada exitosamente"))
}

// @Summary Get active reservations
// @Description Get active reservations with filters
// @Tags Reservas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idsocio query string false "User ID"
// @Param idlibro query string false "Book ID"
// @Param fechareserva query string false "Reservation date"
// @Param nombre_socio query string false "User name"
// @Success 200 {array} data.Reservation
// @Router /api/admin/reservations [get]
func (app *application) getActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {
	
	idsocio := r.URL.Query().Get("idsocio")
	idlibro := r.URL.Query().Get("idlibro")
	fechareserva := r.URL.Query().Get("fechareserva")
	nombreSocio := r.URL.Query().Get("nombre_socio")

	
	reservations , err := app.models.Reservation.GetActiveReservations(idsocio , idlibro , fechareserva , nombreSocio)
	if err != nil {
		app.badRequestResponse(w , r , err )
		return 
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"reservations": reservations}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	
}

// @Summary Get user active reservations
// @Description Get active reservations for a specific user
// @Tags Reservas
// @Accept json
// @Produce json
// @Param usuario_id query string true "User ID"
// @Success 200 {array} data.Reservation
// @Router /api/reservations [get]
func (app *application) getUserActiveReservationsHandler(w http.ResponseWriter, r *http.Request) {


	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	reservations , err := app.models.Reservation.GetUserActiveReservations(usuarioID)
	if err != nil{
		app.badRequestResponse( w , r , err )
		return
	}

	if len(reservations) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay reservas activas para este usuario"})
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"reservations": reservations}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
