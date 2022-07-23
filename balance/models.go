package balance

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Storage interface {
	SetBalance(balance Balance) error
	GetBalance(uName string) (Balance, error)
}
