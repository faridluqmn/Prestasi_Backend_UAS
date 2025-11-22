package model

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // jangan tampilkan di JSON
	RoleID   string `json:"role_id"`
	Role     *Role  `json:"role,omitempty"` // optional join
}
