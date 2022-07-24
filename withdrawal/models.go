package withdrawal

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
	User        string
}

type Storage interface {
	Add(withdrawal Withdrawal) error
	Find(order string) (Withdrawal, error)
	FindAll(user string) ([]Withdrawal, error)
}
