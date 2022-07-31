package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
)

type BalanceStorage interface {
	Set(balance models.Balance) *aerror.AppError
	Get(user string) (models.Balance, *aerror.AppError)
}

type DBBalance struct {
	db *sql.DB
}

func NewDBBalanceStorage(db *sql.DB) (DBBalance, *aerror.AppError) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS balance ("user" VARCHAR(50), current REAL, withdrawn REAL, UNIQUE("user"))`)
	if err != nil {
		return DBBalance{}, aerror.NewError(aerror.BalanceTableCreate, err)
	}
	return DBBalance{db: db}, aerror.NewEmptyError()
}

func (r DBBalance) Set(b models.Balance) *aerror.AppError {
	_, err := r.db.Exec(`INSERT INTO balance ("user", current, withdrawn) VALUES ($1, $2, $3) ON CONFLICT("user") DO UPDATE SET current = excluded.current, withdrawn = excluded.withdrawn`, b.User, b.Current, b.Withdrawn)
	if err != nil {
		return aerror.NewError(aerror.BalanceSet, err)
	}
	return nil
}

func (r DBBalance) Get(user string) (models.Balance, *aerror.AppError) {
	var b models.Balance
	err := r.db.QueryRow(`SELECT * FROM balance WHERE "user" = $1`, user).Scan(&b.User, &b.Current, &b.Withdrawn)
	if err != nil {
		return b, aerror.NewError(aerror.BalanceGet, err)
	}
	return b, aerror.NewEmptyError()
}
