package model

import (
	"github.com/google/uuid"
	stringUtil "github.com/vinicius.csantos/nagini-api/internal/util/string"
	"time"
)

type Address struct {
	ID           uint
	Name         string
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

type UpdateAddress struct {
	Name         string
	Cep          string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Country      string
}

func CheckIfAddressEntityIsValid(address *Address) (bool, string) {

	if address.Cep == "" {
		return false, "Cep"
	}
	if address.AddressLine1 == "" {
		return false, "AddressLine1"
	}
	if address.City == "" {
		return false, "City"
	}
	if address.State == "" {
		return false, "State"
	}
	if address.Country == "" {
		return false, "Country"
	}
	if address.UserID == uuid.Nil {
		return false, "UserID"
	}

	return true, ""
}

func PrepareAddressToUpdate(address *Address, updateAddress *UpdateAddress) *Address {

	if updateAddress.Name != "" {
		address.Name = updateAddress.Name
	}

	if updateAddress.Cep != "" {
		address.Cep = updateAddress.Cep
	}

	if updateAddress.AddressLine1 != "" {
		address.AddressLine1 = updateAddress.AddressLine1
	}

	if updateAddress.AddressLine2 != "" {
		address.AddressLine2 = updateAddress.AddressLine2
	}

	if updateAddress.City != "" {
		address.City = updateAddress.City
	}

	if updateAddress.State != "" {
		address.State = updateAddress.State
	}

	if updateAddress.Country != "" {
		address.Country = updateAddress.Country
	}

	return address
}

func MessageAddress(genericMessage string) string {
	return stringUtil.FormatGenericMessagesString(genericMessage, "Address")
}
