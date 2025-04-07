package handler

import (
	"data_service/database"
	"data_service/model/entity"
	"log"

	"github.com/gofiber/fiber/v2"
)

func DataHandlerGetAll(ctx *fiber.Ctx) error {
	tradesInfo := ctx.Locals("tradesInfo")
	log.Println("tradesInfo data::", tradesInfo)
	var trades []entity.Trades
	result := database.DB.Find(&trades)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return ctx.JSON(
		trades,
	)
}
