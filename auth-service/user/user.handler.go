package user

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/auth-service/mlog"
)

type userHandler struct {
	userService UserService
	logger      *slog.Logger
}

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Profile(c *gin.Context)
	VerifyToken(c *gin.Context)
	Refresh(c *gin.Context)
}

func NewUserHandler(userService UserService, logger *slog.Logger) UserHandler {
	return &userHandler{userService, logger}
}

func (u *userHandler) Refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Refresh token")

	var body RefreshTokenResponse
	var response Response[*TokenResponse]

	if err := c.BindJSON(&body); err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := u.userService.RefreshToken(ctx, logger, body.RefreshToken)
	if err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Message = "Refresh token success"
	response.Status = "success"
	response.Data = result

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) Profile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Get user profile")

	claims := c.MustGet("token").(*RegisteredClaims)

	response := Response[IUser]{}

	result, err := u.userService.GetUser(ctx, logger, claims.Subject)
	if err != nil {
		logger.Error(err.Error())
		response.Message = "not found"
		response.Status = "error"

		c.JSON(http.StatusNotFound, response)
		return
	}

	response.Message = "Login success"
	response.Status = "success"

	data := IUser{
		ID:       result.ID,
		Href:     c.Request.URL.Path + "/" + result.ID,
		Username: result.Username,
		Email:    result.Email,
		Name:     result.Name,
		Roles:    result.Roles,
	}

	response.Data = data

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Create user")

	var body User
	response := Response[IUser]{}

	if err := c.BindJSON(&body); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := u.userService.CreateUser(ctx, logger, body)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := IUser{
		ID:       result.ID,
		Href:     c.Request.URL.Path + "/" + result.ID,
		Username: result.Username,
		Email:    result.Email,
		Name:     result.Name,
		Roles:    result.Roles,
	}

	response.Message = "Create user success"
	response.Status = "success"
	response.Data = data

	c.JSON(http.StatusCreated, response)
}

func (u *userHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Login user")

	var body Login
	var response Response[*TokenResponse]
	if err := c.BindJSON(&body); err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"
		c.JSON(400, response)
		return
	}

	result, err := u.userService.Login(ctx, logger, body)
	if err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"
		c.JSON(400, response)
		return
	}

	response.Message = "Login success"
	response.Status = "success"
	response.Data = result

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) VerifyToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Verify token")

	var body TokenResponse
	var response Response[*TokenResponse]
	if err := c.BindJSON(&body); err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := u.userService.VerifyAccessToken(logger, body.AccessToken)
	if err != nil {
		logger.Error(err.Error())
		response.Message = err.Error()
		response.Status = "error"

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Message = "success"
	response.Status = "success"
	response.Data = result

	c.JSON(http.StatusOK, response)
}
