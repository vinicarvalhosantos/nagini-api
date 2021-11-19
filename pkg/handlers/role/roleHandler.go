package roleHandler

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"vcsxsantos/nagini-api/pkg/model"
	roleRepoImpl "vcsxsantos/nagini-api/pkg/repository/role/impl"
	constantUtils "vcsxsantos/nagini-api/pkg/util/constant"
	roleUtil "vcsxsantos/nagini-api/pkg/util/role"
	stringUtil "vcsxsantos/nagini-api/pkg/util/string"
)

var repo = roleRepoImpl.NewRepository()

func CreateRole(c *fiber.Ctx) error {
	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	role := new(model.Role)

	err := c.BodyParser(&role)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	role, dtoError = repo.SaveRole(c, role)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constantUtils.StatusSuccess, "message": message(constantUtils.GenericCreateSuccessMessage), "data": role})
}

func GetRoles(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.SUPPORT)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	roles, dtoError := repo.FindRoles(c)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	return c.JSON(fiber.Map{"status": constantUtils.StatusSuccess, "message": message(constantUtils.GenericFoundSuccessMessage), "data": roles})
}

func GetRole(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.SUPPORT)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	idString := c.Params("roleId")
	id, err := strconv.Atoi(idString)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	role, dtoError := repo.FindRoleById(c, id)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	return c.JSON(fiber.Map{"status": constantUtils.StatusSuccess, "message": message(constantUtils.GenericFoundSuccessMessage), "data": role})
}

func UpdateRole(c *fiber.Ctx) error {
	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	var role *model.Role
	idString := c.Params("roleId")
	id, err := strconv.Atoi(idString)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	role, dtoErr := repo.FindRoleById(c, id)

	if dtoErr != nil {
		return c.Status(dtoError.Status).JSON(dtoErr.Map)
	}

	var updateRoleData *model.UpdateRole
	err = c.BodyParser(&updateRoleData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	role, dtoErr = repo.UpdateRole(c, id, role, updateRoleData)
	if dtoErr != nil {
		return c.Status(dtoError.Status).JSON(dtoErr.Map)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": constantUtils.StatusSuccess, "message": message(constantUtils.GenericUpdateSuccessMessage), "data": role})
}

func DeleteRole(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	idString := c.Params("roleId")
	id, err := strconv.Atoi(idString)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	dtoError = repo.DeleteRole(c, id)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})

}

func message(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "Role")
}
