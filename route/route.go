package route

import (
    "prestasi_backend/middleware"
    "prestasi_backend/app/service"

	"fmt"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	fmt.Println("🔥 ROUTES LOADED")

	api := app.Group("/api/v1")

	// ============================================
	// 5.1 AUTHENTICATION
	// ============================================
	api.Post("/auth/login", service.AuthLogin)
	api.Post("/auth/refresh", middleware.JWTRequired(), service.AuthRefresh)
	api.Post("/auth/logout", middleware.JWTRequired(), service.AuthLogout)
	api.Get("/auth/profile", middleware.JWTRequired(), service.AuthProfile)

	// ============================================
	// 5.2 USERS (Admin Only)
	// ============================================
	users := api.Group("/users", middleware.JWTRequired())

	users.Get("/", middleware.PermissionRequired("user:manage"), service.UserList)
	users.Get("/:id", middleware.PermissionRequired("user:manage"), service.UserDetail)
	users.Post("/", middleware.PermissionRequired("user:manage"), service.UserCreate)
	users.Put("/:id", middleware.PermissionRequired("user:manage"), service.UserUpdate)
	users.Delete("/:id", middleware.PermissionRequired("user:manage"), service.UserDelete)
	users.Put("/:id/role", middleware.PermissionRequired("user:manage"), service.UserUpdateRole)

	// ============================================
	// 5.4 ACHIEVEMENTS
	// ============================================
	ach := api.Group("/achievements", middleware.JWTRequired())

	ach.Get("/", service.AchievementList)
	ach.Get("/:id", service.AchievementDetail)

	ach.Post("/", middleware.PermissionRequired("achievement:create"), service.AchievementCreate)
	ach.Put("/:id", middleware.PermissionRequired("achievement:update"), service.AchievementUpdate)
	ach.Delete("/:id", middleware.PermissionRequired("achievement:delete"), service.AchievementDelete)

	ach.Post("/:id/submit", middleware.PermissionRequired("achievement:update"), service.AchievementSubmit)
	ach.Post("/:id/verify", middleware.PermissionRequired("achievement:verify"), service.AchievementVerify)
	ach.Post("/:id/reject", middleware.PermissionRequired("achievement:verify"), service.AchievementReject)

	ach.Get("/:id/history", service.AchievementHistory)
	ach.Post("/:id/attachments", middleware.PermissionRequired("achievement:update"), service.AchievementUploadAttachment)

	// ============================================
	// 5.5 STUDENTS
	// ============================================
	students := api.Group("/students", middleware.JWTRequired())

	students.Get("/", service.StudentList)
	students.Get("/:id", service.StudentDetail)
	students.Get("/:id/achievements", service.StudentAchievements)
	students.Put("/:id/advisor", middleware.PermissionRequired("user:manage"), service.StudentSetAdvisor)

	// ============================================
	// 5.5 LECTURERS
	// ============================================
	lect := api.Group("/lecturers", middleware.JWTRequired())

	lect.Get("/", service.LecturerList)
	lect.Get("/:id/advisees", service.LecturerAdvisees)

	// ============================================
	// 5.8 REPORTS
	// ============================================
	reports := api.Group("/reports", middleware.JWTRequired())

	reports.Get("/statistics", service.ReportStatistics)
	reports.Get("/student/:id", service.ReportStudent)
}
