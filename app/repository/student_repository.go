package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

func GetAllStudents() ([]model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
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
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func GetStudentByID(id string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1;
	`

	var s model.Student
	err := database.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Ambil student berdasarkan user_id (dipakai untuk mapping dari JWT user)
func GetStudentByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1;
	`

	var s model.Student
	err := database.DB.QueryRow(query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Set dosen wali / advisor
func SetStudentAdvisor(studentID, lecturerID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, lecturerID, studentID)
	return err
}

// Ambil semua mahasiswa bimbingan dari seorang dosen wali
func GetStudentsByAdvisor(lecturerID string) ([]model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE advisor_id = $1
		ORDER BY student_id;
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
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
