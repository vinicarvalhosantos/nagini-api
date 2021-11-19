package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Username     string    `gorm:"uniqueIndex"`
	UserFullName string
	Email        string `gorm:"unique"`
	CpfCNPJ      string `gorm:"uniqueIndex"`
	Password     string `json:"-"`
	Birthdate    string
	RoleID       int `gorm:"default:1" json:"-"`
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpdateUser struct {
	Username     string
	UserFullName string
	Email    string
	CpfCNPJ  string
	Password string
	Birthdate    string
	Role         Role
}
