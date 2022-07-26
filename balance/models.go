package balance

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	User      string
}

type Storage interface {
	Set(balance Balance) error
	Get(uName string) (Balance, error)
}
