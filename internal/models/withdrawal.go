package models

import "time"

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
	User        string  `json:"-"`
}

func NewWithdrawalFromWithdrawal(wr Withdrawal, user string) Withdrawal {
	return Withdrawal{
		Order:       wr.Order,
		Sum:         wr.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
		User:        user,
	}
}
