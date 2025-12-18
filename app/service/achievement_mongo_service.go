package service

import (
	"time"

	"prestasi_backend/app/model"
	"prestasi_backend/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// ==================================================================
// GET /api/v1/achievements  — list all (role filtered)
// ==================================================================
func AchievementList(c *fiber.Ctx) error {

	role := c.Locals("role").(string)
	userID := c.Locals("userId").(string)

	// ================================================================
	// ADMIN
	// ================================================================
	if role == "Admin" {
		refs, err := repository.GetAllAchievementReferences()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil achievement"})
		}

		return buildAchievementResponse(c, refs)
	}

	// ================================================================
	// DOSEN WALI
	// ================================================================
	if role == "Dosen Wali" {
		lecturer, err := repository.GetLecturerByUserID(userID)
		if err != nil {
			return c.JSON(fiber.Map{"success": true, "count": 0, "data": []any{}})
		}

		refs, err := repository.GetAchievementReferencesByAdvisor(lecturer.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil achievement dosen wali"})
		}

		// ⬅ RETURN DI SINI (WAJIB)
		return buildAchievementResponse(c, refs)
	}

	// ================================================================
	// MAHASISWA
	// ================================================================
	if role == "Mahasiswa" {
		student, err := repository.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Data mahasiswa tidak ditemukan"})
		}

		refs, err := repository.GetAchievementReferencesByStudentID(student.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil achievement mahasiswa"})
		}

		return buildAchievementResponse(c, refs)
	}

	return c.Status(403).JSON(fiber.Map{"error": "Role tidak dikenali"})
}

// Join Postgre + Mongo lalu format JSON response
func buildAchievementResponse(c *fiber.Ctx, refs []model.AchievementReference) error {
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

// ==================================================================
// GET /api/v1/achievements/:id — detail achievement
// ==================================================================
func AchievementDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	ref, err := repository.GetAchievementReferenceByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Achievement tidak ditemukan",
		})
	}

	doc, err := repository.GetAchievementByID(ref.MongoAchievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data achievement (MongoDB)",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reference":   ref,
			"achievement": doc,
		},
	})
}

// ==================================================================
// POST /api/v1/achievements  — create achievement
// ==================================================================
func AchievementCreate(c *fiber.Ctx) error {

	var req model.AchievementCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "student_id wajib diisi"})
	}

	// Insert ke MongoDB
	doc := model.AchievementMongo{
		StudentID:       req.StudentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}

	mongoID, err := repository.CreateAchievement(&doc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan achievement ke MongoDB"})
	}

	// Insert ke Postgre
	ref := model.AchievementReference{
		ID:                 uuid.NewString(),
		StudentID:          req.StudentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
	}

	if err := repository.CreateAchievementReference(&ref); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan reference achievement"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement berhasil dibuat",
		"data": fiber.Map{
			"reference":   ref,
			"achievement": doc,
		},
	})
}

// ==================================================================
// PUT /api/v1/achievements/:id — update achievement (only draft)
// ==================================================================
func AchievementUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req model.AchievementUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	ref, err := repository.GetAchievementReferenceByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement tidak ditemukan"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "Hanya achievement draft yang dapat diupdate"})
	}

	update := bson.M{}

	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Details != nil {
		update["details"] = req.Details
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}
	if req.Points != 0 {
		update["points"] = req.Points
	}

	if err := repository.UpdateAchievement(ref.MongoAchievementID, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update achievement"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement berhasil diupdate",
	})
}

// ==================================================================
// DELETE /api/v1/achievements/:id — soft delete
// ==================================================================
func AchievementDelete(c *fiber.Ctx) error {
    refID := c.Params("id")
    userID := c.Locals("userId").(string)
    role := c.Locals("role").(string)

    // Get reference
    ref, err := repository.GetAchievementReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data achievement tidak ditemukan"})
    }

    // ========= ROLE VALIDATION =========
    switch role {
    case "Admin":
        // Admin bebas — tidak dicek status

    case "Mahasiswa":
        student, _ := repository.GetStudentByUserID(userID)
        if student == nil || student.ID != ref.StudentID {
            return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh menghapus prestasi orang lain"})
        }

        if ref.Status != "draft" {
            return c.Status(403).JSON(fiber.Map{
                "error": "Mahasiswa hanya boleh menghapus prestasi berstatus draft",
            })
        }

    default:
        // Dosen atau role lain tidak boleh hapus
        return c.Status(403).JSON(fiber.Map{"error": "Role ini tidak diizinkan melakukan hapus"})
    }

    // ============ Soft delete Postgre ============
    err = repository.SoftDeleteAchievementReference(refID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus achievement"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Prestasi berhasil dihapus",
    })
}

// ==================================================================
// POST /api/v1/achievements/:id/submit — submit by mahasiswa
// ==================================================================
func AchievementSubmit(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("userId").(string)

	// ambil mahasiswa via user id
	student, err := repository.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Hanya mahasiswa yang dapat submit",
		})
	}

	// ambil reference prestasi
	ref, err := repository.GetAchievementReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement tidak ditemukan"})
	}

	// validasi kepemilikan
	if ref.StudentID != student.ID {
		return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh submit achievement milik orang lain"})
	}

	// validasi status
	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "Hanya status draft yang bisa submit"})
	}

	// update status → submitted
	err = repository.SubmitAchievementReference(achievementID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal submit achievement"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement berhasil dikirim untuk diverifikasi",
	})
}

// ==================================================================
// POST /api/v1/achievements/:id/verify — verify by dosen wali
// ==================================================================
func AchievementVerify(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("userId").(string)

	// ambil dosen via user id
	lecturer, err := repository.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya dosen wali yang dapat memverifikasi"})
	}

	// ambil reference
	ref, err := repository.GetAchievementReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement tidak ditemukan"})
	}

	// cek apakah mahasiswa tersebut memang bimbingannya
	ok, err := repository.IsStudentAdvisedBy(lecturer.ID, ref.StudentID)
	if err != nil || !ok {
		return c.Status(403).JSON(fiber.Map{"error": "Mahasiswa bukan bimbingan anda"})
	}

	// update → verified
	err = repository.VerifyAchievementReference(achievementID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memverifikasi achievement"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement berhasil diverifikasi",
	})
}

// ==================================================================
// POST /api/v1/achievements/:id/reject — reject by dosen wali
// ==================================================================
func AchievementReject(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("userId").(string)

	var req model.AchievementRejectRequest
	if err := c.BodyParser(&req); err != nil || req.Note == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Alasan penolakan wajib diisi"})
	}

	lecturer, err := repository.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Hanya dosen wali yang dapat menolak"})
	}

	ref, err := repository.GetAchievementReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement tidak ditemukan"})
	}

	ok, err := repository.IsStudentAdvisedBy(lecturer.ID, ref.StudentID)
	if err != nil || !ok {
		return c.Status(403).JSON(fiber.Map{"error": "Mahasiswa bukan bimbingan anda"})
	}

	err = repository.RejectAchievementReference(achievementID, userID, req.Note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menolak achievement"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement berhasil ditolak",
	})
}

// ==================================================================
// GET /api/v1/achievements/:id/history — lihat status
// ==================================================================
func AchievementHistory(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	ref, err := repository.GetAchievementReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Achievement tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status":         ref.Status,
			"submitted_at":   ref.SubmittedAt,
			"verified_at":    ref.VerifiedAt,
			"verified_by":    ref.VerifiedBy,
			"rejection_note": ref.RejectionNote,
		},
	})
}

// ==================================================================
// POST /api/v1/achievements/:id/attachments — upload file
// ==================================================================
func AchievementUploadAttachment(c *fiber.Ctx) error {
    refID := c.Params("id")
    userID := c.Locals("userId").(string)
    role := c.Locals("role").(string)

    // Ambil reference dari Postgre
    ref, err := repository.GetAchievementReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data achievement tidak ditemukan"})
    }

    // Cek hak upload
    switch role {
    case "Admin":
        // boleh semua
    case "Mahasiswa":
        student, _ := repository.GetStudentByUserID(userID)
        if student == nil || student.ID != ref.StudentID {
            return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh mengupload file ke prestasi orang lain"})
        }
    default:
        return c.Status(403).JSON(fiber.Map{"error": "Role ini tidak diperbolehkan upload attachment"})
    }

    // Ambil file upload
    f, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "File tidak ditemukan"})
    }

    filename := time.Now().Format("20060102_150405") + "_" + f.Filename
    path := "./uploads/" + filename

    if err := c.SaveFile(f, path); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file"})
    }

    // Tambah attachment ke dokumen Mongo
    err = repository.AddAchievementAttachment(ref.MongoAchievementID, filename)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan attachment"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Attachment berhasil ditambahkan",
        "filename": filename,
    })
}