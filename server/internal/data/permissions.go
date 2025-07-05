package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)


type Permissions []string

type PermissionModel struct {
	DB *sql.DB
}

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}


func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		INNER JOIN socio ON users_permissions.user_id = socio.idsocio
		WHERE socio.idsocio = ?
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	if len(codes) == 0 {
		return nil
	}
	placeholders := make([]string, len(codes))
	args := make([]interface{}, 0, len(codes)+1)
	args = append(args, userID)

	for i, code := range codes {
		placeholders[i] = "?"
		args = append(args, code)
	}

	query := fmt.Sprintf(`
		INSERT INTO users_permissions (user_id, permission_id)
		SELECT ?, id FROM permissions WHERE code IN (%s)
	`, strings.Join(placeholders, ", "))


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	
	return err
}
