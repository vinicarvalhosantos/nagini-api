package model

import (
	"github.com/google/uuid"
	"gitlab.com/vinicius.csantos/nagini-api/internal/util/encrypt"
	stringUtil "gitlab.com/vinicius.csantos/nagini-api/internal/util/string"
	"time"
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Support UserRole = "support"
	UserR   UserRole = "user"
)

type User struct {
	ID           uuid.UUID
	Username     string `gorm:"index;unique;not null;"`
	UserFullName string `gorm:"not null;"`
	Email        string `gorm:"index;unique;not null;"`
	CpfCNPJ      string `gorm:"index;unique;not null;"`
	Password     string `gorm:"not null;"`
	Birthdate    string
	PhoneNumber  string    `gorm:"unique;not null;"`
	Role         UserRole  `gorm:"not null;"`
	Address      []Address `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ReadUser struct {
	ID           uuid.UUID
	Username     string
	UserFullName string
	Email        string
	CpfCNPJ      string
	Birthdate    string
	PhoneNumber  string
	Role         UserRole
	Address      []Address
}

type UpdateUser struct {
	Username     string
	UserFullName string
	Email        string
	CpfCNPJ      string
	Password     string
	Birthdate    string
	PhoneNumber  string
	Role         UserRole
}

type Authentication struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Login       string `json:"login"`
	CpfCNPJ     string `json:"cpfCnpj"`
	TokenString string `json:"token"`
}

func CheckIfUserEntityIsValid(user *User) (bool, string) {

	if user.UserFullName == "" {
		return false, "UserFullName"
	}
	if user.Username == "" {
		return false, "Username"
	}
	if user.Email == "" {
		return false, "Email"
	}
	if user.CpfCNPJ == "" {
		return false, "CpfCnpj"
	}
	if user.PhoneNumber == "" {
		return false, "PhoneNumber"
	}
	if user.Role == "" {
		return false, "Role"
	}

	return true, ""
}

func EntityToReadUser(user *User) *ReadUser {
	readUser := new(ReadUser)
	readUser.ID = user.ID
	readUser.Username = user.Username
	readUser.UserFullName = user.UserFullName
	readUser.Email = user.Email
	readUser.CpfCNPJ = user.CpfCNPJ
	readUser.Birthdate = user.Birthdate
	readUser.PhoneNumber = user.PhoneNumber
	readUser.Role = user.Role
	readUser.Address = user.Address

	return readUser
}

func PrepareUserToUpdate(user *User, updateUserData *UpdateUser) *User {

	if updateUserData.Username != "" {
		user.Username = updateUserData.Username
	}
	if updateUserData.UserFullName != "" {
		user.UserFullName = updateUserData.UserFullName
	}
	if updateUserData.Email != "" {
		user.Email = updateUserData.Email
	}
	if updateUserData.CpfCNPJ != "" {
		user.CpfCNPJ = stringUtil.RemoveSpecialCharacters(updateUserData.CpfCNPJ)
	}
	if updateUserData.Password != "" {
		user.Password, _ = encrypt.HashPassword(updateUserData.Password)
	}
	if updateUserData.Birthdate != "" {
		user.Birthdate = updateUserData.Birthdate
	}
	if updateUserData.PhoneNumber != "" {
		user.PhoneNumber = stringUtil.RemoveSpecialCharacters(updateUserData.PhoneNumber)
	}

	return user
}

func MessageUser(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "User")
}
