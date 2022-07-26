package internal

import (
	"database/sql"
	"go-loyalty-system/balance"
	balanceStorage "go-loyalty-system/balance/storage"
	"go-loyalty-system/order"
	orderStorage "go-loyalty-system/order/storage"
	"go-loyalty-system/user"
	userStorage "go-loyalty-system/user/storage"
	"go-loyalty-system/withdrawal"
	withdrawalStorage "go-loyalty-system/withdrawal/storage"
)

type Repo struct {
	User       user.Storage
	Balance    balance.Storage
	Order      order.Storage
	Withdrawal withdrawal.Storage
}

func NewDBRepo(url, driver string) (Repo, error) {
	if driver == "" {
		driver = `pgx`
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		return Repo{}, err
	}

	bs, err := balanceStorage.NewDBBalanceStorage(db)
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

	ws, err := withdrawalStorage.NewDBWithdrawalStorage(db)
	if err != nil {
		return Repo{}, err
	}

	return Repo{
		Balance:    bs,
		Order:      os,
		User:       us,
		Withdrawal: ws,
	}, nil
}
