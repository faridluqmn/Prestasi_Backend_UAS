package main

import (
	"log"

	"prestasi_backend/config"
	"prestasi_backend/database"
	"prestasi_backend/route"

	"github.com/gofiber/fiber/v2"
)

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API Ready!",
		})
	})

	route.SetupRoutes(app)

	port := config.Get("APP_PORT")
	log.Println("🚀 Server running on port", port)

	app.Listen(":" + port)

}