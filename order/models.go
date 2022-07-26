package order

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    string `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
	User       string
}

type Storage interface {
	Add(order Order) error
	Find(number string) (Order, error)
	FindAll(user string) ([]Order, error)
}
