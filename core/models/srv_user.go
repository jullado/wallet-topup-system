package models

type SrvUserWalletModel struct {
	UserID   uint    `json:"user_id" example:"1"`
	Name     string  `json:"name" example:"Julladith Klinloy"`
	Balance  float64 `json:"balance" example:"100.00"`
	Currency string  `json:"currency" example:"THB"`
}
