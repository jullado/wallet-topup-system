package repositories

import (
	"context"
	"time"
	"wallet-topup-system/core/models"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Initialize(payload models.RepoUserModel) (result models.RepoUserModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_result := &models.RepoUserModel{}

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// create initial user
		result := tx.FirstOrCreate(&payload)
		if result.Error != nil {
			return result.Error
		}

		// if user not already exist -> create user's wallet
		if result.RowsAffected > 0 {
			_result = &payload

			wallet := models.RepoWalletModel{
				UserID: payload.UserID,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	result = *_result

	return result, nil
}

func (r *userRepo) Create(payload models.RepoUserModel) (result models.RepoUserModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// create user
		result := tx.Create(&payload)
		if result.Error != nil {
			return result.Error
		}

		// if user not already exist -> create user's wallet
		if result.RowsAffected > 0 {
			wallet := models.RepoWalletModel{
				UserID: payload.UserID,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	result = payload

	return result, nil
}

func (r *userRepo) Get(userID uint) (result models.RepoUserModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Take(&result, userID).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *userRepo) GetWallet(userID uint) (result models.RepoWalletModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).First(&result, "user_id = ?", userID).Error; err != nil {
		return result, err
	}

	return result, nil
}
