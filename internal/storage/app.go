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

func NewDBRepo(url string) (Repo, error) {
	db, sqlerr := sql.Open("pgx", url)
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
	}, nil
}
