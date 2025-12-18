package service

import (
	"github.com/gofiber/fiber/v2"
	"prestasi_backend/app/repository"
)

// GET /students
func StudentList(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	// ADMIN → semua mahasiswa
	if role == "Admin" {
		students, err := repository.GetAllStudents()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil mahasiswa"})
		}
		return c.JSON(fiber.Map{"success": true, "data": students})
	}

	// DOSEN WALI → mahasiswa bimbingannya
	if role == "Dosen Wali" {
		lect, err := repository.GetLecturerByUserID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Data dosen tidak ditemukan",
			})
		}

		students, err := repository.GetStudentsByAdvisor(lect.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal mengambil mahasiswa bimbingan",
			})
		}

		return c.JSON(fiber.Map{"success": true, "data": students})
	}

	// MAHASISWA → tidak boleh lihat semua
	return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
}

// GET /students/:id
func StudentDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	stud, err := repository.GetStudentByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa tidak ditemukan"})
	}

	// mahasiswa hanya boleh lihat dirinya sendiri
	if role == "Mahasiswa" {
		self, _ := repository.GetStudentByUserID(userID)
		if self.ID != stud.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh mengakses data mahasiswa lain"})
		}
	}

	// dosen hanya boleh lihat bimbingannya
	if role == "Dosen Wali" {
		lect, err := repository.GetLecturerByUserID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Data dosen tidak ditemukan",
			})
		}

		if stud.AdvisorID != lect.ID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Mahasiswa bukan bimbingan Anda",
			})
		}
	}

	return c.JSON(fiber.Map{"success": true, "data": stud})
}

// PUT /students/:id/advisor  (Admin only)
func StudentSetAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")
	var body struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input invalid"})
	}

	err := repository.SetStudentAdvisor(studentID, body.AdvisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengupdate dosen wali"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Advisor updated"})
}

// GET /students/:id/achievements
func StudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	// ambil record mahasiswa berdasarkan studentID yang dikirim
	targetStudent, err := repository.GetStudentByID(studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Mahasiswa tidak ditemukan"})
	}

	// 🔐 ROLE VALIDATION
	switch role {
	case "Mahasiswa":
		// mahasiswa hanya boleh akses prestasinya sendiri
		self, _ := repository.GetStudentByUserID(userID)
		if self.ID != targetStudent.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh mengakses prestasi mahasiswa lain"})
		}

	case "Dosen Wali":
		// dosen hanya boleh akses prestasi mahasiswa bimbingannya
		lect, _ := repository.GetLecturerByUserID(userID)
		if targetStudent.AdvisorID != lect.ID {
			return c.Status(403).JSON(fiber.Map{"error": "Mahasiswa bukan bimbingan Anda"})
		}

	case "Admin":
		// admin boleh akses semua → tidak perlu validasi
	default:
		return c.Status(403).JSON(fiber.Map{"error": "Role tidak dikenali"})
	}

	// ambil achievement reference milik mahasiswa
	refs, err := repository.GetAchievementReferencesByStudentID(targetStudent.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil reference achievement"})
	}

	// join ke Mongo
	var results []fiber.Map
	for _, ref := range refs {
		doc, err := repository.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}
		results = append(results, fiber.Map{
			"reference":   ref,
			"achievement": doc,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(results),
		"data":    results,
	})
}