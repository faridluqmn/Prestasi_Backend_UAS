package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// GetAllRoles mengambil semua peran yang ada
func GetAllRoles() ([]model.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY created_at;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Role
	for rows.Next() {
		var r model.Role
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, rows.Err()
}

func GetRoleByID(id string) (*model.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE id = $1;
	`

	var r model.Role
	err := database.DB.QueryRow(query, id).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRoleByName(name string) (*model.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE name = $1;
	`

	var r model.Role
	err := database.DB.QueryRow(query, name).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &r, nil
}