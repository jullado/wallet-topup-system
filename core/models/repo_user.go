package models

import "time"

type RepoUserModel struct {
	UserID    uint   `gorm:"column:user_id;primaryKey"`
	Name      string `gorm:"column:name;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RepoUserModel) TableName() string {
	return "users"
}
