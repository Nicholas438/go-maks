package handler

import (
	"data_service/database"
	"data_service/model/entity"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func DataHandlerGetAll(ctx *fiber.Ctx) error {
	var trades []entity.Trades
	result := database.DB.Raw("SELECT * FROM trades").Scan(&trades)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    trades,
	})
}

func DataHandlerGetByCoinId(ctx *fiber.Ctx) error {
	coinID := ctx.Params("coin_id")

	var trades []entity.Trades

	result := database.DB.Raw("SELECT * FROM trades WHERE coin_id = ?", coinID).Scan(&trades)
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

func GenerateAndStoreRandomData() {
	var coins []entity.Coin
	err := database.DB.Find(&coins).Error
	if err != nil {
		log.Println(err)
	}
	fmt.Println(coins)

	for _, coin := range coins {
		price := rand.Float64()*1000 + 1000

		trade := entity.Trades{
			Price:  price,
			UserID: 0,
			CoinID: coin.ID,
		}

		errGenerateTrade := database.DB.Create(&trade).Error
		if errGenerateTrade != nil {
			log.Println("Failed to store random trade:", errGenerateTrade)
		}
	}

}

func GetLowestData() {
	timeStr, err := database.Rdb.Get(database.Ctx, "lowest_price_24_hrs_time").Result()
	if err == redis.Nil {
		SetLowestData()
	}
	if err != nil {
		log.Println("Failed to get Redis key:", err)
	}
	fmt.Println(timeStr)

	timeRes, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Println("Invalid time value:", err)
	}
	if time.Since(timeRes) > 24*time.Hour {
		SetLowestData()
	}

}

func SetLowestData() {
	var lowest entity.Trades
	err := database.DB.Raw("SELECT price, created_at FROM trades WHERE created_at >= NOW() - INTERVAL '24 HOURS' ORDER BY price ASC LIMIT 1;").Scan(&lowest)
	if err != nil {
		log.Println(err)
	}

	errRedis := database.Rdb.Set(database.Ctx, "lowest_price_24_hrs", lowest.Price, 0).Err()
	if errRedis != nil {
		log.Println("Failed to store lowest price:", errRedis)
	}

	errRedis = database.Rdb.Set(database.Ctx, "lowest_price_24_hrs_time", lowest.CreatedAt.Format(time.RFC3339), 0).Err()
	if errRedis != nil {
		log.Println("Failed to store timestamp:", errRedis)
	}

	log.Println("Lowest price and time stored in Redis")
}
