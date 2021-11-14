package userHandler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"vcsxsantos/nagini-api/Internal/model"
	encrypt "vcsxsantos/nagini-api/Internal/util"
	"vcsxsantos/nagini-api/database"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User

	db.Find(&users)

	lenUsers := len(users)

	if lenUsers == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "not_found", "message": "Any user was founded", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("We found %d users registred", lenUsers), "data": users})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to create a new user", "data": err})
	}

	user.ID = uuid.New()
	user.Password, _ = encrypt.HashPassword(user.Password)

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to create a new user", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User created with successful", "data": user})
}

func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user model.User

	id := c.Params("userId")
	db.Find(&user, "id = ?", id)

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "not_found", "message": "Any user with this ID was founded", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User founded with successful", "data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
		Username     string
		UserFullName string
		Email        string
		Cpf          string
		Password     string
		Birthdate    string
	}

	db := database.DB
	var user model.User
	id := c.Params("userId")

	db.Find(&user, "id = ?", id)

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "not_found", "message": "Any user with this ID was founded", "data": nil})
	}

	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to update this user", "data": err})
	}

	user.Username = updateUserData.Username
	user.UserFullName = updateUserData.UserFullName
	user.Email = updateUserData.Email
	user.Cpf = updateUserData.Cpf
	user.Password, _ = encrypt.HashPassword(updateUserData.Password)
	user.Birthdate = updateUserData.Birthdate

	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User updated with successful", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	var user model.User

	id := c.Params("userId")

	db.Find(&user, "id = ?", id)

	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "not_found", "message": "Any user with this ID was founded", "data": nil})
	}

	err := db.Delete(&user, "id = ?", id).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "internal_server_error", "message": "It was not possible to delete this user", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User deleted with successful"})
}
