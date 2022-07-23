package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func compareHashes(reqPwd, dbPwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(dbPwd), []byte(reqPwd))
	return err == nil, err
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}
