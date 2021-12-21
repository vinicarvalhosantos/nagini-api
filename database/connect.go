package database

import (
	"fmt"
	"gitlab.com/vinicius.csantos/nagini-api/config"
	"gitlab.com/vinicius.csantos/nagini-api/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var DB *gorm.DB

const DbDefault = "nagini-api"

func ConnectDB() {
	var err error
	portString := config.Config("DB_PORT", "3306")

	dbPort, err := strconv.ParseUint(portString, 10, 32)

	dbHost := config.Config("DB_HOST", "localhost")
	dbUser := config.Config("DB_USERNAME", DbDefault)
	dbPass := config.Config("DB_PASSWORD", DbDefault)
	dbName := config.Config("DB_NAME", DbDefault)

	if err != nil {
		log.Println("Port it is not a number!")
	}

	//textToToken := "$viniciuscarvalhomine@gmail.com$zvinniie$44611032850$" + time.Now().Add(time.Hour * 1).String()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Connection Opened to Database")

	fmt.Println("Migrating database")

	err = migrateModels()

	if err != nil {
		log.Panicln("Failed to migrate models: " + err.Error())
	}

	fmt.Println("Database migrated")

}

func migrateModels() error {
	err := DB.AutoMigrate(&model.Address{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}
