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
	"github.com/sing3demons/auth-service/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	store store.Store
	redis redis.IRedis
}

func NewUserService(client store.Store, redisClient redis.IRedis) UserService {
	return &userService{client, redisClient}
}

const (
	ErrTokenInvalid = "token invalid"
)

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

	if err != nil || !t.Valid {
		return nil, errors.New(ErrTokenInvalid)
	}

	return &TokenResponse{AccessToken: token}, nil
}

func (u *userService) verifyRefreshToken(token string) (jwt.Claims, error) {
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

func (u *userService) Login(ctx context.Context, logger *slog.Logger, body Login) (*TokenResponse, error) {
	logger.Info("userService Login")
	db := u.store.Database("auth").Collection("users")
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

func (u *userService) CreateUser(ctx context.Context, logger *slog.Logger, body User) (User, error) {
	logger.Info("userService Create user")
	db := u.store.Database("auth").Collection("users")

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
	col := u.store.Database("auth")
	db := col.Collection("users")
	var user User

	singleResult := db.FindOne(ctx, bson.M{"id": id})
	if err := singleResult.Decode(&user); err != nil {
		logger.Error(err.Error())
		return User{}, err
	}
	logger.Info("Get user success", "id", user.ID)

	return user, nil
}

func (u *userService) RefreshToken(ctx context.Context, logger *slog.Logger, token string) (*TokenResponse, error) {
	c, err := u.verifyRefreshToken(token)
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

func (u *userService) UpdateUser(ctx context.Context, logger *slog.Logger, body UpdateProfile) (any, error) {
	session, err := u.store.StartSession()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	session.StartTransaction()
	defer session.EndSession(ctx)

	dbProfile := u.store.Database("auth").Collection("profileLanguage")

	profileLanguage := []ProfileLanguage{}
	profileTH := &ProfileLanguage{}
	profileEN := &ProfileLanguage{}
	profileEN.LanguageCode = "en"

	if body.FirstNameTH != "" {
		profileTH.FirstName = body.FirstNameTH
	}

	dbUser := u.store.Database("auth").Collection("users")
	users := Profile{}
	if err := dbUser.FindOne(ctx, bson.M{"id": body.ID}).Decode(&users); err != nil {
		logger.Error(err.Error())
		session.AbortTransaction(ctx)
		return nil, err
	}

	if body.FirstName != "" {
		users.FirstName = body.FirstName
		profileEN.FirstName = body.FirstName
	}

	if body.LastName != "" {
		users.LastName = body.LastName
		profileEN.LastName = body.LastName
	}

	if body.Description != "" {
		users.Description = body.Description
		profileEN.Description = body.Description
	}

	if body.Phone != "" {
		users.Phone = body.Phone
	}

	if body.Address != "" {
		users.Address = body.Address
	}

	users.UpdateDate = time.Now().String()

	updateResult, err := dbUser.UpdateOne(ctx, bson.M{"id": users.ID}, &users, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		logger.Error(err.Error())
		session.AbortTransaction(ctx)
		return nil, err
	}

	if body.ProfileImage != "" {
		profileEN.Attachments = append(profileEN.Attachments, Attachment{
			ID:   uuid.New().String(),
			Name: "profileImage",
			URL:  body.ProfileImage,
			Type: "image",
		})
	}

	profileTH.LanguageCode = "th"
	profileLanguage = append(profileLanguage, *profileTH)
	profileLanguage = append(profileLanguage, *profileEN)
	if len(profileLanguage) > 0 {
		for _, lang := range profileLanguage {
			result, err := dbProfile.UpdateOne(ctx, bson.M{"id": lang.ID}, &lang, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
			if err != nil {
				logger.Error(err.Error())
				session.AbortTransaction(ctx)
				return nil, err
			}
			logger.Info("Update profile language success", "id", result.UpsertedID)
		}
	}

	session.CommitTransaction(ctx)

	return updateResult, nil
}

// func (u *userService) DeleteUser() {}

// func (u *userService) AddRole() {}

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
