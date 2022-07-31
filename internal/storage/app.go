package storage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go-loyalty-system/internal/aerror"
)

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

	db, sqlerr := sql.Open(driver, url)
	if sqlerr != nil {
		return Repo{}, aerror.NewError(aerror.RepoCreate, sqlerr)
	}

	bs, err := NewDBBalanceStorage(db)
	if err != nil {
		return Repo{}, err
	}

	os, err := NewDBOrderStorage(db)
	if err != nil {
		return Repo{}, err
	}

	us, err := NewDBUserStorage(db)
	if err != nil {
		return Repo{}, err
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
	}, aerror.NewEmptyError()
}
