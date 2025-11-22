package model

type Lecturer struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	NIP      string `json:"nip"`
	Name     string `json:"name"`
	Position string `json:"position"` // misal: Lektor, Asisten Ahli, dll
}
