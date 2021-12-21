package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/vinicius.csantos/nagini-api/config"
	"github.com/vinicius.csantos/nagini-api/database"
	"github.com/vinicius.csantos/nagini-api/router"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	appPort := config.Config("APPLICATION_PORT", "5000")

	database.ConnectDB()

	router.SetupRoutes(app)

	err := app.Listen(":" + appPort)

	if err != nil {
		log.Panicln(err.Error())
	}
}
