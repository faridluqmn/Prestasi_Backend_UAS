package service

import (
	"fmt"
	"prestasi_backend/app/repository"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func (s *AuthService) Login(username, password string) (string, error) {
	// TODO: validasi user & generate JWT akan dikerjakan di versi final
	return "", fmt.Errorf("login belum diimplementasikan")
}
