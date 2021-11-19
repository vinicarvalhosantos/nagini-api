package userRepo

import (
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
)

type UserRepository interface {
	// FindUsers Find all users ins our database
	FindUsers(ctx *fiber.Ctx) ([]model.User, *dto.MyError)

	// FindUserById Find user by uuid
	FindUserById(ctx *fiber.Ctx, id string) (*model.User, *dto.MyError)

	// FindUserByCpfCNPJ Find user by cpf or CNPJ
	FindUserByCpfCNPJ(ctx *fiber.Ctx, cpfCNPJ string) (*model.User, *dto.MyError)

	// FindUserByUsername Find user by username
	FindUserByUsername(ctx *fiber.Ctx, username string) (*model.User, *dto.MyError)

	// FindUserByEmail Find user by email
	FindUserByEmail(ctx *fiber.Ctx, email string) (*model.User, *dto.MyError)

	// FindUserByEmailORCpfCNPJORUsername Find user by email or cpf and cnpj or username
	FindUserByEmailORCpfCNPJORUsername(ctx *fiber.Ctx, email, cpfCNPJ, username string) (*model.User, *dto.MyError)

	// SaveUser Save a user
	SaveUser(ctx *fiber.Ctx, user *model.User) (*model.User, *dto.MyError)

	// UpdateUser Update a user
	UpdateUser(ctx *fiber.Ctx, id string, user *model.User, updateUser *model.UpdateUser) (*model.User, *dto.MyError)

	// DeleteUser Delete a user
	DeleteUser(ctx *fiber.Ctx, id string) *dto.MyError
}
