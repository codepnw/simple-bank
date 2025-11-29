package userhandler

import (
	"errors"

	"github.com/codepnw/simple-bank/internal/features/user"
	userusecase "github.com/codepnw/simple-bank/internal/features/user/usecase"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/codepnw/simple-bank/pkg/utils/validate"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(uc userusecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) Register(c *gin.Context) {
	req := new(RegisterReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
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

func (h *userHandler) Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
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
