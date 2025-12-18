package service

import (
	"prestasi_backend/app/model"
	"prestasi_backend/app/repository"
	"prestasi_backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ==================================================================
// LIST USERS
// ==================================================================

// UserList godoc
// @Summary      Lihat Semua User
// @Description  Menampilkan daftar seluruh user di sistem (Biasanya Admin Only).
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users [get]
func UserList(c *fiber.Ctx) error {

	users, err := repository.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

// ==================================================================
// DETAIL USER
// ==================================================================

// UserDetail godoc
// @Summary      Detail User
// @Description  Melihat detail satu user berdasarkan ID.
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users/{id} [get]
func UserDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := repository.GetUserByID(id)
	if err != nil {
		if repository.IsNoRows(err) {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil detail user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// ==================================================================
// CREATE USER
// ==================================================================

// UserCreate godoc
// @Summary      Buat User Baru
// @Description  Menambahkan user baru secara manual (Admin). Password akan otomatis di-hash.
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.UserCreateRequest true "Data User Baru"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users [post]
func UserCreate(c *fiber.Ctx) error {

	var req model.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	user := model.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     req.IsActive,
	}

	if err := repository.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User berhasil dibuat",
		"data":    user,
	})
}

// ==================================================================
// UPDATE USER
// ==================================================================

// UserUpdate godoc
// @Summary      Update Data User
// @Description  Mengubah data user (Username, Email, Nama, Password, Status Aktif).
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path  string                  true  "User ID"
// @Param        request body  model.UserUpdateRequest true  "Data Update"
// @Success      200     {object} map[string]interface{}
// @Failure      400     {object} map[string]interface{}
// @Failure      404     {object} map[string]interface{}
// @Failure      500     {object} map[string]interface{}
// @Router       /users/{id} [put]
func UserUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		if repository.IsNoRows(err) {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil user"})
	}

	// hash password jika dikirim
	if req.Password != "" {
		hash, _ := utils.HashPassword(req.Password)
		user.PasswordHash = hash
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.IsActive = req.IsActive

	if err := repository.UpdateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User berhasil diupdate",
		"data":    user,
	})
}

// ==================================================================
// DELETE USER
// ==================================================================

// UserDelete godoc
// @Summary      Hapus User
// @Description  Menghapus user dari database secara permanen.
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users/{id} [delete]
func UserDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeleteUser(id)
	if err != nil {
		if repository.IsNoRows(err) {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User berhasil dihapus",
	})
}

// ==================================================================
// UPDATE ROLE
// ==================================================================

// UserUpdateRole godoc
// @Summary      Ganti Role User
// @Description  Mengubah role/hak akses user (Misal: dari Mahasiswa ke Admin).
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path  string                     true  "User ID"
// @Param        request body  model.UserUpdateRoleRequest true  "Role ID Baru"
// @Success      200     {object} map[string]interface{}
// @Failure      400     {object} map[string]interface{}
// @Failure      404     {object} map[string]interface{}
// @Failure      500     {object} map[string]interface{}
// @Router       /users/{id}/role [put]
func UserUpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req model.UserUpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	// cek user ada atau tidak
	_, err := repository.GetUserByID(id)
	if err != nil {
		if repository.IsNoRows(err) {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil user"})
	}

	// update role
	if err := repository.UpdateUserRole(id, req.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update role user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Role user berhasil diperbarui",
	})
}