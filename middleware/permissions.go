package middleware

import (
    "prestasi_backend/app/model"
    "github.com/gofiber/fiber/v2"
)

// PermissionRequired memastikan user punya izin tertentu
func PermissionRequired(required string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        // Ambil claims dari middleware JWTRequired
        claims, ok := c.Locals("user").(*model.JWTClaims)
        if !ok {
            return c.Status(403).JSON(fiber.Map{
                "error": "Unauthorized: invalid token claims",
            })
        }

        // ambil permission list dari token
        userPerms := claims.Permissions

        // cek apakah permission sesuai
        for _, p := range userPerms {
            if p == required {
                return c.Next()
            }
        }

        return c.Status(403).JSON(fiber.Map{
            "error": "Forbidden: missing permission (" + required + ")",
        })
    }
}