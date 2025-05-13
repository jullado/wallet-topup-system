package repositories

import "wallet-topup-system/core/models"

type UserRepository interface {
	// สำหรับสร้าง user เริ่มต้นของระบบ
	Initialize(payload models.RepoUserModel) (result models.RepoUserModel, err error)

	// สำหรับสร้าง user
	Create(payload models.RepoUserModel) (result models.RepoUserModel, err error)

	// สำหรับดึงข้อมูล user
	Get(userID uint) (result models.RepoUserModel, err error)

	// สำหรับดึงข้อมูล wallet ของ user
	GetWallet(userID uint) (result models.RepoWalletModel, err error)
}
