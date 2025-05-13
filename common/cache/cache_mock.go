package cache

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type appCacheMock struct {
	mock.Mock
}

func NewAppCacheMock() *appCacheMock {
	return &appCacheMock{}
}

func (m *appCacheMock) Get(key string) ([]byte, error) {
	args := m.Called(key)
	res, ok := args.Get(0).([]byte)
	if !ok {
		return nil, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *appCacheMock) Set(key string, val []byte, exp time.Duration) error {
	return m.Called(key, val, exp).Error(0)
}

func (m *appCacheMock) Delete(key string) error {
	return m.Called(key).Error(0)
}

func (m *appCacheMock) ExpiredEvent(callback func(key string) error) error {
	return m.Called(callback).Error(0)
}
