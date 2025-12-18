package repository

import (
	"prestasi_backend/database"
)

type AchievementStats struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
}

// GetAchievementStats
func GetAchievementStats() ([]AchievementStats, error) {
	query := `
		SELECT status, COUNT(*) as total
		FROM achievement_references
		WHERE status <> 'deleted'
		GROUP BY status;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []AchievementStats
	for rows.Next() {
		var s AchievementStats
		if err := rows.Scan(&s.Status, &s.Total); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetStudentAchievementStats
func GetStudentAchievementStats(studentID string) ([]AchievementStats, error) {
	query := `
		SELECT status, COUNT(*) as total
		FROM achievement_references
		WHERE student_id = $1
		  AND status <> 'deleted'
		GROUP BY status;
	`

	rows, err := database.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []AchievementStats
	for rows.Next() {
		var s AchievementStats
		if err := rows.Scan(&s.Status, &s.Total); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}