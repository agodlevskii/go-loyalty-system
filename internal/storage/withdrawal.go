package storage

import (
	"database/sql"
	"go-loyalty-system/internal/models"
)

type WithdrawalStorage interface {
	Add(withdrawal models.Withdrawal) error
	Find(order string) (models.Withdrawal, error)
	FindAll(user string) ([]models.Withdrawal, error)
}

type DBWithdrawal struct {
	db *sql.DB
}

func NewDBWithdrawalStorage(db *sql.DB) (DBWithdrawal, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS withdrawals ("order" VARCHAR(50), sum REAL, processed_at TIMETZ, "user" VARCHAR(50), UNIQUE("order"))`)
	return DBWithdrawal{db: db}, err
}

func (r DBWithdrawal) Add(w models.Withdrawal) error {
	_, err := r.db.Exec(`INSERT INTO withdrawals ("order", sum, processed_at, "user") VALUES ($1, $2, $3, $4)`, w.Order, w.Sum, w.ProcessedAt, w.User)
	return err
}

func (r DBWithdrawal) Find(order string) (models.Withdrawal, error) {
	var w models.Withdrawal
	err := r.db.QueryRow(`SELECT * FROM withdrawals WHERE "order" = $1`, order).Scan(&w)
	return w, err
}

func (r DBWithdrawal) FindAll(user string) ([]models.Withdrawal, error) {
	ws := make([]models.Withdrawal, 0)
	rows, err := r.db.Query(`SELECT * FROM withdrawals WHERE user = $1`, user)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var w models.Withdrawal
		err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt, &w.User)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}

	return ws, nil
}
