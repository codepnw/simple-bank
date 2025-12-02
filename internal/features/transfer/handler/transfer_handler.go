package transferhandler

import (
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	"github.com/codepnw/simple-bank/pkg/utils/errs"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type transferHandler struct {
	uc transferusecase.TransferUsecase
}

func NewTransferHandler(uc transferusecase.TransferUsecase) *transferHandler {
	return &transferHandler{uc: uc}
}

func (h *transferHandler) CreateTransfer(c *gin.Context) {
	req := new(TransferReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := &transferusecase.TransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}
	result, err := h.uc.Transfer(c.Request.Context(), input)
	if err != nil {
		switch err {
		case errs.ErrAccountNotFound:
			response.NotFound(c, err.Error())
			return
		case errs.ErrTransferToSelf:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrCurrencyMismatch:
			response.BadRequest(c, err.Error())
			return
		case errs.ErrMoneyNotEnough:
			response.BadRequest(c, err.Error())
			return
		default:
			response.InternalServerError(c, err)
			return
		}
	}
	response.Success(c, "transfer success", result)
}
