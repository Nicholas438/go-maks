package migration

import (
	"fmt"
	"log"
	"trade_service/database"
	"trade_service/model/entity"
)

func RunMigration() {
	err := database.DB.AutoMigrate(&entity.Trades{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database migrated")
}
