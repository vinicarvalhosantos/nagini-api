package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"time"
	"vcsxsantos/nagini-api/config"
)

func GenerateToken(userFullName, username, email, roleName string, roleId int) (string, error) {
	secretKey := []byte(config.GetSecretKey("SECRET_KEY"))
	claims := jwt.MapClaims{
		"authorized":   true,
		"userfullname": userFullName,
		"username":     username,
		"email":        email,
		"role":         roleName,
		"roleid":       roleId,
		"now":          time.Now(),
		"exp":          time.Now().Add(time.Hour * 1).Unix(), //Valid for 1 hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		err = fmt.Errorf("Is was not possible to create a valid token for this user\n Error: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func Protected() func(ctx *fiber.Ctx) error {
	jwtConfig := jwtMiddleware.Config{
		SigningKey:   []byte(config.GetSecretKey("SECRET_KEY")),
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
	}
	return jwtMiddleware.New(jwtConfig)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "bad_request", "message": "Unauthorized", "data": err.Error()})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "unauthorized", "message": "Unauthorized", "data": nil})

}
