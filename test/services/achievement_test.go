package services

import (
	"testing"
	"prestasi_backend/app/model"
	"prestasi_backend/app/service"
	"prestasi_backend/test/repo"
)

func TestSubmitAchievement_Success(t *testing.T) {
	// 1. Setup Data Dummy (SESUAI MODEL BARU)
	input := model.AchievementMongo{
		StudentID: "student-123",
		Title:     "Juara 1 Lomba Coding",
		// Masukkan data spesifik ke dalam Map 'Details'
		Details: map[string]interface{}{
			"competitionName": "Lomba Coding Nasional",
			"rank":            1,
		},
	}

	// 2. Inject Mock
	repo.MockCreateAchievement("mongo-id-999", nil)

	// 3. Panggil Service
	// Pakai &input karena biasanya repo butuh pointer, atau sesuaikan repo kamu
	id, err := service.SubmitAchievement(input)

	// 4. Assert
	if err != nil {
		t.Errorf("Harusnya sukses, tapi error: %v", err)
	}
	if id != "mongo-id-999" {
		t.Errorf("ID salah: %s", id)
	}
}

func TestSubmitAchievement_InvalidInput(t *testing.T) {
	// Case: Competition Name kosong di dalam Details
	input := model.AchievementMongo{
		Title: "Judul Doang",
		Details: map[string]interface{}{
			"competitionName": "", // Kosong
			"rank":            1,
		},
	}

	_, err := service.SubmitAchievement(input)

	if err == nil {
		t.Error("Harusnya error karena competitionName kosong")
	}
}

func TestGetAchievements_Found(t *testing.T) {
	// 1. Siapkan Data Dummy (Isinya 2 prestasi)
	mockData := []model.AchievementMongo{
		{Title: "Juara 1"},
		{Title: "Juara 2"},
	}

	// 2. Inject Mock
	repo.MockGetByStudent(mockData, nil)

	// 3. Panggil Service
	result, err := service.GetStudentAchievements("student-123")

	// 4. Assert
	if err != nil {
		t.Errorf("Harusnya sukses, tapi error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Harusnya dapet 2 prestasi, tapi dapet %d", len(result))
	}
	if result[0].Title != "Juara 1" {
		t.Error("Data urutan pertama salah")
	}
}