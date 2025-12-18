package model

import "time"

// Role represents a user role in the RBAC system (e.g., Admin, Mahasiswa, Dosen Wali)
type Role struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Name        string    `json:"name" example:"Mahasiswa"`
	Description string    `json:"description" example:"Pengguna yang berhak melaporkan prestasi"`
	CreatedAt   time.Time `json:"created_at" swaggerignore:"true"`
}