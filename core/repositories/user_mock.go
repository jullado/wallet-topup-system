package repositories

import (
	"wallet-topup-system/core/models"

	"github.com/stretchr/testify/mock"
)

type userRepoMock struct {
	mock.Mock
}

func NewUserRepoMock() *userRepoMock {
	return &userRepoMock{}
}

func (m *userRepoMock) Initialize(payload models.RepoUserModel) (result models.RepoUserModel, err error) {
	args := m.Called(payload)
	res, ok := args.Get(0).(models.RepoUserModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *userRepoMock) Create(payload models.RepoUserModel) (result models.RepoUserModel, err error) {
	args := m.Called(payload)
	res, ok := args.Get(0).(models.RepoUserModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *userRepoMock) Get(userID uint) (result models.RepoUserModel, err error) {
	args := m.Called(userID)
	res, ok := args.Get(0).(models.RepoUserModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *userRepoMock) GetWallet(userID uint) (result models.RepoWalletModel, err error) {
	args := m.Called(userID)
	res, ok := args.Get(0).(models.RepoWalletModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}
