package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/auth-service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	ContentTypeJSON = "application/json"
	ContentType     = "Content-Type"
	email           = "test@test.com"
)

func TestRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	url := "/refresh"

	t.Run("success", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(&user.TokenResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil)

		body := map[string]string{"refresh_token": "validRefreshToken"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Refresh(ctx)

		assert.Equal(t, http.StatusOK, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "Refresh token success")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})

	t.Run("refresh: error validate body", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(nil, "xx")

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Refresh(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "invalid character")

	})

	t.Run("refresh: error refresh token", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("RefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(&user.TokenResponse{}, errors.New("error"))

		body := map[string]string{"refresh_token": "validRefreshToken"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Refresh(ctx)

		// Assert that the mock expectations were met
		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "error")

		mockUserService.AssertExpectations(t)
	})
}

func TestProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	url := "/profile"

	t.Run("profile: success", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user.User{
			ID:    "1",
			Email: email,
		}, nil)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		token := user.RegisteredClaims{
			UserName: "test",
			Email:    email,
		}
		token.Subject = "1"
		ctx.Set("token", &token)
		handler := user.NewUserHandler(mockUserService, logger)
		handler.Profile(ctx)

		assert.Equal(t, http.StatusOK, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "Login success")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})

	t.Run("profile: error get user", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user.User{}, errors.New("error"))

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		token := user.RegisteredClaims{
			UserName: "test",
			Email:    email,
		}
		token.Subject = "1"
		ctx.Set("token", &token)
		handler := user.NewUserHandler(mockUserService, logger)
		handler.Profile(ctx)

		assert.Equal(t, http.StatusNotFound, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "not found")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	url := "/register"
	t.Run("register: success", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(user.User{
			ID:    "1",
			Email: email,
		}, nil)

		body := map[string]string{"username": "test", "password": "test", "email": email}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Register(ctx)

		assert.Equal(t, http.StatusCreated, respRecorder.Code)
		assert.NotNil(t, respRecorder.Body.String())
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})

	t.Run("register: error validate body", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(user.User{}, errors.New("error"))

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Register(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.NotNil(t, respRecorder.Body.String())

	})

	t.Run("register: error create user", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(user.User{}, errors.New("error"))

		body := map[string]string{"username": "test", "password": "test", "email": email}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Register(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "error")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})
}

func TestLoginV1(t *testing.T) {
	gin.SetMode(gin.TestMode)
	url := "/login"
	t.Run("login: success", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(&user.TokenResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil)

		body := map[string]string{"username": "test", "password": "test"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Login(ctx)

		assert.Equal(t, http.StatusOK, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "Login success")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})

	t.Run("login: error validate body", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Login(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.NotEmpty(t, respRecorder.Body.String())

	})

	t.Run("login: error login", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(&user.TokenResponse{}, errors.New("error"))

		body := map[string]string{"username": "test", "password": "test"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.Login(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "error")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})
}

func TestVerifyToken(t *testing.T) {
	url := "/verify"
	t.Run("verify: success", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(&user.TokenResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil)

		body := map[string]string{"access_token": "validAccessToken", "refresh_token": "validRefreshToken"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.VerifyToken(ctx)

		assert.Equal(t, http.StatusOK, respRecorder.Code)
		assert.NotNil(t, respRecorder.Body.String())
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})

	t.Run("verify: error validate body", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.VerifyToken(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.NotEmpty(t, respRecorder.Body.String())

	})

	t.Run("verify: error verify token", func(t *testing.T) {
		mockUserService := new(user.MockUserService)
		logger := slog.Default()
		mockUserService.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(&user.TokenResponse{}, errors.New("error"))

		body := map[string]string{"access_token": "validAccessToken", "refresh_token": "validRefreshToken"}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
		req.Header.Set(ContentType, ContentTypeJSON)

		respRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(respRecorder)
		ctx.Request = req

		handler := user.NewUserHandler(mockUserService, logger)
		handler.VerifyToken(ctx)

		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Contains(t, respRecorder.Body.String(), "error")
		// Assert that the mock expectations were met
		mockUserService.AssertExpectations(t)
	})
}
