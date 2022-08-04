package auth

import (
	"github.com/dgrijalva/jwt-go"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"strings"
	"time"
)

type Claims struct {
	User string
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key0")

func GetTokenFromUser(user models.User) (string, *aerror.AppError) {
	exp := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		User: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	if err != nil {
		return ``, aerror.NewError(aerror.UserTokenGeneration, err)
	}
	return token, nil
}

func GetUserFromToken(tokenStr string) (models.User, *aerror.AppError) {
	if tokenStr == `` {
		return models.User{}, aerror.NewError(aerror.UserTokenIncorrect, nil)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFn)
	if err == nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return models.User{Login: claims.User}, nil
		}
	}
	return models.User{}, aerror.NewError(aerror.UserTokenIncorrect, err)
}

func IsTokenValid(tokenStr string) (bool, *aerror.AppError) {
	if tokenStr == `` {
		return false, aerror.NewError(aerror.UserTokenInvalid, nil)
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFn)
	if err != nil {
		return false, aerror.NewError(aerror.UserTokenInvalid, nil)
	}
	return token.Valid, nil
}

func keyFn(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

func GetBearer(token string) string {
	if token == `` {
		return ``
	}
	return `Bearer ` + token
}

func GetTokenFromBearer(bearer string) (string, *aerror.AppError) {
	res := strings.Split(bearer, `Bearer `)
	if len(res) != 2 || res[0] != `` {
		return ``, aerror.NewError(aerror.UserTokenIncorrect, nil)
	}

	return res[1], nil
}
