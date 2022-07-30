package internal

import (
	"database/sql"
	"go-loyalty-system/balance"
	balanceStorage "go-loyalty-system/balance/storage"
	"go-loyalty-system/order"
	orderStorage "go-loyalty-system/order/storage"
	"go-loyalty-system/user"
	userStorage "go-loyalty-system/user/storage"
)

type Repo struct {
	User       user.Storage
	Account    balance.AccountStorage
	Order      order.Storage
	Withdrawal balance.WithdrawalStorage
}

func NewDBRepo(url, driver string) (Repo, error) {
	if driver == "" {
		driver = `pgx`
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		return Repo{}, err
	}

	bs, err := balanceStorage.NewDBAccountStorage(db)
	if err != nil {
		return Repo{}, err
	}

	os, err := orderStorage.NewDBOrderStorage(db)
	if err != nil {
		return Repo{}, err
	}

	us, err := userStorage.NewDBUserStorage(db)
	if err != nil {
		return Repo{}, err
	}

	ws, err := balanceStorage.NewDBWithdrawalStorage(db)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		Account:    bs,
		Order:      os,
		User:       us,
		Withdrawal: ws,
	}, nil
}
