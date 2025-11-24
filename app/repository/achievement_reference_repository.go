package repository

import (
	"database/sql"
	"prestasi_backend/app/model"
	"prestasi_backend/database"
	"time"
)

// =====================================
// ACHIEVEMENT_REFERENCES (Postgre)
// =====================================

// Insert reference baru (ketika mahasiswa membuat prestasi)
func CreateAchievementReference(ref *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by, rejection_note,
			created_at, updated_at
		)
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

// Ambil reference berdasarkan ID
func GetAchievementReferenceByID(id string) (*model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id = $1;
	`

	var ref model.AchievementReference
	var submittedAt, verifiedAt sql.NullTime
	var verifiedBy, rejectionNote sql.NullString

	err := database.DB.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&submittedAt,
		&verifiedAt,
		&verifiedBy,
		&rejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if submittedAt.Valid {
		t := submittedAt.Time
		ref.SubmittedAt = &t
	}
	if verifiedAt.Valid {
		t := verifiedAt.Time
		ref.VerifiedAt = &t
	}
	if verifiedBy.Valid {
		s := verifiedBy.String
		ref.VerifiedBy = &s
	}
	if rejectionNote.Valid {
		s := rejectionNote.String
		ref.RejectionNote = &s
	}

	return &ref, nil
}

// Ambil daftar reference milik 1 mahasiswa
func GetAchievementReferencesByStudentID(studentID string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		  AND status <> 'deleted'
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
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy, rejectionNote sql.NullString

		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&submittedAt,
			&verifiedAt,
			&verifiedBy,
			&rejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if submittedAt.Valid {
			t := submittedAt.Time
			ref.SubmittedAt = &t
		}
		if verifiedAt.Valid {
			t := verifiedAt.Time
			ref.VerifiedAt = &t
		}
		if verifiedBy.Valid {
			s := verifiedBy.String
			ref.VerifiedBy = &s
		}
		if rejectionNote.Valid {
			s := rejectionNote.String
			ref.RejectionNote = &s
		}

		list = append(list, ref)
	}
	return list, rows.Err()
}

// Daftar semua reference (untuk admin)
func GetAllAchievementReferences() ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy, rejectionNote sql.NullString

		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&submittedAt,
			&verifiedAt,
			&verifiedBy,
			&rejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if submittedAt.Valid {
			t := submittedAt.Time
			ref.SubmittedAt = &t
		}
		if verifiedAt.Valid {
			t := verifiedAt.Time
			ref.VerifiedAt = &t
		}
		if verifiedBy.Valid {
			s := verifiedBy.String
			ref.VerifiedBy = &s
		}
		if rejectionNote.Valid {
			s := rejectionNote.String
			ref.RejectionNote = &s
		}

		list = append(list, ref)
	}
	return list, rows.Err()
}

// Daftar reference untuk dosen wali (mahasiswa bimbingannya)
func GetAchievementReferencesByAdvisor(lecturerID string) ([]model.AchievementReference, error) {
	// join students untuk filter advisor_id
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
		       ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
		       ar.created_at, ar.updated_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.advisor_id = $1
		  AND ar.status <> 'deleted'
		ORDER BY ar.created_at DESC;
	`

	rows, err := database.DB.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy, rejectionNote sql.NullString

		if err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&submittedAt,
			&verifiedAt,
			&verifiedBy,
			&rejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if submittedAt.Valid {
			t := submittedAt.Time
			ref.SubmittedAt = &t
		}
		if verifiedAt.Valid {
			t := verifiedAt.Time
			ref.VerifiedAt = &t
		}
		if verifiedBy.Valid {
			s := verifiedBy.String
			ref.VerifiedBy = &s
		}
		if rejectionNote.Valid {
			s := rejectionNote.String
			ref.RejectionNote = &s
		}

		list = append(list, ref)
	}
	return list, rows.Err()
}

// Update status umum
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

// Submit: ubah jadi submitted + timestamp
func SubmitAchievementReference(id string) error {
	now := time.Now()
	query := `
		UPDATE achievement_references
		SET status = 'submitted',
		    submitted_at = $1,
		    updated_at = NOW()
		WHERE id = $2;
	`
	_, err := database.DB.Exec(query, now, id)
	return err
}

// Verifikasi: status verified + verified_at + verified_by
func VerifyAchievementReference(id, verifierUserID string) error {
	now := time.Now()
	query := `
		UPDATE achievement_references
		SET status = 'verified',
		    verified_at = $1,
		    verified_by = $2,
		    updated_at = NOW()
		WHERE id = $3;
	`
	_, err := database.DB.Exec(query, now, verifierUserID, id)
	return err
}

// Reject: status rejected + note
func RejectAchievementReference(id, verifierUserID, note string) error {
	now := time.Now()
	query := `
		UPDATE achievement_references
		SET status = 'rejected',
		    verified_at = $1,
		    verified_by = $2,
		    rejection_note = $3,
		    updated_at = NOW()
		WHERE id = $4;
	`
	_, err := database.DB.Exec(query, now, verifierUserID, note, id)
	return err
}

// Soft delete: status deleted
func SoftDeleteAchievementReference(id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'deleted',
		    updated_at = NOW()
		WHERE id = $1;
	`
	_, err := database.DB.Exec(query, id)
	return err
}