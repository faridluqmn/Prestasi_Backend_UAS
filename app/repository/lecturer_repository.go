package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// ambil semua dosen
func GetAllLecturers() ([]model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		ORDER BY lecturer_id;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

// ambil dosen berdasarkan id
func GetLecturerByID(id string) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1;
	`
	var l model.Lecturer
	err := database.DB.QueryRow(query, id).Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// mapping JWT → dosen
func GetLecturerByUserID(userID string) (*model.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1;
	`
	var l model.Lecturer
	err := database.DB.QueryRow(query, userID).Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}