package user

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       string    `json:"id" bson:"id,omitempty"`
	Href     string    `json:"href" bson:"-"`
	Username string    `json:"username,omitempty" bson:"username,omitempty"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"password" bson:"password"`
	Name     string    `json:"name" bson:"name,omitempty"`
	Roles    []string  `json:"roles" bson:"roles,omitempty"`
	CreateAt time.Time `json:"-" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"-" bson:"update_at,omitempty"`
}

type IUser struct {
	ID       string   `json:"id"`
	Href     string   `json:"href"`
	Username string   `json:"username,omitempty"`
	Email    string   `json:"email"`
	Name     string   `json:"name,omitempty"`
	Roles    []string `json:"roles"`
}

type Attachment struct {
	ID      string `json:"id,omitempty" bson:"id,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	URL     string `json:"url,omitempty" bson:"url,omitempty"`
	Type    string `json:"type,omitempty" bson:"type,omitempty"`
	Display struct {
		Type  string   `json:"type,omitempty" bson:"type,omitempty"`
		Value []string `json:"value,omitempty" bson:"value,omitempty"`
	} `json:"display,omitempty" bson:"display,omitempty"`
}

type ProfileLanguage struct {
	ID           string       `json:"id" bson:"id"`
	Type         string       `json:"@Type,omitempty" bson:"@Type,omitempty"`
	Ref          string       `json:"ref,omitempty" bson:"ref,omitempty"`
	Href         string       `json:"href"`
	LanguageCode string       `json:"languageCode" bson:"languageCode"`
	Name         string       `json:"name,omitempty" bson:"name,omitempty"`
	Description  string       `json:"description,omitempty" bson:"description,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty" bson:"attachments,omitempty"`
	CreateDate   string       `json:"createDate,omitempty" bson:"createDate,omitempty"`
	UpdateDate   string       `json:"updateDate,omitempty" bson:"updateDate,omitempty"`
}

type Profile struct {
	ID          string            `json:"id,omitempty" bson:"id,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	Description string            `json:"description,omitempty" bson:"description,omitempty"`
	Phone       string            `json:"phone,omitempty" bson:"phone,omitempty"`
	Address     string            `json:"address,omitempty" bson:"address,omitempty"`
	Languages   []ProfileLanguage `json:"languages,omitempty" bson:"languages,omitempty"`
	CreateDate  string            `json:"createDate,omitempty" bson:"createDate,omitempty"`
	UpdateDate  string            `json:"updateDate,omitempty" bson:"updateDate,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
}

type Login struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}
type Response[T any] struct {
	Message  string `json:"message"`
	Status   string `json:"status"`
	Data     T      `json:"data"`
	Total    uint   `json:"total,omitempty"`
	PageSize uint   `json:"pageSize,omitempty"`
	Page     uint   `json:"page,omitempty"`
}

type RegisteredClaims struct {
	jwt.RegisteredClaims
	UserName string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}
