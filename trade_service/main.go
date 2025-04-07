package main

import (
	"trade_service/database"
	"trade_service/database/migration"
	"trade_service/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.DatabaseInit()
	migration.RunMigration()

	app := fiber.New()

	route.RouteInit(app)

	app.Listen(":3002")
}
