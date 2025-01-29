package request

type UserCreateRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email" validate:"required,email"`
}

type UserUpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,password"`
}
