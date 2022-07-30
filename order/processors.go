package order

import "time"

func getOrderFromAccrual(accrual AccrualOrder, user string) Order {
	return Order{
		Number:     accrual.Order,
		Status:     accrual.Status,
		Accrual:    accrual.Accrual,
		UploadedAt: time.Now(),
		User:       user,
	}
}

func combineOrderAndAccrual(order Order, accrual AccrualOrder) Order {
	return Order{
		Number:     order.Number,
		Status:     accrual.Status,
		Accrual:    accrual.Accrual,
		UploadedAt: order.UploadedAt,
		User:       order.User,
	}
}
