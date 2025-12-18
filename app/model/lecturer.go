package model

import "time"

// Lecturer represents the lecturer profile data in PostgreSQL
type Lecturer struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	UserID     string    `json:"user_id" example:"uuid-user-456"`
	// NIDN atau Nomor Induk Dosen
	LecturerID string    `json:"lecturer_id" example:"198801012015011001"`
	Department string    `json:"department" example:"Teknik Informatika"`
	CreatedAt  time.Time `json:"created_at" swaggerignore:"true"`
}