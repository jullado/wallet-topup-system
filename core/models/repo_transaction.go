package models

import (
	"time"

	"github.com/google/uuid"
)

type RepoTransactionModel struct {
	TransactionID uuid.UUID `gorm:"column:transaction_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uint      `gorm:"column:user_id;index:idx_user_status;not null;constraint:OnDelete:CASCADE"`
	Amount        float64   `gorm:"column:amount;type:numeric(12,2);default:0.00"`
	Currency      string    `gorm:"column:currency;type:varchar(10)"`
	PaymentMethod string    `gorm:"column:payment_method;not null"`
	Status        string    `gorm:"column:status;index:idx_user_status;not null;default:'verified'"`
	ExpiresAt     time.Time `gorm:"column:expires_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	Users RepoUserModel `gorm:"foreignKey:UserID;-"`
}

func (RepoTransactionModel) TableName() string {
	return "transactions"
}

type RepoTransactionResModel struct {
	TransactionID string    `gorm:"column:transaction_id"`
	UserID        uint      `gorm:"column:user_id"`
	Amount        float64   `gorm:"column:amount"`
	Currency      string    `gorm:"column:currency"`
	PaymentMethod string    `gorm:"column:payment_method"`
	Status        string    `gorm:"column:status"`
	ExpiresAt     time.Time `gorm:"column:expires_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (RepoTransactionResModel) TableName() string {
	return "transactions"
}

type RepoFilterTransactionModel struct {
	TransactionID string `gorm:"column:transaction_id"`
	UserID        uint   `gorm:"column:user_id"`
	PaymentMethod string `gorm:"column:payment_method"`
	Status        string `gorm:"column:status"`
}

func (RepoFilterTransactionModel) TableName() string {
	return "transactions"
}
