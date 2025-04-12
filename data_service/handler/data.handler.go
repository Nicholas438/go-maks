package handler

import (
	"data_service/database"
	"data_service/model/entity"
	"data_service/model/request"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func LowestTradeHandler(ctx *fiber.Ctx) error {
	price, time := GetLowestData()
	return ctx.Status(200).JSON(fiber.Map{
		"price":         price,
		"time_recorded": time,
	})
}

func AveragePriceHandler(ctx *fiber.Ctx) error {
	coinID := ctx.Query("coin_id")
	var price float64
	var result *gorm.DB

	if coinID == "" {
		result = database.DB.Raw("SELECT AVG(price) FROM trades WHERE created_at >= NOW() - INTERVAL '24 HOURS'").Scan(&price)

	} else {
		result = database.DB.Raw("SELECT AVG(price) FROM trades WHERE created_at >= NOW() - INTERVAL '24 HOURS' AND coin_id=?", coinID).Scan(&price)
	}

	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Trades query failed",
			"error":   result.Error,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    price,
	})
}

func CoinHandlerCreate(ctx *fiber.Ctx) error {
	coin := new(request.CoinCreateRequest)
	if err := ctx.BodyParser(coin); err != nil {
		return err
	}

	validate := validator.New()
	errValidate := validate.Struct(coin)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	newCoin := entity.Coin{
		Name: coin.Name,
	}

	errCreateCoin := database.DB.Create(&newCoin).Error
	if errCreateCoin != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Failed to store data",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    newCoin,
	})
}

func RSICalculator(ctx *fiber.Ctx) error {
	coinIDQuery := ctx.Query("coin_id")
	periodQuery := ctx.Query("period", "14")
	var price []float64
	var result *gorm.DB

	coinID, err := strconv.Atoi(coinIDQuery)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid coin_id"})
	}
	period, err := strconv.Atoi(periodQuery)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid period"})
	}

	result = database.DB.Raw("SELECT price FROM trades WHERE coin_id = ? ORDER BY created_at ASC", coinID).Scan(&price)

	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Trades query failed",
			"error":   result.Error,
		})
	}

	rsi, err := CalculateRSI(price, period)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "RSI Calculation failed"})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"rsi":     rsi[len(rsi)-1],
		"coin_id": coinID,
		"period":  period,
	})
}

func CalculateRSI(prices []float64, period int) ([]float64, error) {
	if len(prices) <= period {
		return nil, fmt.Errorf("not enough data to calculate RSI (need more than %d points)", period)
	}

	rsi := make([]float64, len(prices))
	gains := 0.0
	losses := 0.0

	// Initial average gain/loss
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change // subtract negative to get positive loss
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// First RSI value
	rs := avgGain / avgLoss
	rsi[period] = 100 - (100 / (1 + rs))

	// Continue smoothing RSI
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		var gain, loss float64

		if change > 0 {
			gain = change
			loss = 0
		} else {
			gain = 0
			loss = -change
		}

		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs = avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi, nil
}

func CoinHandlerGet(ctx *fiber.Ctx) error {
	var coins []entity.Coin

	result := database.DB.Raw("SELECT * FROM coins").Scan(&coins)

	if result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "Coin query failed",
			"error":   result.Error,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Success",
		"data":    coins,
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

		lowestPriceRedis, timeStrRedis := GetLowestData()
		if trade.Price < lowestPriceRedis {
			SetLowestData(trade.Price, time.Now().Format(time.RFC3339))
		}

		parsedTime, err := time.Parse(time.RFC3339, timeStrRedis)
		if err != nil {
			log.Println("Invalid time format in Redis:", err)
			UpdateLowestData()
		}
		if time.Since(parsedTime) > 24*time.Hour {
			log.Println("Time passed while generating data")
			FindLowestData()
		}

		errGenerateTrade := database.DB.Create(&trade).Error
		if errGenerateTrade != nil {
			log.Println("Failed to store random trade:", errGenerateTrade)
		}
	}

}

func UpdateLowestData() {
	timeStr, err := database.Rdb.Get(database.Ctx, "lowest_price_24_hrs_time").Result()
	if err == redis.Nil {
		log.Println("No data")
		FindLowestData()
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
		log.Println("Time passed while updating data")
		FindLowestData()
	}

}

func FindLowestData() {
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

func GetLowestData() (float64, string) {
	price, err := database.Rdb.Get(database.Ctx, "lowest_price_24_hrs").Result()
	if err == redis.Nil {
		log.Println("No trade info")
		return 0, ""
	}
	if err != nil {
		log.Println("Failed to get Redis key:", err)
	}

	timeRes, err := database.Rdb.Get(database.Ctx, "lowest_price_24_hrs_time").Result()
	if err != nil {
		log.Println("Failed to get Redis key:", err)
	}

	fprice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println("Conversion failed:", err)
	}
	return fprice, timeRes

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
