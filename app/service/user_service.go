package service

import (
	"prestasi_backend/app/model"
	"prestasi_backend/app/repository"
	"prestasi_backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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