package services

import (
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/utils/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaims struct {
	ID       uint32 `json:"id"`
	Nickname string `json:"nickname"`
	jwt.StandardClaims
}

func GenerateUserToken(cfg config.Config, user models.User) (string, error) {
	secret := []byte(cfg.Secret)
	claims := UserClaims{
		ID:       user.Id,
		Nickname: user.Nickname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
