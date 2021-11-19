package login

import (
	"github.com/gofiber/fiber/v2"
	loginHandler "vcsxsantos/nagini-api/pkg/handlers/login"
)

func SetupLoginRoutes(router fiber.Router) {
	login := router.Group("/login")

	login.Post("/", loginHandler.Login)

}
