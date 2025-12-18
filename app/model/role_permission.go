package model

// RolePermission represents the mapping between roles and their assigned permissions (Many-to-Many)
type RolePermission struct {
	RoleID       string `json:"role_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	PermissionID string `json:"permission_id" example:"550e8400-e29b-41d4-a716-446655440005"`
}