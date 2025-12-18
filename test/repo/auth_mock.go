package repo

import (
	"prestasi_backend/app/model"
	"prestasi_backend/app/repository" // Pastikan import ini benar
)

// Fungsi ini buat nipu service biar gak manggil DB asli
func MockGetUserByUsername(mockUser *model.User, mockErr error) {
	repository.GetUserByUsername = func(username string) (*model.User, error) {
		return mockUser, mockErr
	}
}

func MockGetPermissions(mockPerms []string, mockErr error) {
	repository.GetPermissionsByRoleID = func(roleID string) ([]string, error) {
		return mockPerms, mockErr
	}
}