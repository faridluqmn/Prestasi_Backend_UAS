package model

import "time"

// AchievementReference represents the relational mapping and status of an achievement
type AchievementReference struct {
	ID                 string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StudentID          string     `json:"student_id" example:"uuid-student-123"`
	MongoAchievementID string     `json:"mongo_achievement_id" example:"654321098765432109876543"`
	Status             string     `json:"status" example:"submitted"` 
	SubmittedAt        *time.Time `json:"submitted_at" swaggerignore:"true"`
	VerifiedAt         *time.Time `json:"verified_at" swaggerignore:"true"`
	VerifiedBy         *string    `json:"verified_by" example:"uuid-lecturer-456"`
	RejectionNote      *string    `json:"rejection_note" example:"Bukti sertifikat tidak terbaca atau buram"`
	CreatedAt          time.Time  `json:"created_at" swaggerignore:"true"`
	UpdatedAt          time.Time  `json:"updated_at" swaggerignore:"true"`
}