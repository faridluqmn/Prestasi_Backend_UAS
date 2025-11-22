package model

type Student struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StudentID string `json:"student_id"` // NIM
	AdvisorID string `json:"advisor_id"` // relasi ke lecturer
}
