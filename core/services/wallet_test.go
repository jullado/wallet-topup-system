package services_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
	"wallet-topup-system/common/cache"
	"wallet-topup-system/common/logs"
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/repositories"
	"wallet-topup-system/core/services"
	"wallet-topup-system/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopUpVerified(t *testing.T) {
	// input
	type arg struct {
		userID        uint
		amount        float64
		paymentMethod string
	}

	// mock data
	now := time.Now()
	mockTransactionID := "123e4567-e89b-12d3-a456-426655440000"

	// test case
	testcases := []struct {
		name       string
		arg        arg
		wantResult models.SrvTopUpVerifiedResModel
		wantErr    error
	}{
		{
			name: "success",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{
				TransactionID: mockTransactionID,
				UserID:        1,
				Amount:        100.20,
				PaymentMethod: "credit_card",
				Status:        "verified",
				ExpiresAt:     now.Add(1 * time.Minute),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUserIDIsNotExist,
			},
		},
		{
			name: "get user unexpected error",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "amount must be positive",
			arg: arg{
				userID:        1,
				amount:        -100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrAmountMustBePositive,
			},
		},
		{
			name: "payment method is invalid",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "invalid",
			},
			wantResult: models.SrvTopUpVerifiedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrPaymentMethodIsInvalid,
			},
		},
		{
			name: "create transaction unexpected error",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "set cache unexpected error",
			arg: arg{
				userID:        1,
				amount:        100.20,
				paymentMethod: "credit_card",
			},
			wantResult: models.SrvTopUpVerifiedResModel{
				TransactionID: mockTransactionID,
				UserID:        1,
				Amount:        100.20,
				PaymentMethod: "credit_card",
				Status:        "verified",
				ExpiresAt:     now.Add(1 * time.Minute),
			},
			wantErr: nil,
		},
	}

	// -------------------- run test --------------------

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()
			transactionRepo := repositories.NewTransactionRepoMock()
			log := logs.NewAppLogsMock()
			cache := cache.NewAppCacheMock()

			// mock Get user
			switch tt.name {
			case "get user unexpected error":
				userRepo.On("Get", tt.arg.userID).Return(nil, errors.New(""))
			case "user not found":
				userRepo.On("Get", tt.arg.userID).Return(nil, errors.New("record not found"))
			default:
				userRepo.On("Get", tt.arg.userID).Return(models.RepoUserModel{UserID: tt.arg.userID, Name: "root"}, nil)
			}

			// mock Create transaction
			switch tt.name {
			case "create transaction unexpected error":
				transactionRepo.On("Create", mock.AnythingOfType("models.RepoTransactionModel")).Return(nil, errors.New(""))
			default:
				transactionRepo.On("Create", mock.AnythingOfType("models.RepoTransactionModel")).
					Return(models.RepoTransactionResModel{
						TransactionID: mockTransactionID,
						UserID:        tt.arg.userID,
						Amount:        tt.arg.amount,
						PaymentMethod: tt.arg.paymentMethod,
						Status:        "verified",
						CreatedAt:     now,
						UpdatedAt:     now,
						ExpiresAt:     now.Add(1 * time.Minute),
					}, nil)
			}

			// mock Set cache
			switch tt.name {
			case "set cache unexpected error":
				cache.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(errors.New(""))
			default:
				cache.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)
			}

			// ------------------- Act (ดําเนินการ) --------------------
			walletSrv := services.NewWalletService(log, cache, userRepo, transactionRepo)

			result, err := walletSrv.TopUpVerified(tt.arg.userID, tt.arg.amount, tt.arg.paymentMethod)

			// ------------------- Assert (ตรวจสอบ) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestTopUpConfirmed(t *testing.T) {
	// input
	type arg struct {
		transactionID string
	}

	// mock data
	now := time.Now()
	mockTransaction := models.RepoTransactionModel{
		TransactionID: uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"),
		UserID:        1,
		Amount:        100.20,
		PaymentMethod: "credit_card",
		Status:        "verified",
		CreatedAt:     now,
		UpdatedAt:     now,
		ExpiresAt:     now.Add(1 * time.Minute),
	}
	mockBalance := 100.20
	mockTransactionCache, _ := json.Marshal(mockTransaction)

	// test case
	testcases := []struct {
		name       string
		arg        arg
		wantResult models.SrvTopUpConfirmedResModel
		wantErr    error
	}{
		{
			name: "success with cache",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{
				TransactionID: mockTransaction.TransactionID.String(),
				UserID:        mockTransaction.UserID,
				Amount:        mockTransaction.Amount,
				Status:        models.TransactionStatusConfirmed,
				Balance:       mockBalance + mockTransaction.Amount,
			},
			wantErr: nil,
		},
		{
			name: "success without cache",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{
				TransactionID: mockTransaction.TransactionID.String(),
				UserID:        mockTransaction.UserID,
				Amount:        mockTransaction.Amount,
				Status:        models.TransactionStatusConfirmed,
				Balance:       mockBalance + mockTransaction.Amount,
			},
			wantErr: nil,
		},
		{
			name: "transaction not found",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsNotExist,
			},
		},
		{
			name: "transaction not verified",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsNotVerified,
			},
		},
		{
			name: "transaction expired",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsExpired,
			},
		},
		{
			name: "transaction already confirmed",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrTransactionIsNotExist,
			},
		},
		{
			name: "unmarshal cache unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "get cache unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "get transaction unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "update transaction unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "update wallet balance unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "wallet not found",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrWalletIsNotExist,
			},
		},
		{
			name: "delete cache unexpected error",
			arg: arg{
				transactionID: mockTransaction.TransactionID.String(),
			},
			wantResult: models.SrvTopUpConfirmedResModel{
				TransactionID: mockTransaction.TransactionID.String(),
				UserID:        mockTransaction.UserID,
				Amount:        mockTransaction.Amount,
				Status:        models.TransactionStatusConfirmed,
				Balance:       mockBalance + mockTransaction.Amount,
			},
			wantErr: nil,
		},
	}

	// -------------------- run test --------------------

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()
			transactionRepo := repositories.NewTransactionRepoMock()
			log := logs.NewAppLogsMock()
			cache := cache.NewAppCacheMock()

			// mock Get cache
			switch tt.name {
			case "unmarshal cache unexpected error":
				cache.On("Get", tt.arg.transactionID).Return([]byte("xxx"), nil)
			case "success without cache", "get transaction unexpected error",
				"transaction not found", "transaction not verified", "transaction expired":
				cache.On("Get", tt.arg.transactionID).Return(nil, nil)
			case "get cache unexpected error":
				cache.On("Get", tt.arg.transactionID).Return(nil, errors.New(""))
			default:
				cache.On("Get", tt.arg.transactionID).Return(mockTransactionCache, nil)
			}

			// mock Get transaction
			switch tt.name {
			case "transaction not found":
				transactionRepo.On("Get", tt.arg.transactionID).Return(models.RepoTransactionResModel{}, errors.New("record not found"))
			case "get transaction unexpected error":
				transactionRepo.On("Get", tt.arg.transactionID).Return(models.RepoTransactionResModel{}, errors.New(""))
			case "transaction expired":
				transactionRepo.On("Get", tt.arg.transactionID).Return(models.RepoTransactionResModel{
					TransactionID: mockTransaction.TransactionID.String(),
					UserID:        mockTransaction.UserID,
					Amount:        mockTransaction.Amount,
					Currency:      mockTransaction.Currency,
					PaymentMethod: mockTransaction.PaymentMethod,
					Status:        mockTransaction.Status,
					CreatedAt:     now,
					UpdatedAt:     now,
					ExpiresAt:     now,
				}, nil)
			case "transaction not verified":
				transactionRepo.On("Get", tt.arg.transactionID).Return(models.RepoTransactionResModel{
					TransactionID: mockTransaction.TransactionID.String(),
					UserID:        mockTransaction.UserID,
					Amount:        mockTransaction.Amount,
					Currency:      mockTransaction.Currency,
					PaymentMethod: mockTransaction.PaymentMethod,
					Status:        "xxx",
					CreatedAt:     now,
					UpdatedAt:     now,
					ExpiresAt:     now.Add(1 * time.Minute),
				}, nil)
			default:
				transactionRepo.On("Get", tt.arg.transactionID).Return(models.RepoTransactionResModel{
					TransactionID: mockTransaction.TransactionID.String(),
					UserID:        mockTransaction.UserID,
					Amount:        mockTransaction.Amount,
					Currency:      mockTransaction.Currency,
					PaymentMethod: mockTransaction.PaymentMethod,
					Status:        mockTransaction.Status,
					CreatedAt:     now,
					UpdatedAt:     now,
					ExpiresAt:     now.Add(1 * time.Minute),
				}, nil)
			}

			// mock UpdateStatus
			switch tt.name {
			case "transaction already confirmed":
				transactionRepo.On("UpdateStatus", mock.Anything, models.TransactionStatusConfirmed).Return([]models.RepoTransactionResModel{}, nil)
			case "update transaction unexpected error":
				transactionRepo.On("UpdateStatus", mock.Anything, models.TransactionStatusConfirmed).Return([]models.RepoTransactionResModel{}, errors.New(""))
			default:
				transactionRepo.On("UpdateStatus", mock.Anything, models.TransactionStatusConfirmed).Return([]models.RepoTransactionResModel{
					{
						TransactionID: mockTransaction.TransactionID.String(),
						UserID:        mockTransaction.UserID,
						Amount:        mockTransaction.Amount,
						Currency:      mockTransaction.Currency,
						PaymentMethod: mockTransaction.PaymentMethod,
						Status:        models.TransactionStatusConfirmed,
						CreatedAt:     now,
						UpdatedAt:     now,
						ExpiresAt:     now.Add(1 * time.Minute),
					},
				}, nil)
			}

			// mock UpdateWalletBalance
			switch tt.name {
			case "wallet not found":
				transactionRepo.On("UpdateWalletBalance", mock.Anything, mockTransaction.Amount).Return([]models.RepoWalletResModel{}, nil)
			case "update wallet balance unexpected error":
				transactionRepo.On("UpdateWalletBalance", mock.Anything, mockTransaction.Amount).Return([]models.RepoWalletResModel{}, errors.New(""))
			default:
				transactionRepo.On("UpdateWalletBalance", mock.Anything, mockTransaction.Amount).Return([]models.RepoWalletResModel{
					{
						WalletID:  mockTransaction.TransactionID.String(),
						UserID:    mockTransaction.UserID,
						Balance:   mockBalance + mockTransaction.Amount,
						Currency:  mockTransaction.Currency,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			}

			// mock RunTransaction
			switch tt.name {
			case "run transaction unexpected error":
				transactionRepo.On("RunTransaction", mock.Anything, mock.Anything).Return(errors.New(""))
			default:
				transactionRepo.On("RunTransaction", mock.Anything).Return(nil)
			}

			// mock Delete cache
			switch tt.name {
			case "delete cache unexpected error":
				cache.On("Delete", tt.arg.transactionID).Return(errors.New(""))
			default:
				cache.On("Delete", tt.arg.transactionID).Return(nil)
			}

			// ------------------- Act (ดําเนินการ) --------------------
			walletSrv := services.NewWalletService(log, cache, userRepo, transactionRepo)

			result, err := walletSrv.TopUpConfirmed(tt.arg.transactionID)

			// ------------------- Assert (ตรวจสอบ) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}
