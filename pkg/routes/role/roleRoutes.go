package roleRoutes

import (
	"github.com/gofiber/fiber/v2"
	roleHandler "vcsxsantos/nagini-api/pkg/handlers/role"
	constants "vcsxsantos/nagini-api/pkg/util/constant"
	"vcsxsantos/nagini-api/pkg/util/jwt"
)

func SetupRoleRoutes(router fiber.Router) {
	role := router.Group("/role")

	role.Post("/", jwt.Protected(), roleHandler.CreateRole)

	role.Get("/", jwt.Protected(), roleHandler.GetRoles)

	role.Get(constants.PathRoleIdParam, jwt.Protected(), roleHandler.GetRole)

	role.Put(constants.PathRoleIdParam, jwt.Protected(), roleHandler.UpdateRole)

	role.Delete(constants.PathRoleIdParam, jwt.Protected(), roleHandler.DeleteRole)

}
