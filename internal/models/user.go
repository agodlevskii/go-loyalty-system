package models

const UserKey CtxKey = `user`

type CtxKey string

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
