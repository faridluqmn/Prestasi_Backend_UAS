package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// List dosen
func GetAllLecturers() ([]model.Lecturer, error) {
	query := `
		SELECT id, nip, user_id, name, position, department, created_at, updated_at
		FROM lecturers
		ORDER BY name;
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
			&l.NIP,
			&l.UserID,
			&l.Name,
			&l.Position,
			&l.Department,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

// Detail dosen
func GetLecturerByID(id string) (*model.Lecturer, error) {
	query := `
		SELECT id, nip, user_id, name, position, department, created_at, updated_at
		FROM lecturers
		WHERE id = $1;
	`
	var l model.Lecturer
	err := database.DB.QueryRow(query, id).Scan(
		&l.ID,
		&l.NIP,
		&l.UserID,
		&l.Name,
		&l.Position,
		&l.Department,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// Ambil daftar mahasiswa bimbingan dosen
func GetAdviseesByLecturerID(lecturerID string) ([]model.Student, error) {
	query := `
		SELECT id, student_id, user_id, name, study_program, batch_year, advisor_id, created_at, updated_at
		FROM students
		WHERE advisor_id = $1;
	`
	rows, err := database.DB.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(
			&s.ID,
			&s.StudentID,
			&s.UserID,
			&s.Name,
			&s.StudyProgram,
			&s.BatchYear,
			&s.AdvisorID,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
