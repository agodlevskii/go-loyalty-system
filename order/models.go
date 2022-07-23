package order

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    string `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

type Storage interface {
	AddOrder(order Order) error
	GetOrder(number string) (Order, error)
	GetOrders(user string) ([]Order, error)
}
