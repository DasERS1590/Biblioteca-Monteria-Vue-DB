package main

import (
	"biblioteca/internal/data"
	"biblioteca/internal/validator"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credentials struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type Usuario struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Rol      string `json:"rol"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Nombre          string `json:"nombre"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	TipoSocio       string `json:"tipo_socio"`
	Contrasena      string `json:"contrasena"`
	Rol             string `json:"rol"`
}

const (
	// Permisos para libros
	PermissionBooksRead   = "books:read"
	PermissionBooksWrite  = "books:write"
	PermissionBooksDelete = "books:delete"

	// Permisos para préstamos
	PermissionLoansCreate = "loans:create"
	PermissionLoansView   = "loans:view"
	PermissionLoansManage = "loans:manage"

	// Permisos para usuarios
	PermissionUsersView   = "users:view"
	PermissionUsersManage = "users:manage"

	// Permisos especiales
	PermissionReportsGenerate = "reports:generate"

	// Permisos para multas.
	PermissionFinesRead   = "fines:read"
	PermissionFinesCreate = "fines:create"

	// Para reservas
	PermissionReservationsCreate = "reservations:create"
	PermissionReservationsView   = "reservations:view"

	// Permisos para editoriales
	PublishersRead   = "publishers:read"
	PublishersCreate = "publishers:create"

	// Permisos para autores
	AauthorsRead  = "authors:read"
	AuthorsCreate = "authors:create"
)

// Conjuntos de permisos por rol
var (
	UserPermissions = []string{
		PermissionBooksRead,
		PermissionLoansCreate,
		PermissionLoansView,
		PermissionFinesRead,
		PermissionFinesCreate,
		PermissionReservationsCreate,
		PermissionReservationsView,
		PublishersRead,
		PublishersCreate,
		AauthorsRead,
		AuthorsCreate,
	}

	AdminPermissions = []string{
		PermissionBooksRead,
		PermissionBooksWrite,
		PermissionBooksDelete,
		PermissionLoansCreate,
		PermissionLoansView,
		PermissionLoansManage,
		PermissionUsersView,
		PermissionUsersManage,
		PermissionReportsGenerate,
		PermissionFinesRead,
		PermissionFinesCreate,
		PermissionReservationsView,
		PublishersRead,
		PublishersCreate,
		AauthorsRead,
		AuthorsCreate,
	}
)

// @Summary     Login a los Usuarios
// @Tags        Autenticacion
// @Accept      json
// @Produce     json
// @Param       payload body        Credentials true "Credenciales de login"
// @Router      /api/login [post]
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	err := app.readJSON(w, r, &credentials)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, credentials.Correo)
	data.ValidatePasswordPlaintext(v, credentials.Contrasena)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Usuario.GetUserByEmail(credentials.Correo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(credentials.Contrasena)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	claims := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(app.config.jwt.exp).Unix(),
		"iat":    time.Now().Unix(),
		"nbf":    time.Now().Unix(),
		"iss":    app.config.jwt.iss,
		"aud":    app.config.jwt.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     user.ID,
		"nombre": user.Nombre,
		"rol":    user.Rol,
		"token":  token,
	})

}

// @Summary     Registro Para nuevos usuarios
// @Tags        Autenticacion
// @Accept      json
// @Produce     json
// @Param       payload body  RegisterRequest true "Register credentials"
// @Router      /api/register [post]
func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.RegisterRequest{
		Nombre:          req.Nombre,
		Direccion:       req.Direccion,
		Telefono:        req.Telefono,
		Correo:          req.Correo,
		FechaNacimiento: req.FechaNacimiento,
		TipoSocio:       req.TipoSocio,
		Rol:             req.Rol,
	}

	err = user.Contrasena.Set(req.Contrasena)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateUser(v, user)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Usuario.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var permissionsToAdd []string
	if user.Rol == "administrador" {
		permissionsToAdd = AdminPermissions
	} else {
		permissionsToAdd = UserPermissions
	}

	err = app.models.Permissions.AddForUser(user.ID, permissionsToAdd...)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


