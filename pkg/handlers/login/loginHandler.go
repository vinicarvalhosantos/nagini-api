package loginHandler

import (
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/database"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
	"vcsxsantos/nagini-api/pkg/util/encrypt"
	"vcsxsantos/nagini-api/pkg/util/jwt"
)

func Login(c *fiber.Ctx) error {
	db := database.DB
	var auth dto.Authentication
	var user model.User
	userRole := new(model.Role)

	err := c.BodyParser(&auth)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to signing this user", "data": err.Error()})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "not_found", "message": "Username or password is incorrect", "data": nil})
	}

	passwordMatches := encrypt.CheckPasswordHash(user.Password, auth.Password)

	if !passwordMatches {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "not_found", "message": "Username or password is incorrect", "data": nil})
	}

	db.Find(&userRole, "id = ?", user.RoleID)

	validToken, err := jwt.GenerateToken(user.UserFullName, user.Username, user.Email, userRole.Name, userRole.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to signing this user", "data": err.Error()})
	}

	var token dto.Token
	token.Login = login
	token.RoleName = userRole.Name
	token.RoleID = userRole.ID
	token.TokenString = validToken

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Credentials was correct! Login with successful", "data": token})

}
