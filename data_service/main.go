package main

import (
	"data_service/database"
	"data_service/database/migration"
	"data_service/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.DatabaseInit()
	migration.RunMigration()

	app := fiber.New()

	route.RouteInit(app)

	app.Listen(":3001")
}
