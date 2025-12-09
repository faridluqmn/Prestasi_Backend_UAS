package model

// Request untuk membuat achievement baru
type AchievementCreateRequest struct {
    StudentID       string                 `json:"student_id"`
    AchievementType string                 `json:"achievement_type"`
    Title           string                 `json:"title"`
    Description     string                 `json:"description"`
    Details         map[string]interface{} `json:"details"`
    Tags            []string               `json:"tags"`
    Points          int                    `json:"points"`
}

// Request untuk update achievement
type AchievementUpdateRequest struct {
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Details     map[string]interface{} `json:"details"`
    Tags        []string               `json:"tags"`
    Points      int                    `json:"points"`
}

// Request untuk reject
type AchievementRejectRequest struct {
    Note string `json:"note"`
}
