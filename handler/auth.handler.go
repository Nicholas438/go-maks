package handler

import "github.com/gofiber/fiber/v2"

func UserHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "Hello, World!",
	})

}
