package services

import "wallet-topup-system/core/models"

type WalletService interface {
	// สำหรับส่งข้อมูลการเติมเงินเข้า wallet
	TopUpVerified(userID uint, amount float64, paymentMethod string) (result models.SrvTopUpVerifiedResModel, err error)

	// สำหรับยืนยันการเติมเงิน
	TopUpConfirmed(transactionID string) (result models.SrvTopUpConfirmedResModel, err error)
}
