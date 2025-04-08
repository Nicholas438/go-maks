package route

import (
	"data_service/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/bulk-trades-read", handler.DataHandlerGetAll)
	r.Get("/trades-filter-coin-id/:coin_id", handler.DataHandlerGetByCoinId)
}
