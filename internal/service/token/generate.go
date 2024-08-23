package token

import (
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/dgrijalva/jwt-go"
)

func (s *service) Generate(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Name: info.Name,
		Role: info.Role,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}
