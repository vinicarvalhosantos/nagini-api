package v1

import (
	"github.com/gofiber/fiber/v2"
	addressRoutes "github.com/vinicius.csantos/nagini-api/internal/route/address"
	userRoutes "github.com/vinicius.csantos/nagini-api/internal/route/user"
)

func SetupV1Routes(router fiber.Router) {

	api := router.Group("/v1")

	//User Routes
	userRoutes.SetupUserRoutes(api)

	//Address Routes
	addressRoutes.SetupAddressRoutes(api)

}
