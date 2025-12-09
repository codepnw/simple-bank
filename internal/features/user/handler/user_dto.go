package userhandler

type RegisterReq struct {
	Username  string `json:"username" binding:"required,min=4" example:"johndoe"`
	Password  string `json:"password" binding:"required,min=6" example:"pass1234"`
	FirstName string `json:"first_name" binding:"required" example:"john"`
	LastName  string `json:"last_name" binding:"required" example:"doe"`
	Email     string `json:"email" binding:"required,email" example:"john@mail.com"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email" example:"user1@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
