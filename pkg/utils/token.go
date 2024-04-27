package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-finance-tracker/internal/models"
	"os"
	"time"
)

var (
	ErrEmptyToken            = errors.New("token is empty")
	ErrEmptyTokenClaims      = errors.New("token claims are nil")
	ErrInvalidToken          = errors.New("token is invalid")
	ErrInvalidTokenSignature = errors.New("signature is invalid")
	ErrInvalidParsedToken    = errors.New("parsed token is invalid")
)

func CreateToken(id, username string) (tokenString string, err error) {
	claims := &models.Claims{
		Id:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(token string) (string, string, error) {
	if token == "" {
		return "", "", ErrEmptyToken
	}
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return "", "", ErrInvalidTokenSignature
		}
		return "", "", ErrInvalidToken
	}
	if !parsedToken.Valid {
		return "", "", ErrInvalidParsedToken
	}
	if claims == nil {
		return "", "", ErrEmptyTokenClaims
	}
	return claims.Id, claims.Username, nil
}
