package services

import (
	"testing"
	"prestasi_backend/app/model"
	"prestasi_backend/app/service"
	"prestasi_backend/test/repo"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func TestLogin_Success(t *testing.T) {
	// 1. SETUP PASSWORD YANG PASTI COCOK
	passwordInput := "rahasia123" // Password mentah
	
	// Kita bikin hash ASLI dari password di atas (biar bcrypt gak bingung)
	hashAsli, _ := bcrypt.GenerateFromPassword([]byte(passwordInput), bcrypt.DefaultCost)
	
	// 2. DATA DUMMY
	userPalsu := &model.User{
		Username:     "mhs_test",
		// PENTING: Masukkan hash yang baru kita generate tadi
		PasswordHash: string(hashAsli), 
		RoleID:       "role-mhs",
	}

	// 3. PASANG MOCK
	repo.MockGetUserByUsername(userPalsu, nil)
	repo.MockGetPermissions([]string{"achievement:create"}, nil)

	// 4. JALANKAN LOGIN
	// PENTING: Login pakai password yang SAMA dengan variabel 'passwordInput'
	token, err := service.Login("mhs_test", passwordInput)

	// 5. ASSERTION
	if err != nil {
		t.Fatalf("Login gagal: %v", err) // Kalau ini error, berarti aneh banget
	}
	if token == "" {
		t.Error("Token kosong padahal login sukses")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	// 1. Setup User dengan password hash benar
	passwordBenar := "rahasia123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(passwordBenar), bcrypt.DefaultCost)
	
	userDummy := &model.User{
		Username:     "mhs_test",
		PasswordHash: string(hash),
		RoleID:       "role-mhs",
	}

	// 2. Inject Mock (User ditemukan)
	repo.MockGetUserByUsername(userDummy, nil)

	// 3. Login dengan PASSWORD SALAH
	token, err := service.Login("mhs_test", "password_ngawur")

	// 4. Assert: Harusnya Error
	if err == nil {
		t.Error("Harusnya error karena password salah, tapi malah sukses")
	}
	if token != "" {
		t.Error("Harusnya tidak dapat token")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	// 1. Inject Mock ERROR (User tidak ditemukan)
	// Kita simulasikan repository mengembalikan error
	repo.MockGetUserByUsername(nil, errors.New("record not found"))

	// 2. Login sembarang
	_, err := service.Login("hantu", "123")

	// 3. Assert
	if err == nil {
		t.Error("Harusnya error User Not Found, tapi malah sukses")
	}
}