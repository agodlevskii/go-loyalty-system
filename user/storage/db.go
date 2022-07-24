package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go-loyalty-system/user"
)

type DBUser struct {
	db *sql.DB
}

func NewDBUserStorage(db *sql.DB) (DBUser, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (name VARCHAR(50), password VARCHAR(255), UNIQUE(name))`)
	return DBUser{db: db}, err
}

func (r DBUser) Add(user user.User) error {
	var name string
	return r.db.QueryRow(`INSERT INTO users (name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING name`, user.Login, user.Password).Scan(&name)
}

func (r DBUser) Find(name string) (user.User, error) {
	dbUser := user.User{Login: name}
	err := r.db.QueryRow(`SELECT password FROM users WHERE name = $1`, dbUser.Login).Scan(&dbUser.Password)
	return dbUser, err
}
