package model

import (
	"github.com/google/uuid"
	"time"
)

type Address struct {
	ID           uint
	Cep          string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
	UserID       uuid.UUID
	MainAddress  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
