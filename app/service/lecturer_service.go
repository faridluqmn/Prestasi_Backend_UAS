package service

import (
	"prestasi_backend/app/repository"
	"github.com/gofiber/fiber/v2"
)

// GET /lecturers
func LecturerList(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	switch role {

	case "Admin":
		// ✅ Admin → lihat semua dosen
		lects, err := repository.GetAllLecturers()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal mengambil dosen",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    lects,
		})

	case "Mahasiswa":
		// ✅ Mahasiswa → hanya dosen walinya
		stud, err := repository.GetStudentByUserID(userID)
		if err != nil || stud.AdvisorID == "" {
			return c.JSON(fiber.Map{
				"success": true,
				"data":    []any{},
			})
		}

		lect, err := repository.GetLecturerByID(stud.AdvisorID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Dosen wali tidak ditemukan",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    []any{lect},
		})

	case "Dosen Wali":
		// ❌ Dosen tidak boleh lihat daftar dosen
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden",
		})

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "Role tidak dikenali",
		})
	}
}

// GET /lecturers/:id/advisees
func LecturerAdvisees(c *fiber.Ctx) error {
	lectID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	switch role {

	case "Admin":
		// ✅ Admin boleh lihat advisees dosen mana pun
		// lanjut

	case "Dosen Wali":
		// ✅ Dosen hanya boleh lihat advisees dirinya sendiri
		self, err := repository.GetLecturerByUserID(userID)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Data dosen tidak ditemukan",
			})
		}
		if lectID != self.ID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Tidak boleh melihat mahasiswa bimbingan dosen lain",
			})
		}

	default:
		// ❌ Mahasiswa tidak boleh akses endpoint ini
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	students, err := repository.GetStudentsByAdvisor(lectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil mahasiswa bimbingan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    students,
	})
}