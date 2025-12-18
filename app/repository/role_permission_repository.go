package repository

import "prestasi_backend/database"

// ambil permissions by role ID
var GetPermissionsByRoleID = func(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1;
	`

	rows, err := database.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, rows.Err()
}

// Tambah relasi role-permission
func AddPermissionToRole(roleID, permissionID string) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`
	_, err := database.DB.Exec(query, roleID, permissionID)
	return err
}

// Hapus relasi role-permission
func RemovePermissionFromRole(roleID, permissionID string) error {
	query := `
		DELETE FROM role_permissions
		WHERE role_id = $1 AND permission_id = $2;
	`
	_, err := database.DB.Exec(query, roleID, permissionID)
	return err
}
