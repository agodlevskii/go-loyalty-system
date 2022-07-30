package order

import "time"

const (
	StatusNew        = `NEW`
	StatusInvalid    = `INVALID`
	StatusProcessing = `PROCESSING`
	StatusProcessed  = `PROCESSED`
	ErrOtherUser     = `the order was added by another user`
	ErrSameUser      = `the order is already enqueued by the user`
)

type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type AccrualResponse struct {
	StatusCode int
	Accrual    AccrualOrder
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
	User       string  `json:"-"`
}

type Storage interface {
	Add(order Order) (Order, error)
	Update(order Order) (Order, error)
	Find(number string) (Order, error)
	FindAll(user string) ([]Order, error)
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
