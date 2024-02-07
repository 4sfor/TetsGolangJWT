package main

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// секретный ключ
var secret = []byte("TestTask")

// создание токена доступа
func TokenAccess() (string, error) {
	claimsA := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 3).Unix()}
	tokenA := jwt.NewWithClaims(jwt.SigningMethodHS512, claimsA)

	tokenStringA, err := tokenA.SignedString(secret)
	if err != nil {

	}
	return tokenStringA, err
}

// создание токена обновления
func TokenRefresh() (string, error) {
	claimsR := jwt.MapClaims{"exp": time.Now().Add(time.Hour * 24).Unix(), "type": "refresh"}
	tokekR := jwt.NewWithClaims(jwt.SigningMethodHS512, claimsR)

	tokenStringR, err := tokekR.SignedString(secret)
	if err != nil {

	}
	return tokenStringR, err
}
