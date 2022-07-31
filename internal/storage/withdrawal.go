package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS withdrawals ("order" VARCHAR(50), sum REAL, processed_at VARCHAR(25), "user" VARCHAR(50), UNIQUE("order"))`)
	if err != nil {
		return DBWithdrawal{}, aerror.NewError(aerror.WithdrawalTableCreate, err)
	}
	return DBWithdrawal{db: db}, nil
}

func (r DBWithdrawal) Add(w models.Withdrawal) *aerror.AppError {
	_, err := r.db.Exec(`INSERT INTO withdrawals ("order", sum, processed_at, "user") VALUES ($1, $2, $3, $4)`, w.Order, w.Sum, w.ProcessedAt, w.User)
	if err != nil {
		return aerror.NewError(aerror.WithdrawalAdd, err)
	}
	return nil
}

func (r DBWithdrawal) Find(order string) (models.Withdrawal, *aerror.AppError) {
	var w models.Withdrawal
	err := r.db.QueryRow(`SELECT * FROM withdrawals WHERE "order" = $1`, order).Scan(&w)
	if err != nil {
		return w, aerror.NewError(aerror.WithdrawalFind, err)
	}
	return w, nil
}

func (r DBWithdrawal) FindAll(user string) ([]models.Withdrawal, *aerror.AppError) {
	ws := make([]models.Withdrawal, 0)
	rows, err := r.db.Query(`SELECT * FROM withdrawals WHERE "user" = $1`, user)
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
