package adminSvc

import (
	"errors"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var adminKey = []byte(config.GetEnv("ADMIN_KEY"))

type AdminClaims struct {
	Id int `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateAdminToken(id int, username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 365) // 365 days
	claims := AdminClaims{
		Id: id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(adminKey)
	return tokenString, err
}

func VerifyAdminToken(tokenString string) (*AdminClaims, error) {
	claims := &AdminClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return adminKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
