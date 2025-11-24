package service

import (
    "prestasi_backend/app/model"
    "prestasi_backend/app/repository"
    "prestasi_backend/utils"

	"time"
	"github.com/golang-jwt/jwt/v5"
    "github.com/gofiber/fiber/v2"
)

func AuthLogin(c *fiber.Ctx) error {
    var req model.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
    }

    // ambil user
    user, err := repository.GetUserByUsername(req.Username)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    // cek password
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
    }

    // ambil role
    role, _ := repository.GetRoleByID(user.RoleID)

    // ambil permissions
    perms, _ := repository.GetPermissionsByRoleID(user.RoleID)

    // siapkan claim
    claim := model.JWTClaims{
        UserID:      user.ID,
        RoleName:    role.Name,
        Permissions: perms,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token, err := utils.GenerateToken(claim)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "token":   token,
        "user": fiber.Map{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
            "full_name": user.FullName,
            "role":      role.Name,
        },
        "permissions": perms,
    })
}

func AuthProfile(c *fiber.Ctx) error {
    claims := c.Locals("user").(*model.JWTClaims)
    user, err := repository.GetUserByID(claims.UserID)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "user":    user,
        "permissions": claims.Permissions,
        "role":        claims.RoleName,
    })
}
