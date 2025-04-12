package route

import (
	"data_service/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/data/lowest-trade", handler.LowestTradeHandler)
	r.Get("/data/average-price", handler.AveragePriceHandler)
	r.Get("/data/coins", handler.CoinHandlerGet)
	r.Get("/data/rsi", handler.RSICalculator)
	r.Post("/coin/create", handler.CoinHandlerCreate)
}
