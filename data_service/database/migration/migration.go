package migration

import (
	"auth_service/database"
	"auth_service/model/entity"
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
