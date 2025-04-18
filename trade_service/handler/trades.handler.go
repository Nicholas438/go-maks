package handler

import (
	"log"
	"strconv"
	"time"
	"trade_service/database"
	"trade_service/model/entity"
	"trade_service/model/request"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func TradesHandlerGet(ctx *fiber.Ctx) error {
	coinID := ctx.Query("coin_id")
	var trades []entity.Trades
	var result *gorm.DB

	if coinID == "" {
		result = database.DB.Raw("SELECT * FROM trades").Scan(&trades)

	} else {
		result = database.DB.Raw("SELECT * FROM trades WHERE coin_id = ?", coinID).Scan(&trades)
	}

	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Trades query failed",
			"error":   result.Error,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    trades,
	})
}

func TradeHandlerCreate(ctx *fiber.Ctx) error {
	var trade request.TradeCreateRequest
	if err := ctx.BodyParser(&trade); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(trade); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	lowestPriceRedis := GetLowestData()
	if trade.Price < (0.5 * lowestPriceRedis) {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Cannot trade lower than half of lowest trade price in the past 24 hours",
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
	if trade.Price < lowestPriceRedis {
		SetLowestData(trade.Price, time.Now().Format(time.RFC3339))
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newTrade,
	})
}

func GetLowestData() float64 {
	price, err := database.Rdb.Get(database.Ctx, "lowest_price_24_hrs").Result()
	if err == redis.Nil {
		log.Println("No trade info")
		return 0
	}
	if err != nil {
		log.Println("Failed to get Redis key:", err)
	}

	fprice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println("Conversion failed:", err)
	}
	return fprice

}

func SetLowestData(price float64, createdAt string) {
	errRedis := database.Rdb.Set(database.Ctx, "lowest_price_24_hrs", price, 0).Err()
	if errRedis != nil {
		log.Println("Failed to store lowest price:", errRedis)
	}

	errRedis = database.Rdb.Set(database.Ctx, "lowest_price_24_hrs_time", createdAt, 0).Err()
	if errRedis != nil {
		log.Println("Failed to store timestamp:", errRedis)
	}

	log.Println("Lowest price and time stored in Redis")
}
