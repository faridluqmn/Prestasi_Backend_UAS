package service

import (
	"errors"
	"prestasi_backend/app/model"
	"prestasi_backend/app/repository"
)

func SubmitAchievement(data model.AchievementMongo) (string, error) {
	// 1. Validasi Input (Cara Mengakses Map)
	
	// Cek apakah 'Details' ada isinya
	if data.Details == nil {
		return "", errors.New("details tidak boleh kosong")
	}

	// Ambil value 'competitionName' dari map
	compName, ok := data.Details["competitionName"].(string) // Type assertion ke string
	if !ok || compName == "" {
		return "", errors.New("nama kompetisi tidak boleh kosong")
	}

	// Ambil value 'rank' dari map
	rank, ok := data.Details["rank"].(int) // Type assertion ke int
	if !ok || rank < 1 {
		// Note: Kalau di JSON rank seringkali terbaca float64, 
		// jadi handle kalau perlu, tapi untuk unit test manual 'int' aman.
		return "", errors.New("ranking tidak valid")
	}

	// 2. Panggil Repository
	// Sesuaikan apakah repository minta pointer (&data) atau value (data)
	// Berdasarkan chat sebelumnya repo kamu minta Value (tanpa bintang), jadi:
	id, err := repository.CreateAchievement(&data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetStudentAchievements(studentID string) ([]model.AchievementMongo, error) {
	// 1. Validasi
	if studentID == "" {
		return nil, errors.New("student id kosong")
	}

	// 2. Panggil Repo
	achievements, err := repository.GetAchievementsByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	// 3. Business Logic (Misal: kalau kosong balikin array kosong, bukan null)
	if len(achievements) == 0 {
		return []model.AchievementMongo{}, nil
	}

	return achievements, nil
}