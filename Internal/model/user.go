package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid"`
	Username     string
	UserFullName string
	Email        string
	Cpf          string
	Password     string `json:"-"`
	Birthdate    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
