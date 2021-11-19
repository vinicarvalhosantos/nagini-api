package model

import (
	"time"
)

type Role struct {
	ID          int
	Name        string `gorm:"uniqueIndex"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateRole struct {
	Name        string
	Description string
}
