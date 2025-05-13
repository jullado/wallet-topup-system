package repositories

import (
	"context"
	"time"
	"wallet-topup-system/core/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{
		db: db,
	}
}

func (r *transactionRepo) RunTransaction(fn func(rt TransactionRepository) error) error {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewTransactionRepository(tx))
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepo) Create(payload models.RepoTransactionModel) (result models.RepoTransactionResModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.db.WithContext(ctx).Clauses(clause.Returning{}).
		Create(&payload).
		Scan(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *transactionRepo) Get(transactionID string) (result models.RepoTransactionResModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *transactionRepo) UpdateStatus(filter models.RepoFilterTransactionModel, status string) (result []models.RepoTransactionResModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	model := []models.RepoTransactionModel{}

	err = r.db.WithContext(ctx).Model(&model).
		Clauses(clause.Locking{Strength: "UPDATE"}, clause.Returning{}).
		Where(&filter).
		Update("status", status).Error
	if err != nil {
		return result, err
	}

	// DTO
	result = make([]models.RepoTransactionResModel, 0, len(model))
	for _, v := range model {
		result = append(result, models.RepoTransactionResModel{
			TransactionID: v.TransactionID.String(),
			UserID:        v.UserID,
			Amount:        v.Amount,
			Currency:      v.Currency,
			PaymentMethod: v.PaymentMethod,
			Status:        v.Status,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
			ExpiresAt:     v.ExpiresAt,
		})
	}

	return result, nil
}

func (r *transactionRepo) UpdateWalletBalance(filter models.RepoFilterWalletModel, amount float64) (result []models.RepoWalletResModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	model := []models.RepoWalletModel{}

	err = r.db.WithContext(ctx).Model(&model).
		Clauses(clause.Locking{Strength: "UPDATE"}, clause.Returning{}).
		Where(&filter).Limit(1).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return result, err
	}

	// DTO
	result = make([]models.RepoWalletResModel, 0, len(model))
	for _, v := range model {
		result = append(result, models.RepoWalletResModel{
			WalletID:  v.WalletID.String(),
			UserID:    v.UserID,
			Balance:   v.Balance,
			Currency:  v.Currency,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return result, nil
}
