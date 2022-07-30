package models

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	User      string  `json:"-"`
}

func NewBalance(user string) Balance {
	return Balance{
		Current:   0,
		Withdrawn: 0,
		User:      user,
	}
}
