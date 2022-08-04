package storage

import (
	"database/sql"
	"errors"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
)

const (
	OrderAdd         = `INSERT INTO orders(number, status, accrual, uploaded_at, "user") VALUES($1, $2, $3, $4, $5) ON CONFLICT(number) DO NOTHING RETURNING *`
	OrderFind        = `SELECT * FROM orders WHERE number = $1`
	OrderFindAll     = `SELECT * FROM orders WHERE "user" = $1`
	OrderTableCreate = `CREATE TABLE IF NOT EXISTS orders (number VARCHAR(25), status VARCHAR(15), accrual REAL, uploaded_at VARCHAR(25), "user" VARCHAR(50), UNIQUE(number))`
	OrderUpdate      = `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
)

type OrderStorage interface {
	Add(order models.Order) (models.Order, *aerror.AppError)
	Update(order models.Order) (models.Order, *aerror.AppError)
	Find(number string) (models.Order, *aerror.AppError)
	FindAll(user string) ([]models.Order, *aerror.AppError)
}

type DBOrder struct {
	db *sql.DB
}

func NewDBOrderStorage(db *sql.DB) (DBOrder, *aerror.AppError) {
	_, err := db.Exec(OrderTableCreate)
	if err != nil {
		return DBOrder{}, aerror.NewError(aerror.OrderTableCreate, err)
	}
	return DBOrder{db: db}, nil
}

func (r DBOrder) Add(o models.Order) (models.Order, *aerror.AppError) {
	var newOrder models.Order
	err := r.db.QueryRow(OrderAdd, o.Number, o.Status, o.Accrual, o.UploadedAt, o.User).Scan(&newOrder.Number, &newOrder.Status, &newOrder.Accrual, &newOrder.UploadedAt, &newOrder.User)
	if err != nil {
		return handleOrderAddFailure(r, o, err)
	}
	return newOrder, nil
}

func (r DBOrder) Update(o models.Order) (models.Order, *aerror.AppError) {
	_, err := r.db.Exec(OrderUpdate, o.Status, o.Accrual, o.Number)
	if err != nil {
		return models.Order{}, aerror.NewError(aerror.OrderUpdate, err)
	}
	return o, nil
}

func (r DBOrder) Find(number string) (models.Order, *aerror.AppError) {
	var o models.Order
	err := r.db.QueryRow(OrderFind, number).Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.User)
	if err != nil {
		return o, aerror.NewError(aerror.OrderFind, err)
	}
	return o, nil
}

func (r DBOrder) FindAll(user string) ([]models.Order, *aerror.AppError) {
	os := make([]models.Order, 0)
	rows, err := r.db.Query(OrderFindAll, user)
	if err != nil {
		return nil, aerror.NewError(aerror.OrderFindAll, err)
	}
	if rows.Err() != nil {
		return nil, aerror.NewError(aerror.OrderFindAll, rows.Err())
	}

	for rows.Next() {
		var o models.Order
		if err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.User); err != nil {
			return nil, aerror.NewError(aerror.OrderFindAll, err)
		}
		os = append(os, o)
	}

	return os, nil
}

func handleOrderAddFailure(r DBOrder, o models.Order, err error) (models.Order, *aerror.AppError) {
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		if dbOrder, err := r.Find(o.Number); err == nil {
			if dbOrder.User == o.User {
				return dbOrder, aerror.NewError(aerror.OrderExistsSameUser, nil)
			}

			return models.Order{}, aerror.NewError(aerror.OrderExistsOtherUser, nil)
		}
	}

	return models.Order{}, aerror.NewError(aerror.OrderAdd, err)
}
