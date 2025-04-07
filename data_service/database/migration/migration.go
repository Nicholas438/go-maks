package migration

import (
	"data_service/database"
	"data_service/model/entity"
	"fmt"
	"log"
)

func RunMigration() {
	err := database.DB.AutoMigrate(&entity.Trades{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database migrated")
}
