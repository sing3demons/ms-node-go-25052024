package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/auth-service/logger"
	"github.com/sing3demons/auth-service/mlog"
	"github.com/sing3demons/auth-service/redis"
	"github.com/sing3demons/auth-service/router"
	"github.com/sing3demons/auth-service/store"
	"github.com/sing3demons/auth-service/user"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env.dev"); err != nil {
			panic(err)
		}

	}
}

func main() {
	port := os.Getenv("PORT")

	logger := logger.New()
	logger.Info("Starting the application...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db := store.NewStore(ctx)
	defer db.Disconnect(ctx)

	redisClient := redis.New()
	defer redisClient.Close()

	r := router.New()
	r.Use(mlog.Middleware(logger))
	r.GET("/healthz", func(c *gin.Context) {
		if err := db.Ping(ctx, readpref.Primary()); err != nil {
			logger.Error(err.Error())
			c.JSON(500, "MongoDB is not available")
			return
		}

		_, err := redisClient.Ping(ctx)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(500, "Internal Server Error")
			return
		}
		c.JSON(200, "OK")
	})

	user.Register(r, db, redisClient, logger)

	r.StartHTTP(port)
}
