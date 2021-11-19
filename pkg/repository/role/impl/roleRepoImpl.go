package roleRepoImpl

import (
	"github.com/gofiber/fiber/v2"
	"vcsxsantos/nagini-api/database"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
	roleRepo "vcsxsantos/nagini-api/pkg/repository/role"
	constantUtils "vcsxsantos/nagini-api/pkg/util/constant"
	stringUtil "vcsxsantos/nagini-api/pkg/util/string"
)

type roleRepositoryImpl struct {
}

func (r roleRepositoryImpl) FindRoleById(ctx *fiber.Ctx, id int) (*model.Role, *dto.MyError) {
	db := database.DB
	var role *model.Role

	err := db.Find(&role, "id = ?", id).Error

	dtoErr := checkRoleErr(role, err)

	if dtoErr != nil {
		return nil, dtoErr
	}

	return role, nil
}

func (r roleRepositoryImpl) FindRoles(ctx *fiber.Ctx) ([]model.Role, *dto.MyError) {
	db := database.DB
	var roles []model.Role

	err := db.Find(&roles).Error

	if err == nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return nil, myError
	}

	if len(roles) == 0 {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusNotFound
		myError.Map = fiber.Map{"status": constantUtils.STATUS_NOT_FOUND, "message": message(constantUtils.GENERIC_NOT_FOUND_MESSAGE), "data": nil}
		return nil, myError
	}

	return roles, nil

}

func (r roleRepositoryImpl) FindRoleByName(ctx *fiber.Ctx, name string) (*model.Role, *dto.MyError) {
	db := database.DB
	var role *model.Role

	err := db.Find(&role, "name ilike \"?\"", name).Error

	dtoErr := checkRoleErr(role, err)

	if dtoErr != nil {
		return nil, dtoErr
	}

	return role, nil
}

func (r roleRepositoryImpl) SaveRole(ctx *fiber.Ctx, role *model.Role) (*model.Role, *dto.MyError) {

	rolesExists, dtoErr := r.FindRoleByName(ctx, role.Name)

	if dtoErr != nil && dtoErr.Status != fiber.StatusNotFound {
		return nil, dtoErr
	}

	if rolesExists != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusConflict
		myError.Map = fiber.Map{"status": constantUtils.STATUS_CONFLICT, "message": "This role already exists on our database", "data": nil}
		return nil, myError
	}

	db := database.DB
	err := db.Create(role).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": message(constantUtils.GENERIC_CREATE_ERROR_MESSAGE), "data": err.Error()}
		return nil, myError
	}

	return role, nil
}

func (r roleRepositoryImpl) UpdateRole(ctx *fiber.Ctx, id int, role *model.Role, updateRole *model.UpdateRole) (*model.Role, *dto.MyError) {
	db := database.DB

	if updateRole.Description != "" {
		role.Description = updateRole.Description
	}
	if updateRole.Name != "" {
		role.Name = updateRole.Name
	}

	err := db.Save(role).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": message(constantUtils.GENERIC_UPDATE_ERROR_MESSAGE), "data": err.Error()}
		return nil, myError
	}

	return role, nil

}

func (r roleRepositoryImpl) DeleteRole(ctx *fiber.Ctx, id int) *dto.MyError {
	db := database.DB

	role, dtoErr := r.FindRoleById(ctx, id)

	if dtoErr != nil {
		return dtoErr
	}

	err := db.Delete(role).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": message(constantUtils.GENERIC_DELETE_ERROR_MESSAGE), "data": err.Error()}
		return myError
	}

	return nil

}

func NewRepository() roleRepo.RoleRepository {
	return *&roleRepositoryImpl{}
}

func message(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "Role")
}

func checkRoleErr(role *model.Role, err error) *dto.MyError {
	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return myError
	}

	if role.ID == 0 {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusNotFound
		myError.Map = fiber.Map{"status": constantUtils.STATUS_NOT_FOUND, "message": message(constantUtils.GENERIC_NOT_FOUND_MESSAGE), "data": nil}
		return myError
	}

	return nil
}
