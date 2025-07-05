package data

import (
	"context"
	"database/sql"
	"time"
)

type NewAuthor struct {
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

type Author struct {
	ID           int    `json:"idautor"`
	Nombre       string `json:"nombre"`
	Nacionalidad string `json:"nacionalidad"`
}

type AutorModel struct {
	DB *sql.DB
}

func (m AutorModel) GetAutores() ([]*Author, error) {

	sqlQuery := `SELECT idautor, nombre, nacionalidad FROM autor`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*Author
	for rows.Next() {
		var author Author

		if err := rows.Scan(&author.ID, &author.Nombre, &author.Nacionalidad); err != nil {
			return nil, err
		}
		authors = append(authors, &author)
	}

	return authors, nil

}

func (m AutorModel) CreateAuthor(author *Author) error {

	sqlQuery := `INSERT INTO autor( nombre, nacionalidad ) VALUES(?,?)`

	args := []any{author.Nombre, author.Nacionalidad}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, sqlQuery, args...)

	if err != nil {
		return err
	}

	return nil
}
