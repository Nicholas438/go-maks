package route

import (
	"maks-go/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/", handler.UserHandler)
}
