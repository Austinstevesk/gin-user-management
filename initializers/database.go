package initializers

import (
	"fmt"
	"gin-user-management/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
}

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	fmt.Println("Migration Complete")
}