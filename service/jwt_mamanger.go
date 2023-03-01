package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-usermgmt-grpc/db/models"
	"time"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

func NewJWTManager(secretkey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretkey,
		tokenDuration: tokenDuration,
	}
}

func (manager *JWTManager) Generate(user *models.Users) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
		},
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("nexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}
	return claims, nil
}
