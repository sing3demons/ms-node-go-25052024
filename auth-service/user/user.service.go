package user

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sing3demons/auth-service/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, logger *slog.Logger, body User) (User, error)
	GetUser(ctx context.Context, logger *slog.Logger, id string) (User, error)
	Login(ctx context.Context, logger *slog.Logger, body Login) (*TokenResponse, error)
	RefreshToken(ctx context.Context, logger *slog.Logger, token string) (*TokenResponse, error)
	VerifyAccessToken(logger *slog.Logger, token string) (*TokenResponse, error)
}

type userService struct {
	*mongo.Client
	redis redis.IRedis
}

func NewUserService(client *mongo.Client, redisClient redis.IRedis) UserService {
	return &userService{client, redisClient}
}

const (
	ErrTokenInvalid = "token invalid"
)

func (u *userService) generateAccessToken(user User) (string, error) {
	private := os.Getenv("PRIVATE_ACCESS_KEY")
	if private == "" {
		return "", errors.New("private key not found")
	}
	privateKey, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return "", err
	}
	rsa, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}

	claims := &RegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    os.Getenv("ISSUER"),
			Audience:  jwt.ClaimStrings{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	if user.Email != "" {
		claims.Email = user.Email
	}

	if user.Username != "" {
		claims.UserName = user.Username
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
}

func (u *userService) VerifyAccessToken(logger *slog.Logger, token string) (*TokenResponse, error) {
	public := os.Getenv("PUBLIC_ACCESS_KEY")

	if public == "" {
		return nil, errors.New("public key not found")
	}

	publicKey, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}
	rsa, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	t, err := jwt.ParseWithClaims(token, &RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return rsa, nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New(ErrTokenInvalid)
	}

	return &TokenResponse{AccessToken: token}, nil
}

func (u *userService) VerifyRefreshToken(token string) (jwt.Claims, error) {
	public := os.Getenv("PUBLIC_REFRESH_KEY")
	if public == "" {
		return nil, errors.New("public key not found")
	}

	publicKey, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}
	rsa, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	t, err := jwt.ParseWithClaims(token, &RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return rsa, nil
	})

	if err != nil {
		return nil, err
	}

	return t.Claims, nil
}

func (u *userService) generateRefreshToken(user User) (string, error) {
	private := os.Getenv("PRIVATE_REFRESH_KEY")
	if private == "" {
		return "", errors.New("private key not found")
	}

	privateKey, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return "", err
	}
	rsa, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}

	claims := &RegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    os.Getenv("ISSUER"),
			Audience:  jwt.ClaimStrings{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	if user.Email != "" {
		claims.Email = user.Email
	}

	if user.Username != "" {
		claims.UserName = user.Username
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
}

func (u *userService) Login(ctx context.Context, logger *slog.Logger, body Login) (*TokenResponse, error) {
	logger.Info("userService Login")
	db := u.Client.Database("auth").Collection("users")
	var user User
	if body.Email != "" {
		if err := db.FindOne(ctx, bson.M{"email": body.Email}).Decode(&user); err != nil {
			msg := errors.New("user not found")
			logger.Error(msg.Error())
			return nil, msg
		}
	}

	if body.Username != "" {
		if err := db.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user); err != nil {
			msg := errors.New("user not found")
			logger.Error(msg.Error())
			return nil, msg
		}
	}

	if err := u.comparePassword(user.Password, body.Password); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var token TokenResponse

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("generate access token failed")
	}
	token.AccessToken = accessToken

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("generate refresh token failed")
	}
	token.RefreshToken = refreshToken

	if err := u.redis.SetEx(ctx, refreshToken, "true", time.Minute*60); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &token, nil

}

func (u *userService) hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (u *userService) comparePassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func (u *userService) CreateUser(ctx context.Context, logger *slog.Logger, body User) (User, error) {
	logger.Info("userService Create user")
	db := u.Client.Database("auth").Collection("users")

	if body.Username != "" {
		if err := db.FindOne(ctx, bson.M{"username": body.Username}).Decode(&User{}); err == nil {
			msg := errors.New("username already exists")
			logger.Error(msg.Error())
			return User{}, msg
		}

	}

	if body.Email != "" {
		if err := db.FindOne(ctx, bson.M{"email": body.Email}).Decode(&User{}); err == nil {
			msg := errors.New("email already exists")
			logger.Error(msg.Error())
			return User{}, msg
		}
	}

	hash, err := u.hashPassword(body.Password)
	if err != nil {
		logger.Error(err.Error())
		return User{}, err
	}

	user := User{
		ID:       uuid.New().String(),
		Username: body.Username,
		Password: hash,
		Email:    body.Email,
		Roles: []string{
			"user",
		},
		UpdateAt: time.Now(),
		CreateAt: time.Now(),
	}

	r, err := db.InsertOne(ctx, user)
	if err != nil {
		logger.Error(err.Error())
		return User{}, err
	}
	logger.Info("Create user success", "id", r.InsertedID)

	return user, nil

}

func (u *userService) GetUser(ctx context.Context, logger *slog.Logger, id string) (User, error) {
	logger.Info("userService Get user", "id", id)
	db := u.Client.Database("auth").Collection("users")
	var user User
	if err := db.FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
		logger.Error(err.Error())
		return User{}, err
	}
	logger.Info("Get user success", "id", user.ID)

	return user, nil
}

func (u *userService) RefreshToken(ctx context.Context, logger *slog.Logger, token string) (*TokenResponse, error) {
	c, err := u.VerifyRefreshToken(token)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	customClaims, ok := c.(*RegisteredClaims)
	if !ok {
		logger.Error(ErrTokenInvalid)
		return nil, errors.New(ErrTokenInvalid)
	}

	intCmdVal, err := u.redis.Exists(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if intCmdVal == 0 {
		logger.Error("refresh token not found")
		return nil, errors.New("refresh token not found")
	}

	if err := u.redis.Del(ctx, token); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	user, err := u.GetUser(ctx, logger, customClaims.Subject)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var response TokenResponse

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("generate access token failed")
	}

	response.AccessToken = accessToken

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("generate refresh token failed")
	}

	if err := u.redis.SetEx(ctx, refreshToken, "true", time.Minute*60); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	response.RefreshToken = refreshToken

	return &response, nil
}

// func (u *userService) UpdateUser() {}

// func (u *userService) DeleteUser() {}

// func (u *userService) AddRole() {}
