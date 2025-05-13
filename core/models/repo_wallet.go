package models

import (
	"time"

	"github.com/google/uuid"
)

type RepoWalletModel struct {
	WalletID  uuid.UUID `gorm:"column:wallet_id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uint      `gorm:"column:user_id;index;not null;constraint:OnDelete:CASCADE"`
	Balance   float64   `gorm:"column:balance;type:numeric(12,2);default:0.00"`
	Currency  string    `gorm:"column:currency;type:varchar(10)"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Users RepoUserModel `gorm:"foreignKey:UserID"`
}

func (RepoWalletModel) TableName() string {
	return "wallets"
}

type RepoWalletResModel struct {
	WalletID  string  `gorm:"column:wallet_id"`
	UserID    uint    `gorm:"column:user_id"`
	Balance   float64 `gorm:"column:balance"`
	Currency  string  `gorm:"column:currency"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RepoWalletResModel) TableName() string {
	return "wallets"
}

type RepoFilterWalletModel struct {
	WalletID string `gorm:"column:wallet_id"`
	UserID   uint   `gorm:"column:user_id"`
}

func (RepoFilterWalletModel) TableName() string {
	return "wallets"
}
