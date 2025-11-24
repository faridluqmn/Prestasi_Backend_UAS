package middleware

import (
    "github.com/gofiber/fiber/v2"
)

// PermissionRequired memastikan user punya izin tertentu
func PermissionRequired(required string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        perms := c.Locals("permissions")
        if perms == nil {
            return c.Status(403).JSON(fiber.Map{
                "error": "Permissions not loaded",
            })
        }

        // permissions disimpan sebagai []string
        userPerms := perms.([]string)

        // cek apakah required permission ada di userPerms
        for _, p := range userPerms {
            if p == required {
                return c.Next()
            }
        }

        return c.Status(403).JSON(fiber.Map{
            "error": "You don't have permission to perform this action",
        })
    }
}
