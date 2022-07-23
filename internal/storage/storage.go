package storage

import (
	"database/sql"
	"go-loyalty-system/internal/models"
)

type UserStorage interface {
	AddUser(user models.User) error
	GetUser(name string) (models.User, error)
}

type OrderStorage interface {
	AddOrder(order models.Order) error
	GetOrder(number string) (models.Order, error)
	GetOrders(user string) ([]models.Order, error)
}

type BalanceStorage interface {
	SetBalance(balance models.Balance) error
	GetBalance(uName string) (models.Balance, error)
}

type WithdrawalStorage interface {
	AddWithdrawal(withdrawal models.Withdrawal) error
	GetWithdrawal(order string) (models.Withdrawal, error)
	GetWithdrawals(user string) ([]models.Withdrawal, error)
}

type Repo struct {
	User       UserStorage
	Balance    BalanceStorage
	Order      OrderStorage
	Withdrawal WithdrawalStorage
}

func NewDBRepo(url, driver string) (Repo, error) {
	if driver == "" {
		driver = `pgx`
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		return Repo{}, err
	}

	user, err := NewDBUser(db)
	if err != nil {
		return Repo{}, err
	}

	return Repo{User: user}, nil
}
