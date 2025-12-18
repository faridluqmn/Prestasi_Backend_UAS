package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// GetAllPermissions mengambil semua izin yang ada
func GetAllPermissions() ([]model.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		ORDER BY name;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Permission
	for rows.Next() {
		var p model.Permission
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Resource,
			&p.Action,
			&p.Description,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// GetPermissionByID mengambil izin berdasarkan UUID
func GetPermissionByID(id string) (*model.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE id = $1;
	`

	var p model.Permission
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Resource,
		&p.Action,
		&p.Description,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPermissionByName mengambil izin berdasarkan nama
func GetPermissionByName(name string) (*model.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE name = $1;
	`

	var p model.Permission
	err := database.DB.QueryRow(query, name).Scan(
		&p.ID,
		&p.Name,
		&p.Resource,
		&p.Action,
		&p.Description,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}