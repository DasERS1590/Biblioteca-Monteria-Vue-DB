package main

import (
	"biblioteca/internal/data"
	"encoding/json"
	"net/http"
	"time"
)

// parseDate intenta parsear una fecha en múltiples formatos
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02T15:04:05Z", // ISO con Z
		"2006-01-02T15:04:05",  // ISO sin Z
		"2006-01-02 15:04:05",  // MySQL format
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, &time.ParseError{
		Layout: "multiple formats",
		Value:  dateStr,
		Message: "could not parse date in any supported format",
	}
}

// createLoanHandler maneja las solicitudes POST para crear un nuevo préstamo.
//
// @Summary Crear un nuevo préstamo
// @Description Crea un nuevo préstamo para un usuario autenticado.
// @Tags Prestamos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param loan body data.Loan true "Crear prestamo"
// @Success 201 {object} map[string]interface{} "Préstamo creado exitosamente"
// @Failure 400 {string} string "Solicitud mal formada o parámetros faltantes"
// @Failure 401 {string} string "No autorizado: token JWT inválido o ausente"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/loans [post]
func (app *application) createLoanHandler(w http.ResponseWriter, r *http.Request) {

	var input data.Loan

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if input.UsuarioID == 0 || input.LibroID == 0 || input.FechaPrestamo == "" || input.FechaDevolucion == "" {
		http.Error(w, "Todos los parámetros son obligatorios", http.StatusBadRequest)
		return
	}

	loan := &data.Loan{
		UsuarioID:       input.UsuarioID,
		LibroID:         input.LibroID,
		FechaPrestamo:   input.FechaPrestamo,
		FechaDevolucion: input.FechaDevolucion,
	}

	err = app.models.Loan.CreateLoan(loan)
	
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"loan": loan}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary Extend loan
// @Description Extend the return date of a loan
// @Tags Prestamos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Param request body object true "New return date"
// @Success 200 {string} string "Loan extended successfully"
// @Router /api/loans/extend/{id} [post]
func (app *application) extendLoanHandler(w http.ResponseWriter, r *http.Request) {

	reservationID := r.PathValue("id")
	
	if reservationID == "" {
		http.Error(w, "ID del reserva no proporcionado", http.StatusBadRequest)
		return
	}
	
	var request struct {
		NuevaFechaDevolucion string `json:"nuevafechadevolucion"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	nuevaFecha, err := time.Parse("2006-01-02", request.NuevaFechaDevolucion)
	if err != nil {
		http.Error(w, "Fecha inválida, el formato debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	err = app.models.Loan.ExtendLoand(reservationID , nuevaFecha )
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Préstamo extendido exitosamente"))
}

// getActiveLoansHandler obtiene todos los préstamos activos en un rango de fechas
//
// @Summary Obtener préstamos activos por rango de fechas
// @Description Retorna todos los préstamos activos en un rango de fechas específico para administradores
// @Tags Prestamos
// @Accept json
// @Produce json
// @Param startdate query string true "Fecha de inicio (YYYY-MM-DD)"
// @Param enddate query string true "Fecha de fin (YYYY-MM-DD)"
// @Success 201 {object} map[string]interface{} "Lista de préstamos activos"
// @Failure 400 {string} string "Parámetros de fecha faltantes o inválidos"
// @Failure 403 {string} string "No tienes permisos para acceder a este recurso"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/admin/loans [post]
// @Security BearerAuth
func (app *application) getActiveLoansHandler(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startdate")
	endDate := r.URL.Query().Get("enddate")

	if startDate == "" || endDate == "" {
		http.Error(w, "startdate y enddate son requeridos", http.StatusBadRequest)
		return
	}

	loans, err := app.models.Loan.GetActiveLoansByDateRange(startDate, endDate)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"loans": loans}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get user loan history
// @Description Get loan history for a specific user
// @Tags Prestamos
// @Accept json
// @Produce json
// @Param idsocio query string true "User ID"
// @Success 200 {array} data.Loan
// @Router /api/loans/history [get]
func (app *application) getUserLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario := r.URL.Query().Get("idsocio")

	loans, err := app.models.Loan.GetUserLoanHistory(idUsuario)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos registrados para este usuario"})
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"loanshystory": loans}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


// getUserActiveLoanStatusHandler obtiene los préstamos activos de un usuario.
//
// @Summary Obtener préstamos activos de un usuario
// @Description Retorna todos los préstamos activos asociados a un usuario por su ID.
// @Tags Prestamos
// @Accept json
// @Produce json
// @Param usuario_id query string true "ID del usuario"
// @Success 201 {object} map[string]interface{} "Lista de préstamos activos del usuario"
// @Failure 400 {string} string "Parámetro 'usuario_id' faltante o inválido"
// @Failure 404 {object} map[string]string "No hay préstamos activos para este usuario"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /api/loans [get]
// @Security BearerAuth
func (app *application) getUserActiveLoanStatusHandler(w http.ResponseWriter, r *http.Request) {

	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	loans, err := app.models.Loan.GetUserActiveLoanStatus(usuarioID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos activos para este usuario"})
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"loansactive": loans}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get user completed loan history
// @Description Get completed loan history for a specific user
// @Tags Prestamos
// @Accept json
// @Produce json
// @Param idsocio query string true "User ID"
// @Success 200 {array} data.Loan
// @Router /api/loans/completed [get]
func (app *application) getUserCompletedLoanHistoryHandler(w http.ResponseWriter, r *http.Request) {

	usuarioID := r.URL.Query().Get("usuario_id")
	if usuarioID == "" {
		http.Error(w, "El parámetro 'usuario_id' es obligatorio", http.StatusBadRequest)
		return
	}

	loans, err := app.models.Loan.GetUserCompletedLoanHistory(usuarioID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(loans) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No hay préstamos completados para este usuario"})
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"loanscomplete": loans}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Return loan
// @Description Mark a loan as returned
// @Tags Prestamos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Loan ID"
// @Success 200 {string} string "Loan returned successfully"
// @Router /api/loans/return/{id} [post]
func (app *application) returnLoanHandler(w http.ResponseWriter, r *http.Request) {
	loanID := r.PathValue("id")
	if loanID == "" {
		http.Error(w, "ID del préstamo no proporcionado", http.StatusBadRequest)
		return
	}

	// Obtener información del préstamo
	loan, err := app.models.Loan.GetLoanByID(loanID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Verificar si el préstamo está activo
	if loan.Estado != "activo" {
		http.Error(w, "El préstamo no está activo", http.StatusBadRequest)
		return
	}

	// Verificar si hay retraso
	fechaDevolucion, err := parseDate(loan.FechaDevolucion)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fechaActual := time.Now()
	var multaGenerada bool
	var multaID int

	// Si hay retraso, generar multa
	if fechaActual.After(fechaDevolucion) {
		diasRetraso := int(fechaActual.Sub(fechaDevolucion).Hours() / 24)
		montoMulta := float64(diasRetraso) * 5.0 // $5 por día de retraso

		// Crear la multa
		fine := &data.Fine{
			IDPrestamo: loan.IDPrestamo,
			SaldoPagar: montoMulta,
			FechaMulta: fechaActual.Format("2006-01-02"),
			Estado:     "pendiente",
		}

		err = app.models.Fine.CreateFine(fine)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		multaGenerada = true
		multaID = fine.IDMulta
	}

	// Marcar préstamo como completado
	err = app.models.Loan.CompleteLoan(loanID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Marcar libro como disponible
	err = app.models.Book.UpdateBookStatus(loan.IDLibro, "disponible")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Preparar respuesta
	response := map[string]interface{}{
		"message": "Libro devuelto exitosamente",
		"loan_id": loanID,
	}

	if multaGenerada {
		response["fine_generated"] = true
		response["fine_id"] = multaID
		response["fine_amount"] = float64(int(fechaActual.Sub(fechaDevolucion).Hours()/24)) * 5.0
		response["days_late"] = int(fechaActual.Sub(fechaDevolucion).Hours() / 24)
	} else {
		response["fine_generated"] = false
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"return": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
