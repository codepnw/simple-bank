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

// @Summary Create Account
// @Description user create account
// @Tags accounts
// @Accept       json
// @Produce      json
// @Param request body CreateAccountReq true "Create Account Data"
// @Success 201 {object} account.Account "Create Account Successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid Input"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Account Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router /accounts [post]
func (h *accountHandler) CreateAccount(c *gin.Context) {
	req := new(CreateAccountReq)
	if err := c.ShouldBindJSON(req); err != nil {
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

// @Summary Get Account
// @Description get account by id
// @Tags accounts
// @Accept       json
// @Produce      json
// @Param id path int true "Account ID"
// @Success 200 {object} account.Account "Get Account Successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid Input"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Account Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router /accounts/{id} [get]
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

// @Summary List Accounts
// @Description list account by user
// @Tags accounts
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} account.Account "List Accounts Successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Security     BearerAuth
// @Router /accounts [get]
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
