package redis

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRedis struct {
	mock.Mock
}

// Mock implementation of the Ping method
func (m *MockRedis) Ping(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// Mock implementation of the Set method
func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// Mock implementation of the Get method
func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// Mock implementation of the Del method
func (m *MockRedis) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Mock implementation of the Exists method
func (m *MockRedis) Exists(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

// Mock implementation of the SetEx method
func (m *MockRedis) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Different implementation from Set
	args := m.Called(ctx, key, value, expiration)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

// Mock implementation of the Close method
func (m *MockRedis) Close() error {
	args := m.Called()
	return args.Error(0)
}
