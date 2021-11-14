package main

import (
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/config"
	"vcsxsantos/nagini-api/database"
	"vcsxsantos/nagini-api/router"
)

func main() {
	app := fiber.New()
	appPort := config.Config("APPLICATION_PORT", "8000")

	database.ConnectDB()

	router.SetupRoutes(app)

	app.Listen(":" + appPort)
}
