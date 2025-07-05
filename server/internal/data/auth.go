package data

import (
	"biblioteca/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &RegisterRequest{}

type RegisterRequest struct {
	ID              int64      `json:"id"`
	Nombre          string   `json:"nombre"`
	Direccion       string   `json:"direccion"`
	Telefono        string   `json:"telefono"`
	Correo          string   `json:"correo"`
	FechaNacimiento string   `json:"fecha_nacimiento"`
	TipoSocio       string   `json:"tipo_socio"`
	Contrasena      password `json:"-"`
	Rol             string   `json:"rol"`
}

type Usuario struct {
	ID       int64      `json:"id"`
	Nombre   string   `json:"nombre"`
	Rol      string   `json:"rol"`
	Email    string   `json:"email"`
	Password password `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)

	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "correo", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "correo", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "contrasena", "must be provided")
	v.Check(len(password) >= 8, "contrasena", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "contrasena", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *RegisterRequest) {
	v.Check(user.Nombre != "", "nombre", "debe ser proporcionado")
	v.Check(len(user.Nombre) <= 500, "nombre", "no debe superar los 500 caracteres")

	v.Check(user.Direccion != "", "direccion", "debe ser proporcionada")
	v.Check(user.Telefono != "", "telefono", "debe ser proporcionado")
	v.Check(user.Correo != "", "correo", "debe ser proporcionado")
	v.Check(validator.Matches(user.Correo, validator.EmailRX), "correo", "debe ser un correo valido	")

	v.Check(user.FechaNacimiento != "", "fecha_nacimiento", "debe ser proporcionada")
	_, err := time.Parse("2006-01-02", user.FechaNacimiento)
	v.Check(err == nil, "fecha_nacimiento", "debe tener formato YYYY-MM-DD")

	v.Check(user.TipoSocio != "", "tipo_socio", "debe ser proporcionado")
	v.Check(user.Rol == "usuario" || user.Rol == "administrador", "rol", "debe ser 'usuario' o 'administrador'")

	if user.Contrasena.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Contrasena.plaintext)
	}

	if user.Contrasena.hash == nil {
		panic("faltante hash de contraseña")
	}
}

func (m UserModel) GetUserByEmail(email string) (*Usuario, error) {

	sqlQuery := `
		SELECT socio.idsocio, socio.nombre, socio.rol, usuariopassword.hash_contrasena
		FROM socio
		JOIN usuariopassword ON socio.idsocio = usuariopassword.idusuario
		WHERE socio.correo = ?
	`
	var usuario Usuario

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, email).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Rol,
		&usuario.Password.hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &usuario, nil
}

func (m UserModel) GetUserByID(id int64) (*RegisterRequest, error) {
	sqlQuery := `
		SELECT socio.idsocio, socio.nombre, socio.rol, usuariopassword.hash_contrasena
		FROM socio
		JOIN usuariopassword ON socio.idsocio = usuariopassword.idusuario
		WHERE socio.idsocio = ?
	`
	var usuario RegisterRequest

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, id).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Rol,
		&usuario.Contrasena.hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &usuario, nil
}

func (m UserModel) Insert(user *RegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var existingID int
	queryCheck := "SELECT idsocio FROM socio WHERE correo = ?"
	err := m.DB.QueryRowContext(ctx, queryCheck, user.Correo).Scan(&existingID)
	if err == nil {
		return ErrDuplicateEmail
	} else if err != sql.ErrNoRows {
		return err
	}

	queryInsertSocio := `
		INSERT INTO socio (nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, imagenperfil, rol)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), ?, ?)
	`

	result, err := m.DB.ExecContext(ctx, queryInsertSocio,
		user.Nombre, user.Direccion, user.Telefono, user.Correo,
		user.FechaNacimiento, user.TipoSocio, nil, user.Rol)
	if err != nil {
		return err
	}

	socioID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	// add the ide to the object 
	user.ID = socioID

	queryInsertPassword := `
		INSERT INTO usuariopassword (idusuario, hash_contrasena)
		VALUES (?, ?)
	`

	_, err = m.DB.ExecContext(ctx, queryInsertPassword, socioID, user.Contrasena.hash)
	if err != nil {
		return err
	}

	return nil
}

func (u *RegisterRequest) IsAnonymous() bool {
	return u == AnonymousUser
}
