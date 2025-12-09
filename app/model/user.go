package model

import (
    "time"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	PasswordHash  string `json:"-"`
	FullName  string    `json:"full_name"`
	RoleID    string    `json:"role_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTClaims struct {
	UserID      string   `json:"user_id"`
	RoleName    string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type UserCreateRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	RoleID    string `json:"role_id"`
	IsActive  bool   `json:"is_active"`
}

type UserUpdateRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"` // optional
	FullName  string `json:"full_name"`
	IsActive  bool   `json:"is_active"`
}

type UserUpdateRoleRequest struct {
	RoleID string `json:"role_id"`
}
