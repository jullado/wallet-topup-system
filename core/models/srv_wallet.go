package models

import "time"

type SrvTopUpVerifiedResModel struct {
	TransactionID string    `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426655440000"`
	UserID        uint      `json:"user_id" example:"1"`
	Amount        float64   `json:"amount" example:"100.00"`
	PaymentMethod string    `json:"payment_method" example:"credit_card"`
	Status        string    `json:"status" example:"verified"`
	ExpiresAt     time.Time `json:"expires_at" example:"2025-05-14T00:00:00Z"`
}

type SrvTopUpConfirmedResModel struct {
	TransactionID string  `json:"transaction_id" example:"123e4567-e89b-12d3-a456-426655440000"`
	UserID        uint    `json:"user_id" example:"1"`
	Amount        float64 `json:"amount" example:"100.00"`
	Status        string  `json:"status" example:"confirmed"`
	Balance       float64 `json:"balance" example:"100.00"`
}
