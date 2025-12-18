package service

import (
	"github.com/gofiber/fiber/v2"
	"prestasi_backend/app/repository"
)

// ==================================================================
// LIST STUDENTS
// ==================================================================

// StudentList godoc
// @Summary      Lihat Daftar Mahasiswa
// @Description  Admin melihat semua mahasiswa. Dosen Wali hanya melihat mahasiswa bimbingannya.
// @Tags         Student
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /students [get]
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

// ==================================================================
// STUDENT DETAIL
// ==================================================================

// StudentDetail godoc
// @Summary      Detail Mahasiswa
// @Description  Melihat detail data satu mahasiswa. Mahasiswa hanya bisa lihat diri sendiri. Dosen hanya bimbingannya.
// @Tags         Student
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Student ID (UUID)"
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /students/{id} [get]
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

// ==================================================================
// SET ADVISOR (ADMIN ONLY)
// ==================================================================

// StudentSetAdvisor godoc
// @Summary      Set Dosen Wali (Admin)
// @Description  Admin menetapkan dosen wali untuk mahasiswa tertentu.
// @Tags         Student
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string             true  "Student ID"
// @Param        request  body      map[string]string  true  "Body: { advisor_id: 'uuid-dosen' }"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /students/{id}/advisor [put]
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

// ==================================================================
// STUDENT ACHIEVEMENTS LIST
// ==================================================================

// StudentAchievements godoc
// @Summary      Lihat Prestasi Mahasiswa Tertentu
// @Description  Melihat daftar prestasi milik mahasiswa tertentu berdasarkan ID-nya (Admin, Dosen Wali, & Pemilik Akun).
// @Tags         Student
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Student ID"
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /students/{id}/achievements [get]
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