package auth

import (
	"go-loyalty-system/internal/aerror"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashes(reqPwd, dbPwd string) (bool, *aerror.AppError) {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPwd), []byte(reqPwd)); err != nil {
		return false, aerror.NewError(aerror.UserPasswordIncorrect, err)
	}

	return true, aerror.NewEmptyError()
}

func HashPassword(pwd string) (string, *aerror.AppError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return ``, aerror.NewError(aerror.UserPasswordHash, err)
	}
	return string(hash), aerror.NewEmptyError()
}
