package userHandler

import (
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vinicius.csantos/nagini-api/database"
	"github.com/vinicius.csantos/nagini-api/internal/handler/email"
	"github.com/vinicius.csantos/nagini-api/internal/model"
	constants "github.com/vinicius.csantos/nagini-api/internal/util/constant"
	"github.com/vinicius.csantos/nagini-api/internal/util/cpfCNPJ"
	"github.com/vinicius.csantos/nagini-api/internal/util/encrypt"
	"github.com/vinicius.csantos/nagini-api/internal/util/jwt"
	stringUtil "github.com/vinicius.csantos/nagini-api/internal/util/string"
	"net/mail"
	"strings"
	"time"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []*model.User
	var readUsers []*model.ReadUser
	var userAddress []model.Address

	err := db.Find(&users).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	for _, user := range users {
		userAddress, err = getUserAddresses(c, user.ID)

		if err != nil {
			return err
		}

		if len(userAddress) != 0 {
			user.Address = userAddress
		}

		readUsers = append(readUsers, model.EntityToReadUser(user))
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageUser(constants.GenericFoundSuccessMessage), "data": readUsers})

}

func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user []*model.User
	var readUser []*model.ReadUser
	var userAddress []model.Address

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if len(user) == 0 || user[0].ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": readUser})
	}

	for _, user := range user {

		userAddress, err = getUserAddresses(c, user.ID)

		if err != nil {
			return err
		}

		if len(userAddress) != 0 {
			user.Address = userAddress
		}

		readUser = append(readUser, model.EntityToReadUser(user))
	}

	return c.JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageUser(constants.GenericFoundSuccessMessage), "data": readUser})
}

func Login(c *fiber.Ctx) error {
	db := database.DB
	var auth model.Authentication
	var user model.User

	err := c.BodyParser(&auth)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	condition := ""
	var login string

	if auth.Email != "" {
		condition += constants.EmailCondition
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	var token model.Token
	token.Login = login
	token.CpfCNPJ = stringUtil.RemoveSpecialCharacters(user.CpfCNPJ)
	token.UserID = user.ID
	token.Username = user.Username
	token.Expiration = time.Now().Add(time.Hour * 1)
	token.TokenString = validToken

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Login with successful", "data": token})
}

func RegisterUser(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User
	var newUser *model.User
	var readUser *model.ReadUser

	err := c.BodyParser(&newUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": constants.GenericInternalServerErrorMessage, "data": err.Error()})
	}

	isValid, invalidField := model.CheckIfUserEntityIsValid(newUser)

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": stringUtil.FormatGenericMessagesString(constants.GenericInvalidFieldMessage, invalidField), "data": nil})
	}

	formattedCpfCNPJ := stringUtil.RemoveSpecialCharacters(newUser.CpfCNPJ)
	formattedPhoneNumber := stringUtil.RemoveSpecialCharacters(newUser.PhoneNumber)

	_, err = mail.ParseAddress(newUser.Email)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.EmailInvalidMessage, "data": err.Error()})
	}

	err = db.Find(&user, "email = ? OR cpf_cnpj = ? OR username = ? OR phone_number = ?", newUser.Email,
		formattedCpfCNPJ, newUser.Username, formattedPhoneNumber).Error

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
		} else if user.PhoneNumber == formattedPhoneNumber {
			columnAlreadyExists = "phoneNumber"
		}

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": constants.StatusConflict, "message": stringUtil.FormatGenericMessagesString(constants.GenericAlreadyExistsMessage, columnAlreadyExists), "data": nil})
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
	newUser.PhoneNumber = formattedPhoneNumber

	err = db.Create(&newUser).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericCreateErrorMessage), "data": err.Error()})
	}
	readUser = model.EntityToReadUser(newUser)

	err = email.SendEmailToConfirmAccount(newUser.Email, newUser.Username, newUser.CpfCNPJ, newUser.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageUser(constants.GenericUserCreatedSuccessMessage), "data": readUser})
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User
	var readUser *model.ReadUser

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	var updateUserData *model.UpdateUser

	err = c.BodyParser(&updateUserData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericUpdateErrorMessage), "data": err.Error()})
	}

	readUser = model.EntityToReadUser(user)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": constants.StatusSuccess, "message": model.MessageUser(constants.GenericUpdateSuccessMessage), "data": readUser})

}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User

	id := c.Params("userId")

	err := db.Find(&user, constants.IdCondition, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	err = db.Model(&model.Address{}).Where(constants.UserIdCondition, user.ID).Delete(&model.Address{}).Error

	err = db.Delete(&user).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericDeleteErrorMessage), "data": err.Error()})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func ConfirmEmail(c *fiber.Ctx) error {
	userCache := model.UserCache
	db := database.DB

	token := c.Params("userToken")

	tokenFields, tokenDecrypted, err := getTokenFields(c, token)

	if err != nil {
		return err
	}

	userEmail := tokenFields[0]

	user, err := findUserByEmail(c, userEmail)

	if err != nil {
		return err
	}

	timeEncrypted := tokenFields[3]
	cacheKey := getCacheKey(user, timeEncrypted)

	err = validateUserToken(c, cacheKey, tokenDecrypted)

	if err != nil {
		return err
	}

	err = db.Model(&model.User{}).Where(constants.IdCondition, user.ID).Update("email_verified", true).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	userCache.Remove(cacheKey)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func RecoverPasswordRequest(c *fiber.Ctx) error {
	db := database.DB
	var user *model.User
	var recoverPasswordRequest *model.RecoverPassword

	err := c.BodyParser(&recoverPasswordRequest)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	formattedCpfCnpj := stringUtil.RemoveSpecialCharacters(recoverPasswordRequest.CpfCNPJ)

	err = db.Find(&user, "cpf_cnpj = ?", formattedCpfCnpj).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID != uuid.Nil {

		err = email.SendEmailToRecoverPassword(user.Email, user.Username, formattedCpfCnpj, user.ID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
		}

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "An email with instructions to recover your password was sent to email registered if this Cpf or Cnpj exists in our database! ", "data": nil})
}

func ChangePassword(c *fiber.Ctx) error {
	userCache := model.UserCache
	db := database.DB
	var changePasswordModel *model.ChangePassword

	token := c.Params("userToken")

	err := c.BodyParser(&changePasswordModel)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	tokenFields, tokenDecrypted, err := getTokenFields(c, token)

	if err != nil {
		return err
	}

	userEmail := tokenFields[0]

	user, err := findUserByEmail(c, userEmail)

	if err != nil {
		return err
	}

	timeEncrypted := tokenFields[3]
	cacheKey := getCacheKey(user, timeEncrypted)

	err = validateUserToken(c, cacheKey, tokenDecrypted)

	if err != nil {
		return err
	}

	if changePasswordModel.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": stringUtil.FormatGenericMessagesString(constants.GenericInvalidFieldMessage, "Password"), "data": nil})
	}

	newPassword, err := encrypt.HashPassword(changePasswordModel.NewPassword)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	err = db.Model(&model.User{}).Where(constants.IdCondition, user.ID).Update("password", newPassword).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	userCache.Remove(cacheKey)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": constants.StatusSuccess, "message": "Your password was reset with successful!", "data": nil})
}

func getUserAddresses(c *fiber.Ctx, userId uuid.UUID) ([]model.Address, error) {
	db := database.DB
	var userAddress []model.Address

	err := db.Find(&userAddress, constants.UserIdCondition, userId).Error

	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	return userAddress, nil
}

func validateUserToken(c *fiber.Ctx, cacheKey, tokenDecrypted string) error {
	userCache := model.UserCache
	userToken, err := userCache.Get(cacheKey)

	if userToken == tokenDecrypted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": constants.StatusBadRequest, "message": constants.GenericTokenDoesNotMatch, "data": nil})
	}

	if err != nil {
		if err == ttlcache.ErrNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": constants.StatusForbidden, "message": model.MessageUser(constants.GenericCacheForbiddenMessage), "data": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}
	return nil
}

func getTokenFields(c *fiber.Ctx, token string) ([]string, string, error) {
	tokenDecrypted, err := encrypt.UrlDecrypt(token)

	if err != nil {
		return nil, "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	tokenFields := strings.Split(tokenDecrypted, "$")

	return tokenFields, tokenDecrypted, nil
}

func getCacheKey(user *model.User, timeEncrypted string) string {
	return fmt.Sprintf("%s-%s", user.ID.String(), timeEncrypted)
}

func findUserByEmail(c *fiber.Ctx, userEmail string) (*model.User, error) {
	var user *model.User
	db := database.DB

	err := db.Find(&user, constants.EmailCondition, userEmail).Error

	if err != nil {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": constants.StatusInternalServerError, "message": model.MessageUser(constants.GenericInternalServerErrorMessage), "data": err.Error()})
	}

	if user.ID == uuid.Nil {
		return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": constants.StatusNotFound, "message": model.MessageUser(constants.GenericNotFoundMessage), "data": nil})
	}

	return user, nil
}
