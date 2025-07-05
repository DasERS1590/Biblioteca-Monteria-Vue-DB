package data

import (
	"context"
	"database/sql"
	"time"
)

type Editorial struct {
	Nombre    string `json:"nombre"`
	Direccion string `json:"direccion"`
	PaginaWeb string `json:"paginaweb"`
}

type  editorial struct {
	IdEditorial int    `json:"ideditorial"`
	Nombre      string `json:"nombre"`
	Direccion   string `json:"direccion"`
	PaginaWeb   string `json:"paginaweb"`
}



type EditorialModel struct {
	DB *sql.DB
}

func (m EditorialModel) CreateEditorial(editorial Editorial) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{editorial.Nombre, editorial.Direccion, editorial.PaginaWeb}

	insertEditorialQuery := `
        INSERT INTO editorial (nombre, direccion, paginaweb)
        VALUES ( ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, insertEditorialQuery, args...)
	if err != nil {
		return err
	}
    
	return nil
} 

func (m EditorialModel) GetEditorials() ([]*editorial , error ){

	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext( ctx, "SELECT * FROM editorial")
	if err != nil {
		return nil , err 
	}
	defer rows.Close()


	var editorials []*editorial 

	for rows.Next() {
		var edt editorial 
		err := rows.Scan(&edt.IdEditorial, &edt.Nombre, &edt.Direccion, &edt.PaginaWeb); 
		
		if err != nil {
			return nil , err 
		}
		editorials = append(editorials, &edt)
	}

	return editorials , nil 
}	
