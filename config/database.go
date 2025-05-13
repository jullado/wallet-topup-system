package config

import (
	"wallet-topup-system/core/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewAppInitDBPostgres() *gorm.DB {
	// db connection
	db, err := gorm.Open(
		postgres.Open(Env.DBURI),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		panic(err)
	}

	// create uuid extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic(err)
	}

	// create table if not exist
	if err := db.AutoMigrate(&models.RepoUserModel{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.RepoWalletModel{}, &models.RepoTransactionModel{}); err != nil {
		panic(err)
	}

	return db
}
