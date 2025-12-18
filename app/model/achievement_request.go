package model

// AchievementCreateRequest digunakan oleh Mahasiswa untuk input prestasi baru (FR-003)
type AchievementCreateRequest struct {
	StudentID       string                 `json:"student_id" example:"uuid-student-123"`
	AchievementType string                 `json:"achievement_type" example:"competition"`
	Title           string                 `json:"title" example:"Juara 1 Hackathon Nasional 2025"`
	Description     string                 `json:"description" example:"Memenangkan kompetisi hackathon tingkat nasional"`
	Details         map[string]interface{} `json:"details" swaggertype:"object" example:"competitionName:Indonesia Tech Innovation Challenge,rank:1"`
	Tags            []string               `json:"tags" example:"teknologi,programming"`
	Points          int                    `json:"points" example:"100"`
}

// AchievementUpdateRequest digunakan untuk memperbarui prestasi yang masih berstatus 'draft' (FR-003)
type AchievementUpdateRequest struct {
	Title       string                 `json:"title" example:"Juara 1 Hackathon Nasional 2025 (Updated)"`
	Description string                 `json:"description" example:"Update deskripsi prestasi"`
	Details     map[string]interface{} `json:"details" swaggertype:"object"`
	Tags        []string               `json:"tags" example:"teknologi"`
	Points      int                    `json:"points" example:"100"`
}

// AchievementRejectRequest digunakan oleh Dosen Wali untuk memberikan alasan penolakan (FR-008)
type AchievementRejectRequest struct {
	// Alasan mengapa prestasi ditolak
	Note string `json:"note" example:"Sertifikat tidak valid atau kadaluarsa"`
}