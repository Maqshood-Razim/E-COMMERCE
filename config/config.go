package config

import (
	"furniture-ecommerce/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:razeem19@tcp(127.0.0.1:3306)/db7"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate all models in the correct order
	err = DB.AutoMigrate(
		&models.OrderDetails{},
		&models.User{},
		&models.Cart{},
		&models.Product{},
		&models.Checkout{},
		&models.Payment{},
	)
	if err != nil {
		log.Printf("failed to auto-migrate database: %v", err)
		return
	}

	log.Println("Database connection and migration successful")
}
