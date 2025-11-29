package userhandler

type RegisterReq struct {
	Username  string `json:"username" validate:"required,min=4"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type LoginReq struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}
