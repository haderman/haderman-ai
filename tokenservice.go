package main

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id    string `json:"id"`
	First string `json:"first"`
	jwt.RegisteredClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}

func ParseRefreshToken(refreshToken string) *jwt.RegisteredClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.RegisteredClaims)
}

func NewHadermanBackendServerAccessToken() (string, error) {
	claims := UserClaims{
		Id:    os.Getenv("USER_ID"),
		First: os.Getenv("USER_NAME"),
	}

	signedAccessToken, err := NewAccessToken(claims)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}
