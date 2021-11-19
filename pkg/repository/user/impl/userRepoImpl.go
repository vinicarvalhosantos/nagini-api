package userRepoImpl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"vcsxsantos/nagini-api/database"
	"vcsxsantos/nagini-api/pkg/dto"
	"vcsxsantos/nagini-api/pkg/model"
	"vcsxsantos/nagini-api/pkg/repository/role/impl"
	userRepo "vcsxsantos/nagini-api/pkg/repository/user"
	constantUtils "vcsxsantos/nagini-api/pkg/util/constant"
	cpfCNPJUtil "vcsxsantos/nagini-api/pkg/util/cpfCNPJ"
	"vcsxsantos/nagini-api/pkg/util/encrypt"
	stringUtil "vcsxsantos/nagini-api/pkg/util/string"
)

type userRepositoryImpl struct {
}

var roleRepo = roleRepoImpl.NewRepository()

func (u userRepositoryImpl) FindUsers(ctx *fiber.Ctx) ([]model.User, *dto.MyError) {
	db := database.DB
	var err error
	var users []model.User

	err = db.Find(&users).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return nil, myError
	}

	if len(users) == 0 {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusNotFound
		myError.Map = fiber.Map{"status": constantUtils.STATUS_NOT_FOUND, "message": message(constantUtils.GENERIC_NOT_FOUND_MESSAGE), "data": nil}
		return nil, myError
	}

	for index, user := range users {
		userRole, dtoErr := roleRepo.FindRoleById(ctx, user.RoleID)

		if dtoErr != nil {
			return nil, dtoErr
		}

		users[index].Role = *userRole
	}

	return users, nil
}

func (u userRepositoryImpl) FindUserById(ctx *fiber.Ctx, id string) (*model.User, *dto.MyError) {
	db := database.DB
	var user *model.User

	err := db.Find(&user, "id = ?", id).Error

	dtoErr := checkUserErr(user, err)
	if dtoErr != nil {
		return nil, dtoErr
	}

	user, dtoErr = setUserRole(ctx, user)

	return user, nil
}

func (u userRepositoryImpl) FindUserByCpfCNPJ(ctx *fiber.Ctx, cpfCNPJ string) (*model.User, *dto.MyError) {

	if cpfCNPJUtil.ValidateCpfCNPJ(stringUtil.RemoveSpecialCharacters(cpfCNPJ)) {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusBadRequest
		myError.Map = fiber.Map{"status": constantUtils.STATUS_BAD_REQUEST, "message": constantUtils.CpfCnpjInvalidMessage, "data": nil}
		return nil, myError
	}

	db := database.DB
	var user *model.User

	err := db.Find(&user, "cpfCNPJ = ?", cpfCNPJ).Error

	dtoErr := checkUserErr(user, err)
	if dtoErr != nil {
		return nil, dtoErr
	}

	user, dtoErr = setUserRole(ctx, user)

	return user, nil
}

func (u userRepositoryImpl) FindUserByUsername(ctx *fiber.Ctx, username string) (*model.User, *dto.MyError) {
	db := database.DB
	var user *model.User

	err := db.Find(&user, "username = ?", username).Error

	dtoErr := checkUserErr(user, err)
	if dtoErr != nil {
		return nil, dtoErr
	}

	user, dtoErr = setUserRole(ctx, user)

	return user, nil
}

func (u userRepositoryImpl) FindUserByEmail(ctx *fiber.Ctx, email string) (*model.User, *dto.MyError) {
	db := database.DB
	var user *model.User

	err := db.Find(&user, "email = ?", email).Error

	dtoErr := checkUserErr(user, err)
	if dtoErr != nil {
		return nil, dtoErr
	}

	user, dtoErr = setUserRole(ctx, user)

	return user, nil
}

func (u userRepositoryImpl) FindUserByEmailORCpfCNPJORUsername(ctx *fiber.Ctx, email, cpfCNPJ, username string) (*model.User, *dto.MyError) {
	db := database.DB
	user := new(model.User)

	err := db.Find(&user, "email = ? OR cpf_cnpj = ? OR username = ?", email, stringUtil.RemoveSpecialCharacters(cpfCNPJ), username).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return nil, myError
	}

	if user.ID != uuid.Nil {
		columnAlreadyExists := ""
		if user.Email == email {
			columnAlreadyExists = "email"
		} else if user.CpfCNPJ == stringUtil.RemoveSpecialCharacters(cpfCNPJ) {
			columnAlreadyExists = "cpfCNPJ"
		} else if user.Username == username {
			columnAlreadyExists = "username"
		}

		myError := new(dto.MyError)
		myError.Status = fiber.StatusConflict
		myError.Map = fiber.Map{"status": constantUtils.STATUS_CONFLICT, "message": "This " + columnAlreadyExists + " already exists on our database", "data": nil}
		return nil, myError
	}

	return user, nil
}

func (u userRepositoryImpl) SaveUser(ctx *fiber.Ctx, user *model.User) (*model.User, *dto.MyError) {
	db := database.DB

	_, dtoErr := u.FindUserByEmailORCpfCNPJORUsername(ctx, user.Email, user.CpfCNPJ, user.Username)

	if dtoErr != nil {
		return nil, dtoErr
	}

	userRole, dtoErr := roleRepo.FindRoleById(ctx, user.Role.ID)

	if dtoErr != nil {
		return nil, dtoErr
	}

	if !cpfCNPJUtil.ValidateCpfCNPJ(stringUtil.RemoveSpecialCharacters(user.CpfCNPJ)) {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusBadRequest
		myError.Map = fiber.Map{"status": constantUtils.STATUS_BAD_REQUEST, "message": constantUtils.CpfCnpjInvalidMessage, "data": nil}
		return nil, myError
	}

	user.ID = uuid.New()
	user.Password, _ = encrypt.HashPassword(user.Password)
	user.CpfCNPJ = stringUtil.RemoveSpecialCharacters(user.CpfCNPJ)
	user.RoleID = userRole.ID
	user.Role = *userRole

	err := db.Create(&user).Error
	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": message(constantUtils.GENERIC_CREATE_ERROR_MESSAGE), "data": err.Error()}
		return nil, myError
	}

	return user, nil

}

func (u userRepositoryImpl) UpdateUser(ctx *fiber.Ctx, id string, user *model.User, updateUser *model.UpdateUser) (*model.User, *dto.MyError) {
	db := database.DB
	var roleUser *model.Role

	if updateUser.Role.ID != 0 {
		role, dtoErr := roleRepo.FindRoleById(ctx, updateUser.Role.ID)

		if dtoErr != nil {
			return nil, dtoErr
		}

		roleUser = role
	} else {

		role, dtoErr := roleRepo.FindRoleById(ctx, user.RoleID)
		if dtoErr != nil {
			return nil, dtoErr
		}

		roleUser = role
	}

	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.UserFullName != "" {
		user.UserFullName = updateUser.UserFullName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.CpfCNPJ != "" {
		user.CpfCNPJ = stringUtil.RemoveSpecialCharacters(updateUser.CpfCNPJ)
	}
	if updateUser.Password != "" {
		user.Password, _ = encrypt.HashPassword(updateUser.Password)
	}
	if updateUser.Birthdate != "" {
		user.Birthdate = updateUser.Birthdate
	}

	user.RoleID = roleUser.ID
	user.Role = *roleUser

	if !cpfCNPJUtil.ValidateCpfCNPJ(user.CpfCNPJ) {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusBadRequest
		myError.Map = fiber.Map{"status": constantUtils.STATUS_BAD_REQUEST, "message": constantUtils.CpfCnpjInvalidMessage, "data": nil}
		return nil, myError
	}

	err := db.Save(&user).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return nil, myError
	}

	return user, nil
}

func (u userRepositoryImpl) DeleteUser(ctx *fiber.Ctx, id string) *dto.MyError {
	db := database.DB

	user, dtoErr := u.FindUserById(ctx, id)

	if dtoErr != nil {
		return dtoErr
	}

	err := db.Delete(user).Error

	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": message(constantUtils.GENERIC_DELETE_ERROR_MESSAGE), "data": err.Error()}
		return myError
	}

	return nil
}

func NewRepository() userRepo.UserRepository {
	return *&userRepositoryImpl{}
}

func message(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "User")
}

func checkUserErr(user *model.User, err error) *dto.MyError {
	if err != nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusInternalServerError
		myError.Map = fiber.Map{"status": constantUtils.StatusInternalServerError, "message": constantUtils.GenericInternalServerErrorMessage, "data": err.Error()}
		return myError
	}

	if user.ID == uuid.Nil {
		myError := new(dto.MyError)
		myError.Status = fiber.StatusNotFound
		myError.Map = fiber.Map{"status": constantUtils.STATUS_NOT_FOUND, "message": "Any user with this ID was found", "data": nil}
		return myError
	}

	return nil
}

func setUserRole(ctx *fiber.Ctx, user *model.User) (*model.User, *dto.MyError) {
	var role = new(model.Role)
	role, dtoErr := roleRepo.FindRoleById(ctx, user.RoleID)

	if dtoErr != nil {
		return nil, dtoErr
	}

	user.Role = *role

	return user, nil
}
