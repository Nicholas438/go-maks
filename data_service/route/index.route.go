package route

import (
	"data_service/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/data/lowest-trade", handler.GetLowestTrade)
	r.Post("/coin/create", handler.CoinHandlerCreate)
}
