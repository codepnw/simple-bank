package userhandler

import (
	"errors"

	"github.com/codepnw/simple-bank/internal/features/user"
	userusecase "github.com/codepnw/simple-bank/internal/features/user/usecase"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(uc userusecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

// @Summary Register
// @Description user register
// @Tags users
// @Accept       json
// @Produce      json
// @Param request body RegisterReq true "User Registration Data"
// @Success 201 {object} userusecase.TokenResponse "User created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (h *userHandler) Register(c *gin.Context) {
	req := new(RegisterReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := &user.User{
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	data, err := h.uc.Register(c.Request.Context(), input)
	if err != nil {
		switch err {
		case errs.ErrEmailAlreadyExists:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrUsernameAlreadyExists:
			response.BadRequest(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Created(c, "", data)
}

// @Summary Login
// @Description user login
// @Tags users
// @Accept       json
// @Produce      json
// @Param request body LoginReq true "User Login Data"
// @Success 200 {object} userusecase.TokenResponse "User login successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (h *userHandler) Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.Success(c, "", data)
}

// @Summary Refresh Token
// @Description user refresh token
// @Tags users
// @Accept       json
// @Produce      json
// @Param request body RefreshTokenReq true "User Refresh Token Data"
// @Success 200 {object} userusecase.TokenResponse "Refresh Token successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "Token Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router /users/refresh-token [post]
func (h *userHandler) RefreshToken(c *gin.Context) {
	req := new(RefreshTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case errs.ErrTokenNotFound:
			response.NotFound(c, err.Error())
			return
		case errs.ErrTokenRevoked:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrTokenExpires:
			response.BadRequest(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Success(c, "", data)
}

// @Summary Logout
// @Description user logout
// @Tags users
// @Accept       json
// @Produce      json
// @Param request body RefreshTokenReq true "User Logout Data"
// @Success 204 {object} response.NoContentResponse "successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 404 {object} response.ErrorResponse "Token Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router /users/logout [post]
func (h *userHandler) Logout(c *gin.Context) {
	req := new(RefreshTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.uc.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, errs.ErrTokenNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.NoContent(c)
}
