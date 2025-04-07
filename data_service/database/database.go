package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Rdb *redis.Client
var Ctx = context.Background()

func DatabaseInit() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_DB"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Aiven PostgreSQL: %v", err)
	}

	log.Println("Connected to Aiven PostgreSQL successfully!")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("VALKEY_ADDR"),
		Password: os.Getenv("VALKEY_PASSWORD"),
		DB:       0,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	})

	// Test Redis connection
	pong, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	fmt.Println("Connected to Valkey:", pong)
}
