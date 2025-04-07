package main

import (
	"maks-go/config"
	"maks-go/database"
	"maks-go/database/migration"
	"maks-go/route"

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
