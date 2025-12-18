package model

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// User represents the user account data in PostgreSQL
type User struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username     string    `json:"username" example:"mahasiswa123"`
	Email        string    `json:"email" example:"mahasiswa@univ.ac.id"`
	PasswordHash string    `json:"-" swaggerignore:"true"` // Sembunyikan hash dari dokumentasi
	FullName     string    `json:"full_name" example:"John Doe"`
	RoleID       string    `json:"role_id" example:"uuid-role-mahasiswa"`
	IsActive     bool      `json:"is_active" example:"true"`
	CreatedAt    time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt    time.Time `json:"updated_at" swaggerignore:"true"`
}

// LoginRequest digunakan untuk input kredensial saat login (FR-001)
type LoginRequest struct {
	Username string `json:"username" example:"mahasiswa123"`
	Password string `json:"password" example:"SecurePass123!"`
}

// JWTClaims mendefinisikan struktur payload dalam token JWT (FR-001)
type JWTClaims struct {
	UserID      string   `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoleName    string   `json:"role" example:"Mahasiswa"`
	Permissions []string `json:"permissions" example:"achievement:create,achievement:read"`
	jwt.RegisteredClaims
}

// UserCreateRequest digunakan oleh Admin untuk membuat user baru (FR-009)
type UserCreateRequest struct {
	Username string `json:"username" example:"dosen_wali1"`
	Email    string `json:"email" example:"dosen@univ.ac.id"`
	Password string `json:"password" example:"StrictPass2025!"`
	FullName string `json:"full_name" example:"Dr. Ahmad Yani"`
	RoleID   string `json:"role_id" example:"uuid-role-dosen"`
	IsActive bool   `json:"is_active" example:"true"`
}

// UserUpdateRequest digunakan untuk memperbarui profil user
type UserUpdateRequest struct {
	Username string `json:"username" example:"mahasiswa_edit"`
	Email    string `json:"email" example:"newmail@univ.ac.id"`
	Password string `json:"password" example:"NewSecurePass123!"` // optional
	FullName string `json:"full_name" example:"John Doe Updated"`
	IsActive bool   `json:"is_active" example:"true"`
}

// UserUpdateRoleRequest digunakan oleh Admin untuk mengubah role user (FR-009)
type UserUpdateRoleRequest struct {
	RoleID string `json:"role_id" example:"uuid-role-admin"`
}