package storage

import (
	"database/sql"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/user"
	userStorage "go-loyalty-system/user/storage"
)

type Repo struct {
	User       user.Storage
	Balance    BalanceStorage
	Order      OrderStorage
	Withdrawal WithdrawalStorage
}

func NewDBRepo(url, driver string) (Repo, *aerror.AppError) {
	if driver == "" {
		driver = `pgx`
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, err)
	}

	bs, err := NewDBBalanceStorage(db)
	if err != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, err)
	}

	os, err := NewDBOrderStorage(db)
	if err != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, err)
	}

	us, err := userStorage.NewDBUserStorage(db)
	if err != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, err)
	}

	ws, err := NewDBWithdrawalStorage(db)
	if err != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, err)
	}

	return Repo{
		Balance:    bs,
		Order:      os,
		User:       us,
		Withdrawal: ws,
	}, nil
}
