package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
)

const (
	UserTableCreate = `CREATE TABLE IF NOT EXISTS users (name VARCHAR(50), password VARCHAR(255), UNIQUE(name))`
	UserAdd         = `INSERT INTO users (name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING name`
	UserFind        = `SELECT password FROM users WHERE name = $1`
)

type UserStorage interface {
	Add(user models.User) *aerror.AppError
	Find(name string) (models.User, *aerror.AppError)
}

type DBUser struct {
	db *sql.DB
}

func NewDBUserStorage(db *sql.DB) (DBUser, *aerror.AppError) {
	_, err := db.Exec(UserTableCreate)
	if err != nil {
		return DBUser{}, aerror.NewError(aerror.UserTableCreate, err)
	}
	return DBUser{db: db}, nil
}

func (r DBUser) Add(user models.User) *aerror.AppError {
	var name string
	if err := r.db.QueryRow(UserAdd, user.Login, user.Password).Scan(&name); err != nil {
		return aerror.NewError(aerror.UserAdd, err)
	}
	return nil
}

func (r DBUser) Find(name string) (models.User, *aerror.AppError) {
	dbUser := models.User{Login: name}
	if err := r.db.QueryRow(UserFind, dbUser.Login).Scan(&dbUser.Password); err != nil {
		return models.User{}, aerror.NewError(aerror.UserFind, err)
	}
	return dbUser, nil
}
