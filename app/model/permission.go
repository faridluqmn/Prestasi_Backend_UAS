package model

// Permission represents the access rights for a specific resource and action (RBAC)
type Permission struct {
	ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440005"`
	Name        string `json:"name" example:"achievement:create"`
	Resource    string `json:"resource" example:"achievement"`
	Action      string `json:"action" example:"create"`
	Description string `json:"description" example:"Memberikan akses untuk membuat laporan prestasi baru"`
}