package models

type HandTopUpVerifiedReqModel struct {
	UserID        uint    `json:"user_id" example:"1"`
	Amount        float64 `json:"amount" example:"100.00"`
	PaymentMethod string  `json:"payment_method" example:"credit_card"`
}

type HandTopUpConfirmedReqModel struct {
	TransactionID string `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426655440000"`
}
