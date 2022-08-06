package auth

import (
	"go-loyalty-system/internal/aerror"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashes(reqPwd, dbPwd string) (bool, *aerror.AppError) {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPwd), []byte(reqPwd)); err != nil {
		return false, aerror.NewError(aerror.UserPasswordIncorrect, err)
	}

	return true, nil
}

func HashPassword(pwd string) (string, *aerror.AppError) {
	if pwd == "" {
		return "", aerror.NewError(aerror.UserPasswordIncorrect, nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", aerror.NewError(aerror.UserPasswordHash, err)
	}
	return string(hash), nil
}
