package handler

import (
	"trade_service/database"
	"trade_service/model/entity"
	"trade_service/model/request"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func TradeHandlerCreate(ctx *fiber.Ctx) error {
	trade := new(request.TradeCreateRequest)
	if err := ctx.BodyParser(trade); err != nil {
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(trade)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed Trade ",
			"error":   errValidate.Error(),
		})
	}

	userIDInterface := ctx.Locals("user_id")
	userID, ok := userIDInterface.(int)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid user ID in context",
		})
	}

	newTrade := entity.Trades{
		Price:  trade.Price,
		CoinID: trade.CoinID,
		UserID: userID,
	}

	errCreateTrade := database.DB.Create(&newTrade).Error
	if errCreateTrade != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "failed to create trade",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newTrade,
	})
}
