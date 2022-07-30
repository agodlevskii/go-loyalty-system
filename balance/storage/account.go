package storage

import (
	"database/sql"
	"go-loyalty-system/balance"
)

type DBAccount struct {
	db *sql.DB
}

func NewDBAccountStorage(db *sql.DB) (DBAccount, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS account ("user" VARCHAR(50), current REAL, withdrawn REAL, UNIQUE("user"))`)
	return DBAccount{db: db}, err
}

func (r DBAccount) Set(a balance.Account) error {
	_, err := r.db.Exec(`INSERT INTO account ("user", current, withdrawn) VALUES ($1, $2, $3) ON CONFLICT("user") DO UPDATE SET current = $2, withdrawn = $3`, a.User, a.Current, a.Withdrawn)
	return err
}

func (r DBAccount) Get(user string) (balance.Account, error) {
	var a balance.Account
	err := r.db.QueryRow(`SELECT * FROM account WHERE user = $1`, user).Scan(&a.User, &a.Current, &a.Withdrawn)
	return a, err
}
