package roleRepo

import (
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
)

type RoleRepository interface {
	// FindRoleById Find role by id
	FindRoleById(ctx *fiber.Ctx, id int) (*model.Role, *dto.MyError)

	// FindRoleByName Find role by name
	FindRoleByName(ctx *fiber.Ctx, name string) (*model.Role, *dto.MyError)

	// FindRoles Find all roles
	FindRoles(ctx *fiber.Ctx) ([]model.Role, *dto.MyError)

	// SaveRole Save a role
	SaveRole(ctx *fiber.Ctx, role *model.Role) (*model.Role, *dto.MyError)

	// UpdateRole Update a role
	UpdateRole(ctx *fiber.Ctx, id int, role *model.Role, updateRole *model.UpdateRole) (*model.Role, *dto.MyError)

	// DeleteRole Delete a role
	DeleteRole(ctx *fiber.Ctx, id int) *dto.MyError
}
