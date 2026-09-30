package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, id int, rememberMe bool) (string, error) {
	var expirationDate int64
	if !rememberMe {
		expirationDate = time.Now().Add(24 * time.Hour).Unix()
	} else {
		expirationDate = time.Now().Add(30 * (24 * time.Hour)).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": expirationDate,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
