package main

import (
	"log"

	"prestasi_backend/config"
	"prestasi_backend/database"
	"prestasi_backend/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "prestasi_backend/docs"
)

// @title           Sistem Prestasi Mahasiswa API
// @description     Dokumentasi Backend UAS. <br> **Fitur Utama:** <br> 1. Authentication (PostgreSQL) <br> 2. Manajemen Prestasi Dinamis (MongoDB).
// @termsOfService  http://swagger.io/terms/

// @contact.name    Muhammad Farid Luqman Hakim - 434231058
// @contact.email   muhammad.farid.luqman.410762-2023@vokasi.unair.ac.id

// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	config.LoadEnv()

	postgresDB, err := database.ConnectPostgre()
	if err != nil {
		log.Fatal(err)
	}
	database.DB = postgresDB

	mongoDB, err := database.ConnectMongo()
	if err != nil {
		log.Fatal(err)
	}
	database.MongoDB = mongoDB

	app := fiber.New()

	// Route Cek Health
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API Ready!",
		})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	route.SetupRoutes(app)

	port := config.Get("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)
	app.Listen(":" + port)
}