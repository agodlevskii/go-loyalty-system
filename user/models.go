package user

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Storage interface {
	Add(user User) error
	Find(name string) (User, error)
}
