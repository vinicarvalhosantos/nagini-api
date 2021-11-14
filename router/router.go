package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	userRoutes "vcsxsantos/nagini-api/Internal/routes/user"
)

func SetupRoutes(app *fiber.App){
	api := app.Group("/api", logger.New())

	userRoutes.SetupUserRoutes(api)


}
