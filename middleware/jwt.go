package middleware

import (
    "prestasi_backend/utils"
    "github.com/gofiber/fiber/v2"
)

// JWTRequired memastikan request punya token valid
func JWTRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")

        if tokenString == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization token missing",
            })
        }

        // Format “Bearer <token>”
        userClaims, err := utils.ParseToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }

        // Simpan ke context (FULL claims)
        c.Locals("user", userClaims)

        // Simpan ke context agar service bisa pakai
        c.Locals("userId", userClaims.UserID)
        c.Locals("role", userClaims.RoleName)
        c.Locals("permissions", userClaims.Permissions)

        return c.Next()
    }
}