package models

import "time"

const (
	StatusNew        = `NEW`
	StatusInvalid    = `INVALID`
	StatusProcessing = `PROCESSING`
	StatusProcessed  = `PROCESSED`
)

type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
	User       string  `json:"-"`
}

func NewOrderFromAccrual(accrual AccrualOrder, user string) Order {
	return Order{
		Number:     accrual.Order,
		Status:     accrual.Status,
		Accrual:    accrual.Accrual,
		UploadedAt: time.Now().Format(time.RFC3339),
		User:       user,
	}
}

func NewOrderFromOrderAndAccrual(order Order, accrual AccrualOrder) Order {
	return Order{
		Number:     order.Number,
		Status:     accrual.Status,
		Accrual:    accrual.Accrual,
		UploadedAt: order.UploadedAt,
		User:       order.User,
	}
}
