package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"vcsxsantos/nagini-api/pkg/routes/login"
	"vcsxsantos/nagini-api/pkg/routes/register"
	roleRoutes "vcsxsantos/nagini-api/pkg/routes/role"
	userRoutes "vcsxsantos/nagini-api/pkg/routes/user"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1", logger.New())

	userRoutes.SetupUserRoutes(api)
	roleRoutes.SetupRoleRoutes(api)
	login.SetupLoginRoutes(api)
	register.SetupRegisterRoutes(api)

}
