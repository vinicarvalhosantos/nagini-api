package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"vcsxsantos/nagini-api/Internal/model"
	"vcsxsantos/nagini-api/config"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	portString := config.Config("DB_PORT", "5432")
	if len(portString) == 0 {
		portString = "5432"
	}
	dbPort, err := strconv.ParseUint(portString, 10, 32)

	if err != nil {
		log.Println("Port it is not a number!")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo", config.Config("DB_HOST", "localhost"),
		config.Config("DB_USERNAME", "nagini-api"), config.Config("DB_PASSWORD", "nagini-api"), config.Config("DB_NAME", "nagini-api"), dbPort)

	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic("Failed to migrate models")
	}
	fmt.Println("Database migrated")

}
