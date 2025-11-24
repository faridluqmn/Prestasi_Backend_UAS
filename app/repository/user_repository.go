package repository

import (
	"database/sql"
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// Ambil semua user (untuk list admin)
func GetAllUsers() ([]model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
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
			&u.Email,
			&u.PasswordHash,   // password_hash
			&u.FullName,
			&u.RoleID,
			&u.IsActive,
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
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	var u model.User
	err := database.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
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
		SELECT id, username, email, password_hash, full_name,
		       role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1;
	`

	var u model.User
	err := database.DB.QueryRow(query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
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
		INSERT INTO users (
			id, username, email, password_hash, full_name,
			role_id, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5,
		        $6, $7, NOW(), NOW());
	`
	_, err := database.DB.Exec(
		query,
		u.ID,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.FullName,
		u.RoleID,
		u.IsActive,
	)
	return err
}

// Update data user (kecuali role_id)
func UpdateUser(u *model.User) error {
	query := `
		UPDATE users
		SET username = $1,
		    email = $2,
		    password_hash = $3,
		    full_name = $4,
		    is_active = $5,
		    updated_at = NOW()
		WHERE id = $6;
	`
	_, err := database.DB.Exec(
		query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.FullName,
		u.IsActive,
		u.ID,
	)
	return err
}

// Hapus user (hard delete)
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

// Helper: cek apakah error karena row tidak ditemukan
func IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}