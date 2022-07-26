package storage

import (
	"database/sql"
	"go-loyalty-system/balance"
)

type DBBalance struct {
	db *sql.DB
}

func NewDBBalanceStorage(db *sql.DB) (DBBalance, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS balance ("user" VARCHAR(50), current REAL, withdrawn REAL, UNIQUE("user"))`)
	return DBBalance{db: db}, err
}

func (r DBBalance) Set(b balance.Balance) error {
	_, err := r.db.Exec(`INSERT INTO balance ("user", current, withdrawn) VALUES ($1, $2, $3) ON CONFLICT(current) DO UPDATE SET current = $2, withdrawn = $3`, b.User, b.Current, b.Withdrawn)
	return err
}

func (r DBBalance) Get(user string) (balance.Balance, error) {
	var b balance.Balance
	err := r.db.QueryRow(`SELECT * FROM balance WHERE user = $1`, user).Scan(&b.User, &b.Current, &b.Withdrawn)
	return b, err
}
