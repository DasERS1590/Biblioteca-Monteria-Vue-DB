package models

import (
    "database/sql"
    "errors"
    "biblioteca-api/internal/auth"
    "golang.org/x/crypto/bcrypt"
)

type UserAuth struct {
    ID             int
    Username       string
    PasswordHash   string
    Role           string
}

func (u *UserAuth) HashPassword(password string) error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    u.PasswordHash = string(bytes)
    return err
}

func (u *UserAuth) CheckPasswordHash(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}

func (u *UserAuth) Login(db *sql.DB, username, password string) (string, error) {
    query := "SELECT id, username, hash_contrasena, rol FROM socio WHERE username = ?"
    err := db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role)
    
    if err != nil {
        return "", errors.New("usuario no encontrado")
    }

    if !u.CheckPasswordHash(password) {
        return "", errors.New("contraseña incorrecta")
    }

    token, err := auth.GenerateToken(u.ID, u.Username, u.Role)
    if err != nil {
        return "", errors.New("error generando token")
    }

    return token, nil
}

