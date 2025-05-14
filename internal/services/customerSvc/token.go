package customerSvc

import (
	"errors"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var customerKey = []byte(config.GetEnv("CUSTOMER_KEY"))

type CustomerClaims struct {
	Id int `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateCustomerToken(id int, username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 365)
	claims := CustomerClaims{
		Id: id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(customerKey)	
	return tokenString, err
}

func VerifyCustomerToken(tokenString string) (*CustomerClaims, error) {
	claims := &CustomerClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return customerKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}