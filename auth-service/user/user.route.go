package user

import (
	"log/slog"

	"github.com/sing3demons/auth-service/router"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(r router.MyRouter, client *mongo.Client, logger *slog.Logger) router.MyRouter {
	logger.Info("Register user routes")
	userService := NewUserService(client)
	userHandler := NewUserHandler(userService, logger)
	r.POST("/register", userHandler.Register)

	return r

}
