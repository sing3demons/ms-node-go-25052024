package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt" // Import the package that contains the hashPassword function
)

type UserService interface {
	CreateUser(ctx context.Context, logger *slog.Logger, body User) (*mongo.InsertOneResult, error)
	GetUser(ctx context.Context, logger *slog.Logger, id string)
	Login(ctx context.Context, logger *slog.Logger, body any)
}

type userService struct {
	*mongo.Client
}

func NewUserService(client *mongo.Client) UserService {
	return &userService{client}
}

func (u *userService) Login(ctx context.Context, logger *slog.Logger, body any) {}

func (u *userService) hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash)
}

func (u *userService) comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *userService) ValidateUser(ctx context.Context, logger *slog.Logger, body User) bool {
	logger.Info("userService Validate user")
	db := u.Client.Database("auth").Collection("users")

	var user User
	if body.Email != "" {
		if err := db.FindOne(ctx, bson.M{"email": body.Email}).Decode(&user); err != nil {
			msg := errors.New("user not found")
			logger.Error(msg.Error())
			return false
		}
	}

	if body.Username != "" {
		if err := db.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user); err != nil {
			msg := errors.New("user not found")
			logger.Error(msg.Error())
			return false
		}

	}

	if !u.comparePassword(body.Password, user.Password) {
		msg := errors.New("password not match")
		logger.Error(msg.Error())
		return false
	}

	return true
}

func (u *userService) CreateUser(ctx context.Context, logger *slog.Logger, body User) (*mongo.InsertOneResult, error) {
	logger.Info("userService Create user")
	db := u.Client.Database("auth").Collection("users")

	if body.Username != "" {
		if err := db.FindOne(ctx, bson.M{"username": body.Username}).Decode(&User{}); err == nil {
			msg := errors.New("username already exists")
			logger.Error(msg.Error())
			return nil, msg
		}

	}

	if body.Email != "" {
		if err := db.FindOne(ctx, bson.M{"email": body.Email}).Decode(&User{}); err == nil {
			msg := errors.New("email already exists")
			logger.Error(msg.Error())
			return nil, msg
		}
	}

	user := User{
		ID:       uuid.New().String(),
		Username: body.Username,
		Password: u.hashPassword(body.Password),
		Email:    body.Email,
		UpdateAt: time.Now(),
		CreateAt: time.Now(),
	}
	r, err := db.InsertOne(ctx, user)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.Info("Create user success", "id", r.InsertedID)
	return r, nil

}

func (u *userService) GetUser(ctx context.Context, logger *slog.Logger, id string) {}

func (u *userService) UpdateUser() {}

func (u *userService) DeleteUser() {}

func (u *userService) AddRole() {}
