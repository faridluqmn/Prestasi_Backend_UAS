package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// Ambil semua mahasiswa
func GetAllStudents() ([]model.Student, error) {
	query := `
		SELECT id, student_id, user_id, name, study_program, batch_year, advisor_id, created_at, updated_at
		FROM students
		ORDER BY student_id;
	`
	rows, err := database.DB.Query(query)
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

// Ambil mahasiswa by ID
func GetStudentByID(id string) (*model.Student, error) {
	query := `
		SELECT id, student_id, user_id, name, study_program, batch_year, advisor_id, created_at, updated_at
		FROM students
		WHERE id = $1;
	`
	var s model.Student
	err := database.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.StudentID,
		&s.UserID,
		&s.Name,
		&s.StudyProgram,
		&s.BatchYear,
		&s.AdvisorID,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Ambil mahasiswa berdasarkan user_id (untuk mapping dari JWT)
func GetStudentByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT id, student_id, user_id, name, study_program, batch_year, advisor_id, created_at, updated_at
		FROM students
		WHERE user_id = $1;
	`
	var s model.Student
	err := database.DB.QueryRow(query, userID).Scan(
		&s.ID,
		&s.StudentID,
		&s.UserID,
		&s.Name,
		&s.StudyProgram,
		&s.BatchYear,
		&s.AdvisorID,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Set dosen wali (advisor)
func SetStudentAdvisor(studentID, lecturerID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, lecturerID, studentID)
	return err
}
