package auth

import (
	"github.com/dgrijalva/jwt-go"
	"go-loyalty-system/user"
	"time"
)

type Claims struct {
	User string
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key0")

func getToken(usr user.User) (string, error) {
	exp := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		User: usr.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func isTokenValid(tokenStr string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFn)
	return token.Valid, err
}

func keyFn(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}
