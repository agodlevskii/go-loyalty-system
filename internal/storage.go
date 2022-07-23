package internal

import (
	"database/sql"
	"go-loyalty-system/balance"
	"go-loyalty-system/order"
	"go-loyalty-system/user"
	"go-loyalty-system/user/storage"
	"go-loyalty-system/withdrawal"
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

	user, err := storage.NewDBUserStorage(db)
	if err != nil {
		return Repo{}, err
	}

	return Repo{User: user}, nil
}
