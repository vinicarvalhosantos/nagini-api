package register

import (
	"github.com/gofiber/fiber/v2"
	userHandler "vcsxsantos/nagini-api/pkg/handlers/user"
)

func SetupRegisterRoutes(router fiber.Router) {

	register := router.Group("/user/register")

	register.Post("/", userHandler.RegisterUser)

}
