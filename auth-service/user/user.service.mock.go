package user

import (
	"context"
	"log/slog"

	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func NewProductServiceMock() *ProductServiceMock {
	return &ProductServiceMock{}
}

func (m *ProductServiceMock) CreateUser(ctx context.Context, logger *slog.Logger, body User) (User, error) {
	ret := m.Called(ctx, logger, body)
	return ret.Get(0).(User), ret.Error(1)
}

func (m *ProductServiceMock) GetUser(ctx context.Context, logger *slog.Logger, id string) (User, error) {
	ret := m.Called(ctx, logger, id)
	return ret.Get(0).(User), ret.Error(1)
}

func (m *ProductServiceMock) Login(ctx context.Context, logger *slog.Logger, body Login) (*TokenResponse, error) {
	ret := m.Called(ctx, logger, body)
	return ret.Get(0).(*TokenResponse), ret.Error(1)
}

func (m *ProductServiceMock) RefreshToken(ctx context.Context, logger *slog.Logger, token string) (*TokenResponse, error) {
	ret := m.Called(ctx, logger, token)
	return ret.Get(0).(*TokenResponse), ret.Error(1)
}

func (m *ProductServiceMock) VerifyAccessToken(logger *slog.Logger, token string) (*TokenResponse, error) {
	ret := m.Called(logger, token)
	return ret.Get(0).(*TokenResponse), ret.Error(1)
}
