package constantUtils

const (
	//MESSAGES

	GenericNotFoundMessage            = "Any %_% was found"
	GenericCacheForbiddenMessage      = "Your token has expired or you already used it"
	GenericFoundSuccessMessage        = "%_% found with successful"
	GenericAlreadyExistsMessage       = "This %_% already exists on our database"
	GenericCreateErrorMessage         = "It was not possible to create this %_%"
	GenericCreateSuccessMessage       = "%_% created with successful"
	GenericUserCreatedSuccessMessage  = "User created with successful! Please confirm your email address to access our site"
	GenericUpdateErrorMessage         = "It was not possible to update this %_%"
	GenericUpdateSuccessMessage       = "%_% updated with successful"
	GenericDeleteErrorMessage         = "It was not possible to delete this %_%"
	GenericInternalServerErrorMessage = "It was not possible to perform this action"
	GenericInvalidFieldMessage        = "%_% cannot be null"
	CpfCnpjInvalidMessage             = "This CPF or CNPJ is not valid!"
	RoleInvalidMessage                = "This Role is not valid!"
	EmailInvalidMessage               = "This Email is not valid!"
	GenericTokenDoesNotMatch          = "This token does not match!"

	//STATUS

	StatusSuccess             = "success"
	StatusConflict            = "conflict"
	StatusInternalServerError = "internal_server_error"
	StatusNotFound            = "not_found"
	StatusBadRequest          = "bad_request"
	StatusForbidden           = "forbidden"

	//ROUTES

	PathUserIdParam           = "/:userId"
	PathAddressIdParam        = "/:addressId"
	PathUpdateUserMainAddress = "/:addressId/address/:userId/user"

	//CONDITIONS

	IdCondition     = "id = ?"
	UserIdCondition = "user_id = ?"
	EmailCondition = "email = ?"

	//URL
	GeneralUrlFormat = "%s/%s/%s"

)
