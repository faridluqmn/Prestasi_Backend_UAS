package model

import "time"

// Student represents the student profile data in PostgreSQL
type Student struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440003"`
	UserID       string    `json:"user_id" example:"uuid-user-123"`
	StudentID    string    `json:"student_id" example:"2021101234"`
	ProgramStudy string    `json:"program_study" example:"Teknik Informatika"`
	AcademicYear string    `json:"academic_year" example:"2021/2022"`
	AdvisorID    string    `json:"advisor_id" example:"uuid-lecturer-456"`
	CreatedAt    time.Time `json:"created_at" swaggerignore:"true"`
}