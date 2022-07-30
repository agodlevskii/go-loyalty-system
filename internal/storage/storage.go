package storage

import (
	"database/sql"
	"go-loyalty-system/user"
	userStorage "go-loyalty-system/user/storage"
)

type Repo struct {
	User       user.Storage
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

	bs, err := NewDBBalanceStorage(db)
	if err != nil {
		return Repo{}, err
	}

	os, err := NewDBOrderStorage(db)
	if err != nil {
		return Repo{}, err
	}

	us, err := userStorage.NewDBUserStorage(db)
	if err != nil {
		return Repo{}, err
	}

	ws, err := NewDBWithdrawalStorage(db)
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
