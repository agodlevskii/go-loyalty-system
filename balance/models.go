package balance

import "time"

type Account struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	User      string  `json:"-"`
}

type AccountStorage interface {
	Set(balance Account) error
	Get(uName string) (Account, error)
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
	User        string  `json:"-"`
}

type WithdrawalStorage interface {
	Add(withdrawal Withdrawal) error
	Find(order string) (Withdrawal, error)
	FindAll(user string) ([]Withdrawal, error)
}

func NewAccount(user string) Account {
	return Account{
		Current:   0,
		Withdrawn: 0,
		User:      user,
	}
}

func NewWithdrawalFromRequest(wr Withdrawal, user string) Withdrawal {
	return Withdrawal{
		Order:       wr.Order,
		Sum:         wr.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
		User:        user,
	}
}
