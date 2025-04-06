package migration

import (
	"fmt"
	"log"
	"maks-go/database"
	"maks-go/model/entity"
)

func RunMigration() {
	err := database.DB.AutoMigrate(&entity.User{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database migrated")
}
