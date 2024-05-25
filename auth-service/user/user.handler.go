package user

import (
	"context"
	"log/slog"
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
}

func NewUserHandler(userService UserService, logger *slog.Logger) UserHandler {
	return &userHandler{userService, logger}
}

func (u *userHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	logger := mlog.L(ctx)
	logger.Info("Create user")

	var body User
	if err := c.BindJSON(&body); err != nil {
		logger.Error(err.Error())
		c.JSON(400, "Bad Request")
		return
	}

	result, err := u.userService.CreateUser(ctx, logger, body)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(500, "Internal Server Error")
		return
	}

	c.JSON(200, gin.H{
		"message": "Create user success",
		"status":  "success",
		"result":  result,
	})
}
