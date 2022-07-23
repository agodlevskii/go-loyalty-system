package user

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Storage interface {
	AddUser(user User) error
	GetUser(name string) (User, error)
}
