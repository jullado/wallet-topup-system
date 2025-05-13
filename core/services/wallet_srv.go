package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"
	"wallet-topup-system/common/cache"
	"wallet-topup-system/common/logs"
	"wallet-topup-system/config"
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/repositories"
	"wallet-topup-system/utils"
)

type walletSrv struct {
	log             logs.AppLog
	cache           cache.AppCache
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
}

func NewWalletService(
	log logs.AppLog,
	cache cache.AppCache,
	userRepo repositories.UserRepository,
	transactionRepo repositories.TransactionRepository,
) WalletService {
	return &walletSrv{
		log,
		cache,
		userRepo,
		transactionRepo,
	}
}

func (s *walletSrv) TopUpVerified(userID uint, amount float64, paymentMethod string) (result models.SrvTopUpVerifiedResModel, err error) {
	s.log.Info("top-up verified with user_id: " + fmt.Sprint(userID) + ", amount: " + fmt.Sprint(amount) + ", payment_method: " + paymentMethod)

	// check userID is exists
	_, err = s.userRepo.Get(userID)
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

	// check amount is positive
	if amount < 0 {
		s.log.Error(models.ErrAmountMustBePositive)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrAmountMustBePositive,
		}
	}

	// check payment method is valid
	if !slices.Contains(models.PaymentMethod, paymentMethod) {
		s.log.Error(models.ErrPaymentMethodIsInvalid)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrPaymentMethodIsInvalid,
		}
	}

	// create transaction
	now := time.Now()
	res, err := s.transactionRepo.Create(models.RepoTransactionModel{
		UserID:        userID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		CreatedAt:     now,
		UpdatedAt:     now,
		ExpiresAt:     now.Add(config.Env.TimeoutTopUp),
	})
	if err != nil {
		s.log.Error(err)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrUnexpectedError,
		}
	}

	// set cache transactionID
	byteRes, _ := json.Marshal(res)
	if err := s.cache.Set(res.TransactionID, byteRes, config.Env.TimeoutTopUp); err != nil {
		s.log.Error(err)
	}

	result = models.SrvTopUpVerifiedResModel{
		TransactionID: res.TransactionID,
		UserID:        res.UserID,
		Amount:        res.Amount,
		PaymentMethod: res.PaymentMethod,
		Status:        res.Status,
		ExpiresAt:     res.ExpiresAt,
	}

	return result, nil
}

func (s *walletSrv) TopUpConfirmed(transactionID string) (result models.SrvTopUpConfirmedResModel, err error) {
	s.log.Info("top-up confirmed with transaction_id: " + transactionID)

	// declare variable
	transaction := models.RepoTransactionResModel{}
	_result := &models.SrvTopUpConfirmedResModel{}

	// check transactionID is exists in cache
	res, err := s.cache.Get(transactionID)
	if err != nil {
		s.log.Error(err)
		return result, utils.ErrHandler{
			Code:    400,
			Message: models.ErrUnexpectedError,
		}
	}

	// if transactionID is not exists in cache -> get from database
	if len(res) == 0 {
		transaction, err = s.transactionRepo.Get(transactionID)
		if err != nil {
			if err.Error() == "record not found" {
				s.log.Error(models.ErrTransactionIsNotExist)
				return result, utils.ErrHandler{
					Code:    400,
					Message: models.ErrTransactionIsNotExist,
				}
			}
			s.log.Error(err)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			}
		}

		// ensure the transaction is verified
		if transaction.Status != models.TransactionStatusVerified {
			s.log.Error(models.ErrTransactionIsNotVerified)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsNotVerified,
			}
		}

		if transaction.ExpiresAt.Before(time.Now()) {
			s.log.Error(models.ErrTransactionIsExpired)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsExpired,
			}
		}
	} else {
		if err := json.Unmarshal(res, &transaction); err != nil {
			s.log.Error(err)
			return result, utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			}
		}
	}

	// update transaction with atomic operation
	err = s.transactionRepo.RunTransaction(func(rt repositories.TransactionRepository) error {

		// update transaction status to confirmed
		resTrans, err := rt.UpdateStatus(
			models.RepoFilterTransactionModel{
				TransactionID: transaction.TransactionID,
				Status:        models.TransactionStatusVerified,
			},
			models.TransactionStatusConfirmed,
		)
		if err != nil {
			s.log.Error(err)
			return errors.New(models.ErrUnexpectedError)
		}
		if len(resTrans) == 0 {
			s.log.Error(models.ErrTransactionIsNotExist)
			return errors.New(models.ErrTransactionIsNotExist)
		}

		// update wallet balance
		resWallet, err := rt.UpdateWalletBalance(
			models.RepoFilterWalletModel{
				UserID: transaction.UserID,
			},
			transaction.Amount,
		)
		if err != nil {
			s.log.Error(err)
			return errors.New(models.ErrUnexpectedError)
		}
		if len(resWallet) == 0 {
			s.log.Error(models.ErrWalletIsNotExist)
			return errors.New(models.ErrWalletIsNotExist)
		}

		// set result
		_result = &models.SrvTopUpConfirmedResModel{
			TransactionID: resTrans[0].TransactionID,
			UserID:        resTrans[0].UserID,
			Amount:        resTrans[0].Amount,
			Status:        resTrans[0].Status,
			Balance:       resWallet[0].Balance,
		}

		return nil
	})
	if err != nil {
		s.log.Error(err)
		return result, utils.ErrHandler{
			Code:    400,
			Message: err.Error(),
		}
	}

	// delete transaction from cache
	if err := s.cache.Delete(transactionID); err != nil {
		s.log.Error(err)
	}

	result = *_result

	return result, nil
}
