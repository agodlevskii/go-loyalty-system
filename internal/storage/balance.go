package storage

import (
	"database/sql"
	"go-loyalty-system/internal/models"
)

type BalanceStorage interface {
	Set(balance models.Balance) error
	Get(uName string) (models.Balance, error)
}

type DBBalance struct {
	db *sql.DB
}

func NewDBBalanceStorage(db *sql.DB) (DBBalance, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS balance ("user" VARCHAR(50), current REAL, withdrawn REAL, UNIQUE("user"))`)
	return DBBalance{db: db}, err
}

func (r DBBalance) Set(b models.Balance) error {
	_, err := r.db.Exec(`INSERT INTO balance ("user", current, withdrawn) VALUES ($1, $2, $3) ON CONFLICT("user") DO UPDATE SET current = $2, withdrawn = $3`, b.User, b.Current, b.Withdrawn)
	return err
}

func (r DBBalance) Get(user string) (models.Balance, error) {
	var b models.Balance
	err := r.db.QueryRow(`SELECT * FROM balance WHERE user = $1`, user).Scan(&b.User, &b.Current, &b.Withdrawn)
	return b, err
}
