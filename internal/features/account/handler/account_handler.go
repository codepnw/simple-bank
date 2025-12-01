package accounthandler

import (
	"errors"

	"github.com/codepnw/simple-bank/internal/consts"
	"github.com/codepnw/simple-bank/internal/features/account"
	accountusecase "github.com/codepnw/simple-bank/internal/features/account/usecase"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/helper"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	uc accountusecase.AccountUsecase
}

func NewAccountHandler(uc accountusecase.AccountUsecase) *accountHandler {
	return &accountHandler{uc: uc}
}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	type currencyReq struct {
		Currency string `json:"currency" validate:"required,oneof=THB thb USD usd"`
	}

	req := new(currencyReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := helper.Validate(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.CreateAccount(c.Request.Context(), account.AccountCurrency(req.Currency))
	if err != nil {
		switch err {
		case errs.ErrNoUserID:
			response.Unauthorized(c, err.Error())
			return
		case errs.ErrInvalidCurrency:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrAccountNotFound:
			response.NotFound(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Created(c, "", data)
}

func (h *accountHandler) GetAccount(c *gin.Context) {
	id, err := helper.ParseInt64(c.Param(consts.ParamAccountID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.uc.GetAccount(c.Request.Context(), id)
	if err != nil {
		switch err {
		case errs.ErrAccountNotFound:
			response.NotFound(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Success(c, "", data)
}

func (h *accountHandler) ListAccounts(c *gin.Context) {
	page := helper.ParseInt(c.Query("page"))
	size := helper.ParseInt(c.Query("size"))

	data, err := h.uc.ListAccounts(c.Request.Context(), page, size)
	if err != nil {
		if errors.Is(err, errs.ErrNoUserID) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalServerError(c, err)
		return
	}
	response.Success(c, "", data)
}
