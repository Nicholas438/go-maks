package route

import (
	"auth_service/handler"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Get("/google-login", handler.GoogleLogin)
	r.Get("/google-callback", handler.GoogleCallback)
	r.Post("/auth", handler.AuthorizeUser)
	r.Post("/login", handler.LoginAuth)
	r.Post("/register", handler.RegisterAuth)
	r.Static("/", "./public")
}
