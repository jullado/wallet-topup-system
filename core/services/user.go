package services

import "wallet-topup-system/core/models"

type UserService interface {
	// สำหรับสร้าง user เริ่มต้น
	Initialize() error

	// สำหรับดึงข้อมูล wallet ของ user
	GetUserWallet(userID uint) (result models.SrvUserWalletModel, err error)
}
