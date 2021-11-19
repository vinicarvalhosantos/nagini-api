package roleUtil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"vcsxsantos/nagini-api/config"
	"vcsxsantos/nagini-api/database"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
	stringUtil "vcsxsantos/nagini-api/pkg/util/string"
)

const (
	USER    = 1
	ADMIN   = 2
	SUPPORT = 3
)

func CheckUserPermissions(c *fiber.Ctx, permission int) *dto.MyError {
	userC := c.Request().String()
	token := stringUtil.ExtractTokenFromString(userC)

	if token == "" {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusConflict
		myError.Map = fiber.Map{"status": "conflict", "message": "Authentication token was not found, please contact the administrator to solve this", "data": nil}
		return myError
	}

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetSecretKey("SECRET_KEY")), nil
	})

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": "internal_server_error", "message": "It was not possible to perform this action", "data": err.Error()}
		return myError
	}

	db := database.DB
	role := new(model.Role)
	roleId := claims["roleid"]

	db.Find(&role, "id = ?", roleId)

	if role.ID == 0 {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusConflict
		myError.Map = fiber.Map{"status": "conflict", "message": "Role ID was not found, please contact the administrator to solve this", "data": nil}
		return myError
	}

	if role.ID != ADMIN && permission != USER {
		if role.ID != permission {
			myError := new(dto.MyError)
			myError.Status = fiber.StatusUnauthorized
			myError.Map = fiber.Map{"status": "unauthorized", "message": "This user don't have the rights permissions to perform this action", "data": nil}
			return myError
		}
	}

	return nil
}
