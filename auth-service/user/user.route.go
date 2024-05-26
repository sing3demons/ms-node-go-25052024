package user

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sing3demons/auth-service/redis"
	"github.com/sing3demons/auth-service/router"
	"go.mongodb.org/mongo-driver/mongo"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		if s == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		token := strings.TrimPrefix(s, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		public := os.Getenv("PUBLIC_ACCESS_KEY")

		if public == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		publicKey, err := base64.StdEncoding.DecodeString(public)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		rsa, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		t, err := jwt.ParseWithClaims(token, &RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return rsa, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		claims := t.Claims.(*RegisteredClaims)
		c.Set("token", claims)

		c.Next()
	}
}

func Register(r router.MyRouter, client *mongo.Client, redisClient redis.IRedis, logger *slog.Logger) router.MyRouter {
	logger.Info("Register user routes")
	userService := NewUserService(client, redisClient )
	userHandler := NewUserHandler(userService, logger)
	authMiddleware := Authorization()
	v1 := r.Group("/api/v1")

	v1.POST("/auth/register", userHandler.Register)
	v1.POST("/auth/login", userHandler.Login)
	v1.POST("/auth/verify", userHandler.VerifyToken)
	v1.GET("/auth/profile", authMiddleware, userHandler.Profile)
	v1.POST("/auth/refresh", userHandler.Refresh)

	return r

}
