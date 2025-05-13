package repositories

import "wallet-topup-system/core/models"

type TransactionRepository interface {
	// สำหรับเรียกใช้งาน transaction แบบ atomic
	RunTransaction(fn func(rt TransactionRepository) error) error

	// สำหรับสร้าง transaction
	Create(payload models.RepoTransactionModel) (result models.RepoTransactionResModel, err error)

	// สำหรับดึงข้อมูล transaction ด้วย transactionID
	Get(transactionID string) (result models.RepoTransactionResModel, err error)

	// สำหรับอัพเดท status transaction
	UpdateStatus(filter models.RepoFilterTransactionModel, status string) (result []models.RepoTransactionResModel, err error)

	// สำหรับอัพเดท wallet balance
	UpdateWalletBalance(filter models.RepoFilterWalletModel, amount float64) (result []models.RepoWalletResModel, err error)
}
