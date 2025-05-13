package services_test

import (
	"errors"
	"testing"
	"wallet-topup-system/common/logs"
	"wallet-topup-system/core/models"
	"wallet-topup-system/core/repositories"
	"wallet-topup-system/core/services"
	"wallet-topup-system/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitialize(t *testing.T) {

	// test case
	testcases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name: "database error",
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
	}

	// -------------------- run test --------------------

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Initialize
			switch tt.name {
			case "database error":
				userRepo.On("Initialize", mock.AnythingOfType("models.RepoUserModel")).Return(nil, errors.New(""))
			default:
				userRepo.On("Initialize", mock.AnythingOfType("models.RepoUserModel")).Return(models.RepoUserModel{UserID: 1, Name: "root"}, nil)
			}

			// ------------------- Act (ดําเนินการ) --------------------
			userSrv := services.NewUserService(logs.NewAppLogsMock(), userRepo)

			err := userSrv.Initialize()

			// ------------------- Assert (ตรวจสอบ) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func TestGetUserWallet(t *testing.T) {
	// input
	type arg struct {
		userID uint
	}

	// test case
	testcases := []struct {
		name       string
		arg        arg
		wantResult models.SrvUserWalletModel
		wantErr    error
	}{
		{
			name: "success",
			arg: arg{
				userID: 1,
			},
			wantResult: models.SrvUserWalletModel{
				UserID:   1,
				Name:     "root",
				Balance:  100.00,
				Currency: "THB",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			arg: arg{
				userID: 0,
			},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUserIDIsNotExist,
			},
		},
		{
			name: "wallet not found",
			arg: arg{
				userID: 0,
			},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrWalletIsNotExist,
			},
		},
		{
			name: "get user unexpected error",
			arg: arg{
				userID: 1,
			},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
		{
			name: "get user wallet unexpected error",
			arg: arg{
				userID: 1,
			},
			wantErr: utils.ErrHandler{
				Code:    400,
				Message: models.ErrUnexpectedError,
			},
		},
	}

	// -------------------- run test --------------------

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()
			log := logs.NewAppLogsMock()

			// mock Get user
			switch tt.name {
			case "get user unexpected error":
				userRepo.On("Get", tt.arg.userID).Return(nil, errors.New(""))
			case "user not found":
				userRepo.On("Get", tt.arg.userID).Return(nil, errors.New("record not found"))
			default:
				userRepo.On("Get", tt.arg.userID).Return(models.RepoUserModel{UserID: 1, Name: "root"}, nil)
			}

			// mock Get user wallet
			switch tt.name {
			case "get user wallet unexpected error":
				userRepo.On("GetWallet", tt.arg.userID).Return(nil, errors.New(""))
			case "user not found", "wallet not found":
				userRepo.On("GetWallet", tt.arg.userID).Return(nil, errors.New("record not found"))
			default:
				userRepo.On("GetWallet", tt.arg.userID).Return(models.RepoWalletModel{UserID: 1, Balance: 100.00}, nil)
			}

			// ------------------- Act (ดําเนินการ) --------------------
			userSrv := services.NewUserService(log, userRepo)

			result, err := userSrv.GetUserWallet(tt.arg.userID)

			// ------------------- Assert (ตรวจสอบ) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}
