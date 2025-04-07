package handler

import (
	"data_service/database"
	"data_service/model/entity"
	"fmt"
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

func DataHandlerGetAll(ctx *fiber.Ctx) error {
	var trades []entity.Trades
	result := database.DB.Raw("SELECT * FROM trades").Scan(&trades)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return ctx.JSON(
		trades,
	)
}

func GenerateAndStoreRandomData() {
	var coins []entity.Coin
	err := database.DB.Find(&coins).Error
	if err != nil {
		log.Println(err)
	}
	fmt.Println(coins)

	for _, coin := range coins {
		price := rand.Intn(1000) + 1000

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
