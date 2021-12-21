package address

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vinicius.csantos/nagini-api/database"
	"github.com/vinicius.csantos/nagini-api/internal/model"
	constants "github.com/vinicius.csantos/nagini-api/internal/util/constant"
	stringUtil "github.com/vinicius.csantos/nagini-api/internal/util/string"
)

func GetAddresses(c *fiber.Ctx) error {
	db := database.DB
	var addresses []*model.Address

	err := db.Find(&addresses).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if len(addresses) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericFoundSuccessMessage), "data": addresses})

}

func GetAddressById(c *fiber.Ctx) error {
	db := database.DB
	var address *model.Address

	id := c.Params("addressId")

	err := db.Find(&address, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if address.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericFoundSuccessMessage), "data": address})

}

func GetUserAddressesById(c *fiber.Ctx) error {
	db := database.DB
	var address []*model.Address

	userId := c.Params("userId")

	err := db.Find(&address, constants.UserIdCondition, userId).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if len(address) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericFoundSuccessMessage), "data": address})

}

func RegisterAddress(c *fiber.Ctx) error {
	db := database.DB
	var newAddress *model.Address
	var userAddresses []*model.Address
	var user *model.User

	err := c.BodyParser(&newAddress)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	isValid, invalidField := model.CheckIfAddressEntityIsValid(newAddress)

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": stringUtil.FormatGenericMessagesString(constants.GenericInvalidFieldMessage, invalidField), "data": nil})
	}

	err = db.Find(&user, constants.IdCondition, newAddress.UserID).Error

	if err != nil || user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Find(&userAddresses, constants.UserIdCondition, user.ID).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if len(userAddresses) != 0 {
		newAddress.MainAddress = false //Endereço padrão sempre será falso se exister um ou mais endereços cadastrados para esse usuário
	} else {
		newAddress.MainAddress = true //Endereço padrão sempre será true se nao existir nenhum endereço cadastrado no usuário
	}
	err = db.Create(&newAddress).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageAddress(constants.GenericCreateErrorMessage), "data": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericCreateSuccessMessage), "data": newAddress})

}

func UpdateUserMainAddress(c *fiber.Ctx) error {
	db := database.DB
	var address *model.Address
	var user *model.User

	addressId := c.Params("addressId")
	userId := c.Params("userId")

	err := db.Find(&address, constants.IdCondition, addressId).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if address.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Find(&user, constants.IdCondition, userId).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Model(&model.Address{}).Where(constants.UserIdCondition, user.ID).Update("main_address", false).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	address.MainAddress = true

	err = db.Save(address).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericFoundSuccessMessage), "data": address})
}

func UpdateAddress(c *fiber.Ctx) error {
	db := database.DB
	var updateAddress *model.UpdateAddress
	var address *model.Address

	err := c.BodyParser(&updateAddress)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	id := c.Params("addressId")

	err = db.Find(&address, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if address.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	address = model.PrepareAddressToUpdate(address, updateAddress)

	err = db.Save(address).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericUpdateErrorMessage, "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageAddress(constants.GenericFoundSuccessMessage), "data": address})
}

func DeleteAddress(c *fiber.Ctx) error {
	db := database.DB
	var address model.Address

	id := c.Params("addressId")

	err := db.Find(&address, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if address.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageAddress(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Delete(&address).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericDeleteErrorMessage, "data": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
