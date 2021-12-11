package userHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/vinicius.csantos/nagini-api/database"
	"gitlab.com/vinicius.csantos/nagini-api/internal/model"
	constants "gitlab.com/vinicius.csantos/nagini-api/internal/util/constant"
	"gitlab.com/vinicius.csantos/nagini-api/internal/util/cpfCNPJ"
	"gitlab.com/vinicius.csantos/nagini-api/internal/util/encrypt"
	"gitlab.com/vinicius.csantos/nagini-api/internal/util/jwt"
	stringUtil "gitlab.com/vinicius.csantos/nagini-api/internal/util/string"
	"net/mail"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []*model.User
	var readUsers []*model.ReadUser

	err := db.Find(&users).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": message(constants.GenericNotFoundMessage), "data": users})
	}

	for _, user := range users {
		readUsers = append(readUsers, model.EntityToReadUser(user))
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericFoundSuccessMessage), "data": readUsers})

}

func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user []*model.User
	var readUser []*model.ReadUser

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user[0].ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": message(constants.GenericNotFoundMessage), "data": readUser})
	}

	for _, user := range user {
		readUser = append(readUser, model.EntityToReadUser(user))
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericFoundSuccessMessage), "data": readUser})
}

func Login(c *fiber.Ctx) error {
	db := database.DB
	var auth model.Authentication
	var user model.User

	err := c.BodyParser(&auth)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	condition := ""
	var login string

	if auth.Email != "" {
		condition += "email = ?"
		login = auth.Email
	}

	if auth.Username != "" {
		condition += "username = ?"
		login = auth.Username
	}

	db.Find(&user, condition, login)

	if user.Username == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": "Username or password is incorrect", "data": nil})
	}

	passwordMatches := encrypt.CheckPasswordHash(user.Password, auth.Password)

	if !passwordMatches {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": "Username or password is incorrect", "data": nil})
	}

	validToken, err := jwt.GenerateToken(user.UserFullName, user.Username, user.Email, string(user.Role))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	var token model.Token
	token.Login = login
	token.CpfCNPJ = stringUtil.RemoveSpecialCharacters(user.CpfCNPJ)
	token.TokenString = validToken

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Login with successful", "data": token})
}

func RegisterUser(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User
	var newUser *model.User
	var readUser *model.ReadUser

	err := c.BodyParser(&newUser)

	isValid, invalidField := model.CheckIfEntityIsValid(newUser)

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": stringUtil.FormatGenericMessagesString(constants.GenericInvalidFieldMessage, invalidField), "data": nil})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	formattedCpfCNPJ := stringUtil.RemoveSpecialCharacters(newUser.CpfCNPJ)

	_, err = mail.ParseAddress(newUser.Email)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.EmailInvalidMessage, "data": err.Error()})
	}

	err = db.Find(&user, "email = ? OR cpf_cnpj = ? OR username = ? OR phone_number = ?", newUser.Email,
		formattedCpfCNPJ, newUser.Username, newUser.PhoneNumber).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	if user.ID != uuid.Nil {
		columnAlreadyExists := ""
		if user.Email == newUser.Email {
			columnAlreadyExists = "email"
		} else if user.CpfCNPJ == formattedCpfCNPJ {
			columnAlreadyExists = "cpfCNPJ"
		} else if user.Username == newUser.Username {
			columnAlreadyExists = "username"
		} else if user.PhoneNumber == newUser.PhoneNumber {
			columnAlreadyExists = "phoneNumber"
		}

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": constants.StatusConflict, "message": "This " + columnAlreadyExists + " already exists on our database", "data": nil})
	}

	if !cpfCNPJ.ValidateCpfCNPJ(formattedCpfCNPJ) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.CpfCnpjInvalidMessage, "data": nil})
	}

	if newUser.Role != model.Admin && newUser.Role != model.UserR && newUser.Role != model.Support {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.RoleInvalidMessage, "data": nil})
	}

	newUser.ID = uuid.New()
	newUser.Password, _ = encrypt.HashPassword(newUser.Password)
	newUser.CpfCNPJ = formattedCpfCNPJ

	err = db.Create(&newUser).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericCreateErrorMessage), "data": err.Error()})
	}
	readUser = model.EntityToReadUser(newUser)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericCreateSuccessMessage), "data": readUser})
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User
	var readUser *model.ReadUser

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": message(constants.GenericNotFoundMessage), "data": nil})
	}

	var updateUserData *model.UpdateUser

	err = c.BodyParser(&updateUserData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if updateUserData.Role != "" {
		if updateUserData.Role != model.Admin && updateUserData.Role != model.UserR && updateUserData.Role != model.Support {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.RoleInvalidMessage, "data": nil})
		}
		user.Role = updateUserData.Role
	}

	user = model.PrepareUserToUpdate(user, updateUserData)

	if !cpfCNPJ.ValidateCpfCNPJ(user.CpfCNPJ) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.CpfCnpjInvalidMessage, "data": nil})
	}

	err = db.Save(&user).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	readUser = model.EntityToReadUser(user)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": constants.StatusSuccess, "message": message(constants.GenericUpdateSuccessMessage), "data": readUser})

}

func DeleteUser(c *fiber.Ctx) error{
	db := database.DB
	var user *model.User

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": message(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Delete(&user).Error

	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": message(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func message(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "User")
}
