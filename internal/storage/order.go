package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-loyalty-system/internal/models"
)

type OrderStorage interface {
	Add(order models.Order) (models.Order, error)
	Update(order models.Order) (models.Order, error)
	Find(number string) (models.Order, error)
	FindAll(user string) ([]models.Order, error)
}

type DBOrder struct {
	db *sql.DB
}

func NewDBOrderStorage(db *sql.DB) (DBOrder, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (number VARCHAR(25), status VARCHAR(15), accrual REAL, uploaded_at VARCHAR(25), "user" VARCHAR(50), UNIQUE(number))`)
	return DBOrder{db: db}, err
}

func (r DBOrder) Add(o models.Order) (models.Order, error) {
	var newOrder models.Order

	req := `INSERT INTO orders(number, status, accrual, uploaded_at, "user") VALUES($1, $2, $3, $4, $5) ON CONFLICT(number) DO NOTHING RETURNING *`
	err := r.db.QueryRow(req, o.Number, o.Status, o.Accrual, o.UploadedAt, o.User).Scan(&newOrder.Number, &newOrder.Status, &newOrder.Accrual, &newOrder.UploadedAt, &newOrder.User)
	if errors.Is(err, sql.ErrNoRows) && newOrder.Number == `` {
		dbOrder, err := r.Find(o.Number)
		if err != nil {
			fmt.Println(`DB ADD`, err)
			return newOrder, err
		}

		if dbOrder.User == o.User {
			return dbOrder, errors.New(models.ErrSameUser)
		}

		return newOrder, errors.New(models.ErrOtherUser)
	}

	return newOrder, err
}

func (r DBOrder) Update(o models.Order) (models.Order, error) {
	_, err := r.db.Exec(`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`, o.Status, o.Accrual, o.Number)
	return o, err
}

func (r DBOrder) Find(number string) (models.Order, error) {
	var o models.Order
	err := r.db.QueryRow(`SELECT * FROM orders WHERE number = $1`, number).Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.User)
	return o, err
}

func (r DBOrder) FindAll(user string) ([]models.Order, error) {
	os := make([]models.Order, 0)
	rows, err := r.db.Query(`SELECT * FROM orders WHERE "user" = $1`, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return nil, rows.Err()
	}

	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.User); err != nil {
			fmt.Println(err)
			return nil, err
		}
		os = append(os, o)
	}

	return os, nil
}
