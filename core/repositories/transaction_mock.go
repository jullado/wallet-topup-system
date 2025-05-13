package repositories

import (
	"wallet-topup-system/core/models"

	"github.com/stretchr/testify/mock"
)

type transactionRepoMock struct {
	mock.Mock
}

func NewTransactionRepoMock() *transactionRepoMock {
	return &transactionRepoMock{}
}

func (m *transactionRepoMock) RunTransaction(fn func(rt TransactionRepository) error) error {
	return fn(m)
}

func (m *transactionRepoMock) Create(payload models.RepoTransactionModel) (result models.RepoTransactionResModel, err error) {
	args := m.Called(payload)
	res, ok := args.Get(0).(models.RepoTransactionResModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *transactionRepoMock) Get(transactionID string) (result models.RepoTransactionResModel, err error) {
	args := m.Called(transactionID)
	res, ok := args.Get(0).(models.RepoTransactionResModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *transactionRepoMock) UpdateStatus(filter models.RepoFilterTransactionModel, status string) (result []models.RepoTransactionResModel, err error) {
	args := m.Called(filter, status)
	res, ok := args.Get(0).([]models.RepoTransactionResModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *transactionRepoMock) UpdateWalletBalance(filter models.RepoFilterWalletModel, amount float64) (result []models.RepoWalletResModel, err error) {
	args := m.Called(filter, amount)
	res, ok := args.Get(0).([]models.RepoWalletResModel)
	if !ok {
		return result, args.Error(1)
	}
	return res, args.Error(1)
}
