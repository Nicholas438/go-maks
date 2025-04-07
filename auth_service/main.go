package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/database/migration"
	"auth_service/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.DatabaseInit()
	migration.RunMigration()

	app := fiber.New()
	config.GoogleConfig()

	route.RouteInit(app)

	app.Listen(":3000")
}
