package repo

import (
	"prestasi_backend/app/model"
	"prestasi_backend/app/repository" // Import repository asli
)

func MockCreateAchievement(mockID string, mockErr error) {
	repository.CreateAchievement = func(doc *model.AchievementMongo) (string, error) {
		return mockID, mockErr
	}
}

func MockGetByStudent(mockList []model.AchievementMongo, mockErr error) {
	repository.GetAchievementsByStudentID = func(studentID string) ([]model.AchievementMongo, error) {
		return mockList, mockErr
	}
}