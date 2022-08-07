package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
)

const (
	WithdrawalTableCreate = `CREATE TABLE IF NOT EXISTS withdrawals ("order" VARCHAR(50), sum REAL, processed_at VARCHAR(25), "user" VARCHAR(50), UNIQUE("order"))`
	WithdrawalAdd         = `INSERT INTO withdrawals ("order", sum, processed_at, "user") VALUES ($1, $2, $3, $4)`
	WithdrawalFind        = `SELECT * FROM withdrawals WHERE "order" = $1`
	WithdrawalFindAll     = `SELECT * FROM withdrawals WHERE "user" = $1`
)

type WithdrawalStorage interface {
	Add(withdrawal models.Withdrawal) *aerror.AppError
	Find(order string) (models.Withdrawal, *aerror.AppError)
	FindAll(user string) ([]models.Withdrawal, *aerror.AppError)
}

type DBWithdrawal struct {
	db *sql.DB
}

func NewDBWithdrawalStorage(db *sql.DB) (DBWithdrawal, *aerror.AppError) {
	_, err := db.Exec(WithdrawalTableCreate)
	if err != nil {
		return DBWithdrawal{}, aerror.NewError(aerror.WithdrawalTableCreate, err)
	}
	return DBWithdrawal{db: db}, nil
}

func (r DBWithdrawal) Add(w models.Withdrawal) *aerror.AppError {
	var b models.Balance
	tx, err := r.db.Begin()
	if err != nil {
		return aerror.NewError(aerror.WithdrawalAdd, err)
	}
	if err = tx.QueryRow(BalanceGet, w.User).Scan(&b.User, &b.Current, &b.Withdrawn); err != nil {
		return aerror.NewError(aerror.BalanceGet, err)
	}
	if b.User == `` {
		b = models.NewBalance(w.User)
	}
	if b.Current < w.Sum {
		return aerror.NewError(aerror.BalanceInsufficient, nil)
	}

	if _, err = tx.Exec(WithdrawalAdd, w.Order, w.Sum, w.ProcessedAt, w.User); err == nil {
		_, err = tx.Exec(BalanceSet, b.User, b.Current-w.Sum, b.Withdrawn+w.Sum)
	}

	if err != nil {
		if err = tx.Rollback(); err != nil {
			return aerror.NewError(aerror.SystemRollback, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return aerror.NewError(aerror.SystemCommit, err)
	}

	return nil
}

func (r DBWithdrawal) Find(order string) (models.Withdrawal, *aerror.AppError) {
	var w models.Withdrawal
	err := r.db.QueryRow(WithdrawalFind, order).Scan(&w)
	if err != nil {
		return w, aerror.NewError(aerror.WithdrawalFind, err)
	}
	return w, nil
}

func (r DBWithdrawal) FindAll(user string) ([]models.Withdrawal, *aerror.AppError) {
	ws := make([]models.Withdrawal, 0)
	rows, err := r.db.Query(WithdrawalFindAll, user)
	if err != nil {
		return nil, aerror.NewError(aerror.WithdrawalFindAll, err)
	}
	if rows.Err() != nil {
		return nil, aerror.NewError(aerror.WithdrawalFindAll, rows.Err())
	}

	for rows.Next() {
		var w models.Withdrawal
		err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt, &w.User)
		if err != nil {
			return nil, aerror.NewError(aerror.WithdrawalFindAll, err)
		}
		ws = append(ws, w)
	}

	return ws, nil
}
