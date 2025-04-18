package route

import (
	"trade_service/handler"
	"trade_service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/trade", middleware.AuthMiddleware(), handler.TradesHandlerGet)

	r.Post("/trade", middleware.AuthMiddleware(), handler.TradeHandlerCreate)
}
