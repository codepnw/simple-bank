package transferhandler

type TransferReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1" example:"1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1" example:"2"`
	Amount        int64  `json:"amount" binding:"required,gt=0" example:"10"`
	Currency      string `json:"currency" binding:"required,oneof=THB USD" example:"THB"`
}
