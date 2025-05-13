package models

// error message
var (
	ErrUnexpectedError          = "unexpected error"
	ErrUserIDIsNotExist         = "user_id is not exist"
	ErrAmountMustBePositive     = "amount must be positive"
	ErrTransactionIsNotVerified = "transaction is not verified"
	ErrTransactionIsExpired     = "transaction is expired"
	ErrTransactionIsNotExist    = "transaction is not exist"
	ErrWalletIsNotExist         = "wallet is not exist"
	ErrPaymentMethodIsInvalid   = "payment method is invalid"
	ErrUnauthorized             = "unauthorized"
)

// payment method
var PaymentMethod = []string{"credit_card"}

// transaction status
var (
	TransactionStatusVerified  = "verified"
	TransactionStatusConfirmed = "confirmed"
)
