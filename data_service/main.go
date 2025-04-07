package main

import (
	"data_service/database"
	"data_service/database/migration"
	"data_service/handler"
	"data_service/route"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.DatabaseInit()
	migration.RunMigration()

	go startDataGeneration()

	app := fiber.New()

	route.RouteInit(app)

	app.Listen(":3001")
}

func startDataGeneration() {
	handler.GetLowestData()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			handler.GenerateAndStoreRandomData()
			handler.GetLowestData()
		}
	}
}
