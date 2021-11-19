package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"vcsxsantos/nagini-api/config"
	"vcsxsantos/nagini-api/pkg/model"
)

var DB *gorm.DB

const DB_DEFAULT = "nagini-api"

func ConnectDB() {
	var err error
	portString := config.Config("DB_PORT", "5432")

	dbPort, err := strconv.ParseUint(portString, 10, 32)

	if err != nil {
		log.Println("Port it is not a number!")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/Sao_Paulo", config.Config("DB_HOST", "localhost"),
		config.Config("DB_USERNAME", DB_DEFAULT), config.Config("DB_PASSWORD", DB_DEFAULT), config.Config("DB_NAME", DB_DEFAULT), dbPort)

	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Connection Opened to Database")

	err = DB.AutoMigrate(&model.Role{})
	err = DB.AutoMigrate(&model.User{})

	if err != nil {
		panic("Failed to migrate models")
	}
	fmt.Println("Database migrated")

}
