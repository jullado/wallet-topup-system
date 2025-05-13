package services

import (
	"fmt"
	"wallet-topup-system/common/logs"
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/repositories"
	"wallet-topup-system/utils"
)

type userSrv struct {
	log      logs.AppLog
	userRepo repositories.UserRepository
}

func NewUserService(log logs.AppLog, userRepo repositories.UserRepository) UserService {
	return &userSrv{
		log,
		userRepo,
	}
}

func (s *userSrv) Initialize() error {

	// initialize user
	result, err := s.userRepo.Initialize(models.RepoUserModel{UserID: 1, Name: "root"})
	if err != nil {
		s.log.Error(err)
		return utils.ErrHandler{
			Code:    400,
			Message: models.ErrUnexpectedError,
		}
	}

	if result.UserID > 0 {
		s.log.Info("Initialize user success with userID: " + fmt.Sprint(result.UserID))
	}

	return nil
}

func (s *userSrv) GetUserWallet(userID uint) (result models.SrvUserWalletModel, err error) {
	s.log.Info("get user wallet with user_id: " + fmt.Sprint(userID))

	// get user data
	resUser, err := s.userRepo.Get(userID)
	if err != nil {
		if err.Error() == "record not found" {
			s.log.Error(models.ErrUserIDIsNotExist)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrUserIDIsNotExist,
			}
		}
		s.log.Error(err)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrUnexpectedError,
		}
	}

	// get wallet data
	resWallet, err := s.userRepo.GetWallet(userID)
	if err != nil {
		if err.Error() == "record not found" {
			s.log.Error(models.ErrWalletIsNotExist)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrWalletIsNotExist,
			}
		}
		s.log.Error(err)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrUnexpectedError,
		}
	}

	// set default currency
	if resWallet.Currency == "" {
		resWallet.Currency = "THB"
	}

	result = models.SrvUserWalletModel{
		UserID:   resUser.UserID,
		Name:     resUser.Name,
		Balance:  resWallet.Balance,
		Currency: resWallet.Currency,
	}

	return result, nil
}
