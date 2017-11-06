package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaims struct {
	UserName string `json:"userName"`
	UserId   int    `json:"userId"`
	jwt.StandardClaims
}

type TokenService struct {
	// token有效时长
	Duration  time.Duration
	SignedKey string
}

func GetTokenService() TokenService {
	return TokenService{time.Minute * 30, "test"}
}

func (ts *TokenService) Generate(userId int, userName string) (string, error) {
	expiredAt := time.Now().Add(ts.Duration).Unix()
	claims := UserClaims{userName,
		userId,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ts.SignedKey))
}

func (ts *TokenService) Parse(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.SignedKey), nil
	})
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
