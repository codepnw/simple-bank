package accounthandler

type CreateAccountReq struct {
	Currency string `json:"currency" binding:"required,oneof=THB thb USD usd" example:"THB"`
}
