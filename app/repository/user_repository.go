package repository

import (
	"database/sql"
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// Ambil semua user (untuk admin list)
func GetAllUsers() ([]model.User, error) {
	query := `
		SELECT id, username, password, role_id, created_at, updated_at
		FROM users
		ORDER BY created_at DESC;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.RoleID,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

// Ambil user berdasarkan ID
func GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, username, password, role_id, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	var u model.User
	err := database.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Ambil user berdasarkan username (untuk login)
func GetUserByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, password, role_id, created_at, updated_at
		FROM users
		WHERE username = $1;
	`

	var u model.User
	err := database.DB.QueryRow(query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Tambah user baru
func CreateUser(u *model.User) error {
	query := `
		INSERT INTO users (id, username, password, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW());
	`
	_, err := database.DB.Exec(query, u.ID, u.Username, u.Password, u.RoleID)
	return err
}

// Update user (tanpa ganti role)
func UpdateUser(id string, u *model.User) error {
	query := `
		UPDATE users
		SET username = $1,
		    password = $2,
		    updated_at = NOW()
		WHERE id = $3;
	`
	_, err := database.DB.Exec(query, u.Username, u.Password, id)
	return err
}

// Delete user
func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := database.DB.Exec(query, id)
	return err
}

// Update role user
func UpdateUserRole(userID, roleID string) error {
	query := `
		UPDATE users
		SET role_id = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, roleID, userID)
	return err
}

// (opsional) hitung total user
func CountUsers() (int64, error) {
	var total int64
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&total)
	return total, err
}

// helper kalau mau cek err no rows
func IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}
