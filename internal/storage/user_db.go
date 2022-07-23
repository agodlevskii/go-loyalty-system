package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go-loyalty-system/internal/models"
)

type DBUser struct {
	db *sql.DB
}

func NewDBUser(db *sql.DB) (DBUser, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (name VARCHAR(50), password VARCHAR(255), UNIQUE(name))`)
	return DBUser{db: db}, err
}

func (r DBUser) AddUser(user models.User) error {
	_, err := r.db.Exec(`INSERT INTO users (name, password) VALUES ($1, $2)`, user.Login, user.Password)
	return err
}

func (r DBUser) GetUser(name string) (models.User, error) {
	user := models.User{Login: name}
	err := r.db.QueryRow(`SELECT password FROM users WHERE name = &1`, user.Login).Scan(&user.Password)
	return user, err
}
