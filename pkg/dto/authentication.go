package dto

type Authentication struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Login       string `json:"login"`
	RoleName    string `json:"role"`
	RoleID      int    `json:"roleId"`
	TokenString string `json:"token"`
}
