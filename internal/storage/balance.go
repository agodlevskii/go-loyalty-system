package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
)

const (
	BalanceTableCreate = `CREATE TABLE IF NOT EXISTS balance ("user" VARCHAR(50), current REAL, withdrawn REAL, UNIQUE("user"))`
	BalanceSet         = `INSERT INTO balance ("user", current, withdrawn) VALUES ($1, $2, $3) ON CONFLICT("user") DO UPDATE SET current = excluded.current, withdrawn = excluded.withdrawn`
	BalanceGet         = `SELECT * FROM balance WHERE "user" = $1`
)

type BalanceStorage interface {
	Set(balance models.Balance) *aerror.AppError
	Get(user string) (models.Balance, *aerror.AppError)
}

type DBBalance struct {
	db *sql.DB
}

func NewDBBalanceStorage(db *sql.DB) (DBBalance, *aerror.AppError) {
	_, err := db.Exec(BalanceTableCreate)
	if err != nil {
		return DBBalance{}, aerror.NewError(aerror.BalanceTableCreate, err)
	}
	return DBBalance{db: db}, nil
}

func (r DBBalance) Set(b models.Balance) *aerror.AppError {
	_, err := r.db.Exec(BalanceSet, b.User, b.Current, b.Withdrawn)
	if err != nil {
		return aerror.NewError(aerror.BalanceSet, err)
	}
	return nil
}

func (r DBBalance) Get(user string) (models.Balance, *aerror.AppError) {
	var b models.Balance
	err := r.db.QueryRow(BalanceGet, user).Scan(&b.User, &b.Current, &b.Withdrawn)
	if err != nil {
		return b, aerror.NewError(aerror.BalanceGet, err)
	}
	return b, nil
}
