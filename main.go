package main

import (
	"log"

	"prestasi_backend/config"
	"prestasi_backend/database"
	"prestasi_backend/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Environment
	config.LoadEnv()

	// 2. Connect PostgreSQL
	postgresDB, err := database.ConnectPostgre()
	if err != nil {
		log.Fatal(err)
	}
	database.DB = postgresDB

	// 3. Connect MongoDB
	mongoDB, err := database.ConnectMongo()
	if err != nil {
		log.Fatal(err)
	}
	database.MongoDB = mongoDB

	// Start Fiber server
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API Ready!",
		})
	})

	// 4. REGISTER ROUTES
	route.SetupRoutes(app)

	port := config.Get("APP_PORT")
	log.Println("🚀 Server running on port", port)

	app.Listen(":" + port)

}