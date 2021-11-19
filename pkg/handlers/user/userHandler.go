package userHandler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/pkg/model"
	"vcsxsantos/nagini-api/pkg/repository/user/impl"
	constants "vcsxsantos/nagini-api/pkg/util/constant"
	roleUtil "vcsxsantos/nagini-api/pkg/util/role"
	stringUtil "vcsxsantos/nagini-api/pkg/util/string"
)

var repo = userRepoImpl.NewRepository()

func GetUsers(c *fiber.Ctx) error {

	dtoErr := roleUtil.CheckUserPermissions(c, roleUtil.SUPPORT)

	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	users, repoError := repo.FindUsers(c)

	if repoError != nil {
		return c.Status(repoError.Status).JSON(repoError.Map)
	}

	lenUsers := len(users)

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": fmt.Sprintf("We found %d users registred", lenUsers), "data": users})
}

func RegisterUser(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}
	user := new(model.User)

	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	user, dtoErr := repo.SaveUser(c, user)
	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericCreateSuccessMessage), "data": user})
}

func GetUser(c *fiber.Ctx) error {

	dtoErr := roleUtil.CheckUserPermissions(c, roleUtil.USER)

	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	id := c.Params("userId")
	user, dtoErr := repo.FindUserById(c, id)

	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericFoundSuccessMessage), "data": user})
}

func UpdateUser(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	id := c.Params("userId")

	user, dtoError := repo.FindUserById(c, id)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	var updateUserData = new(model.UpdateUser)
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	user, dtoErr := repo.UpdateUser(c, id, user, updateUserData)

	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericUpdateSuccessMessage), "data": user})
}

func DeleteUser(c *fiber.Ctx) error {

	dtoError := roleUtil.CheckUserPermissions(c, roleUtil.ADMIN)

	if dtoError != nil {
		return c.Status(dtoError.Status).JSON(dtoError.Map)
	}

	id := c.Params("userId")

	dtoErr := repo.DeleteUser(c, id)

	if dtoErr != nil {
		return c.Status(dtoErr.Status).JSON(dtoErr.Map)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func message(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "User")
}
