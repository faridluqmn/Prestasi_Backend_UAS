package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementMongo represents the dynamic achievement data stored in MongoDB
type AchievementMongo struct {
    ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty" swaggerignore:"true"`
    StudentID       string                 `json:"studentId" bson:"studentId" example:"uuid-student-123"`
    AchievementType string                 `json:"achievementType" bson:"achievementType" example:"competition"` 
    Title           string                 `json:"title" bson:"title" example:"Juara 1 Hackathon Nasional 2025"`
    Description     string                 `json:"description" bson:"description" example:"Memenangkan kompetisi hackathon tingkat nasional"`
    Details         map[string]interface{} `json:"details" bson:"details" swaggertype:"object"` 
    Attachments     []string               `json:"attachments" bson:"attachments" example:"sertifikat_juara.pdf"`
    Tags            []string               `json:"tags" bson:"tags" example:"teknologi,programming"`
    Points          int                    `json:"points" bson:"points" example:"100"`
    CreatedAt       time.Time              `json:"createdAt" bson:"createdAt" swaggerignore:"true"`
    UpdatedAt       time.Time              `json:"updatedAt" bson:"updatedAt" swaggerignore:"true"`
}