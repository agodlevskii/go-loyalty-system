package withdrawal

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type Storage interface {
	AddWithdrawal(withdrawal Withdrawal) error
	GetWithdrawal(order string) (Withdrawal, error)
	GetWithdrawals(user string) ([]Withdrawal, error)
}
