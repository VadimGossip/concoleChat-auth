package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Name string `json:"username"`
	Role string `json:"role"`
}
