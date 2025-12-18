package service

import (
	"prestasi_backend/app/repository"
	"github.com/gofiber/fiber/v2"
)

// ==================================================================
// LIST LECTURERS
// ==================================================================

// LecturerList godoc
// @Summary      Lihat Daftar Dosen
// @Description  Menampilkan data dosen. Admin bisa melihat semua dosen, Mahasiswa hanya melihat dosen walinya sendiri.
// @Tags         Lecturer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /lecturers [get]
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

// ==================================================================
// LIST ADVISEES (MAHASISWA BIMBINGAN)
// ==================================================================

// LecturerAdvisees godoc
// @Summary      Lihat Mahasiswa Bimbingan
// @Description  Melihat daftar mahasiswa yang dibimbing oleh dosen tertentu. Dosen hanya bisa melihat bimbingannya sendiri.
// @Tags         Lecturer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Lecturer ID"
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /lecturers/{id}/advisees [get]
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