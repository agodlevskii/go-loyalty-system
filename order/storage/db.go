package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-loyalty-system/order"
)

type DBOrder struct {
	db *sql.DB
}

func NewDBOrderStorage(db *sql.DB) (DBOrder, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (number VARCHAR(25), status VARCHAR(15), accrual REAL, uploaded_at timestamptz, "user" VARCHAR(50), UNIQUE(number))`)
	return DBOrder{db: db}, err
}

func (r DBOrder) Add(o order.Order) (order.Order, error) {
	var newOrder order.Order

	req := `INSERT INTO orders(number, status, accrual, uploaded_at, "user") VALUES($1, $2, $3, $4, $5) ON CONFLICT(number) DO NOTHING RETURNING *`
	err := r.db.QueryRow(req, o.Number, o.Status, o.Accrual, o.UploadedAt, o.User).Scan(&newOrder.Number, &newOrder.Status, &newOrder.Accrual, &newOrder.UploadedAt, &newOrder.User)
	if errors.Is(err, sql.ErrNoRows) && newOrder.Number == `` {
		dbOrder, err := r.Find(o.Number)
		if err != nil {
			fmt.Println(`DB ADD`, err)
			return newOrder, err
		}

		if dbOrder.User == o.User {
			return dbOrder, errors.New(order.ErrSameUser)
		}

		return newOrder, errors.New(order.ErrOtherUser)
	}

	return newOrder, err
}

func (r DBOrder) Find(number string) (order.Order, error) {
	var o order.Order
	err := r.db.QueryRow(`SELECT * FROM orders WHERE number = $1`, number).Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.User)
	return o, err
}

func (r DBOrder) FindAll(user string) ([]order.Order, error) {
	os := make([]order.Order, 0)
	rows, err := r.db.Query(`SELECT * FROM orders WHERE "user" = $1`, user)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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
