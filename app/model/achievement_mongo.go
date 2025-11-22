package model

import (
	"time"
)

type AchievementMongo struct {
	ID              string                 `json:"id" bson:"_id,omitempty"`
	StudentID       string                 `json:"studentId" bson:"studentId"`
	AchievementType string                 `json:"achievementType" bson:"achievementType"`
	Title           string                 `json:"title" bson:"title"`
	Description     string                 `json:"description" bson:"description"`
	Details         map[string]interface{} `json:"details" bson:"details"`
	Attachments     []string               `json:"attachments" bson:"attachments"`
	Tags            []string               `json:"tags" bson:"tags"`
	Points          int                    `json:"points" bson:"points"`
	CreatedAt       time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt" bson:"updatedAt"`
}
