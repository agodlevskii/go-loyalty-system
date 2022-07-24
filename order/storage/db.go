package storage

import (
	"database/sql"
	"go-loyalty-system/order"
)

type DBOrder struct {
	db *sql.DB
}

func NewDBOrderStorage(db *sql.DB) (DBOrder, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (number VARCHAR(25), status VARCHAR(15), accrual REAL, uploaded_at TIMETZ, "user" VARCHAR(50), UNIQUE(number))`)
	return DBOrder{db: db}, err
}

func (r DBOrder) Add(o order.Order) error {
	_, err := r.db.Exec(`INSERT INTO orders(number, status, accrual, uploaded_at, user) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING`, o.Number, o.Status, o.Accrual, o.UploadedAt, o.User)
	return err
}

func (r DBOrder) Find(number string) (order.Order, error) {
	var o order.Order
	err := r.db.QueryRow(`SELECT * FROM orders WHERE number = $1`, number).Scan(&o.Number, &o.Status, &o.Accrual, o.UploadedAt, o.User)
	return o, err
}

func (r DBOrder) FindAll(user string) ([]order.Order, error) {
	os := make([]order.Order, 0)
	rows, err := r.db.Query(`SELECT * FROM orders WHERE user = $1`, user)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o order.Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, o.UploadedAt, o.User)
		if err != nil {
			return nil, err
		}
		os = append(os, o)
	}

	return os, nil
}
