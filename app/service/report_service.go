package service

import (
	"prestasi_backend/app/repository"

	"github.com/gofiber/fiber/v2"
)


// =======================================================
// GET /api/v1/reports/statistics
// Admin only
// =======================================================
func ReportStatistics(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	// 🔒 hanya admin
	if role != "Admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Hanya admin yang dapat mengakses statistik",
		})
	}

	stats, err := repository.GetAchievementStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil statistik prestasi",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// =======================================================
// GET /api/v1/reports/student/:id
// Admin | Dosen Wali | Mahasiswa (ownership enforced)
// =======================================================
func ReportStudent(c *fiber.Ctx) error {
	targetStudentID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	// ambil mahasiswa target
	targetStudent, err := repository.GetStudentByID(targetStudentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Mahasiswa tidak ditemukan",
		})
	}

	// ===================================================
	// ROLE & OWNERSHIP VALIDATION
	// ===================================================
	switch role {

	case "Mahasiswa":
		// mahasiswa hanya boleh lihat laporan dirinya sendiri
		self, err := repository.GetStudentByUserID(userID)
		if err != nil || self.ID != targetStudent.ID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Tidak boleh melihat laporan mahasiswa lain",
			})
		}

	case "Dosen Wali":
		// dosen hanya boleh lihat mahasiswa bimbingannya
		lect, err := repository.GetLecturerByUserID(userID)
		if err != nil || targetStudent.AdvisorID != lect.ID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Mahasiswa bukan bimbingan Anda",
			})
		}

	case "Admin":
		// admin bebas

	default:
		return c.Status(403).JSON(fiber.Map{
			"error": "Role tidak dikenali",
		})
	}

	// ===================================================
	// DATA AGGREGATION
	// ===================================================
	stats, err := repository.GetStudentAchievementStats(targetStudent.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil statistik mahasiswa",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"student": targetStudent,
		"stats":   stats,
	})
}