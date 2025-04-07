package route

import (
	"trade_service/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Post("/trade", handler.TradeCoinHandler)
}
