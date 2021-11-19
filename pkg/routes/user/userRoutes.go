package userRoutes

import (
	"github.com/gofiber/fiber/v2"
	userHandler "vcsxsantos/nagini-api/pkg/handlers/user"
	constants "vcsxsantos/nagini-api/pkg/util/constant"
	"vcsxsantos/nagini-api/pkg/util/jwt"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")

	user.Get("/", jwt.Protected(), userHandler.GetUsers)

	user.Get(constants.PathUserIdParam, jwt.Protected(), userHandler.GetUser)

	user.Put(constants.PathUserIdParam, jwt.Protected(), userHandler.UpdateUser)

	user.Delete(constants.PathUserIdParam, jwt.Protected(), userHandler.DeleteUser)

}
