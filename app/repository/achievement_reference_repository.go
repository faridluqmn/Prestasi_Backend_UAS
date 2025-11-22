package repository

import (
	"prestasi_backend/app/model"
	"prestasi_backend/database"
)

// Insert reference baru (ketika mahasiswa create prestasi)
func CreateAchievementReference(ref *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references
		(id, student_id, mongo_achievement_id, status,
		 submitted_at, verified_at, verified_by, rejection_note,
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4,
		        $5, $6, $7, $8,
		        NOW(), NOW());
	`
	_, err := database.DB.Exec(
		query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.SubmittedAt,
		ref.VerifiedAt,
		ref.VerifiedBy,
		ref.RejectionNote,
	)
	return err
}

// Ambil reference by ID
func GetAchievementReferenceByID(id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id = $1;
	`
	var ref model.AchievementReference
	err := database.DB.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

// Ambil semua reference milik mahasiswa
func GetAchievementReferencesByStudentID(studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := database.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, ref)
	}
	return list, rows.Err()
}

// Update status (submit / verify / reject)
func UpdateAchievementStatus(id, status string) error {
	query := `
		UPDATE achievement_references
		SET status = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, status, id)
	return err
}

// Khusus verifikasi (isi verified_at & verified_by)
func VerifyAchievementReference(id, lecturerID string) error {
	query := `
		UPDATE achievement_references
		SET status = 'verified',
		    verified_at = NOW(),
		    verified_by = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, lecturerID, id)
	return err
}

// Khusus reject (isi rejection_note)
func RejectAchievementReference(id, note string) error {
	query := `
		UPDATE achievement_references
		SET status = 'rejected',
		    rejection_note = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, note, id)
	return err
}
