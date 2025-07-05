package data

import (
	"context"
	"database/sql"
	"time"
)

type User_ struct {
	ID              int    `json:"id"`
	Nombre          string `json:"nombre"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo"`
	FechaNacimiento string `json:"fechanacimiento"`
	TipoSocio       string `json:"tiposocio"`
	FechaRegistro   string `json:"fecharegistro"`
	Rol             string `json:"rol"`
}

type UserModel_ struct {
	DB *sql.DB
}

func (m UserModel_)GetUserByType(userType string) ([]*User_, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `
		SELECT 
			idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, rol 
		FROM socio 
		WHERE tiposocio = ?`

	rows, err := m.DB.QueryContext( ctx , query, userType)

	if err != nil {
		return nil , err 
	}

	defer rows.Close()

	var users []*User_

	for rows.Next() {
		var user User_
		err := rows.Scan(&user.ID, &user.Nombre, &user.Direccion, &user.Telefono, &user.Correo, &user.FechaNacimiento, &user.TipoSocio, &user.FechaRegistro, &user.Rol)
		if err != nil {
			return nil , err 
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil , err 
	}

	return users , nil 
}

func (m UserModel_)GetAllUsers() ([]*User_, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `
		SELECT 
			idsocio, nombre, direccion, telefono, correo, fechanacimiento, tiposocio, fecharegistro, rol 
		FROM socio`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil , err 
	}

	defer rows.Close()

	var users []*User_

	for rows.Next() {
		var user User_
		err := rows.Scan(&user.ID, &user.Nombre, &user.Direccion, &user.Telefono, &user.Correo, &user.FechaNacimiento, &user.TipoSocio, &user.FechaRegistro, &user.Rol)
		if err != nil {
			return nil , err 
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil , err 
	}

	return users , nil 
}

