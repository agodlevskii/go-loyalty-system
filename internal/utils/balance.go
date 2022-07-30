package utils

import (
	"database/sql"
	"errors"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
)

func GetBalance(bs storage.BalanceStorage, user string) (models.Balance, error) {
	b, err := bs.Get(user)
	if errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
		err = bs.Set(b)
	}
	return b, err
}

func UpdateBalanceWithAccrual(bs storage.BalanceStorage, user string, accrual float64) (models.Balance, error) {
	b, err := bs.Get(user)
	if errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
	} else if err != nil {
		return b, err
	}

	b.Current += accrual
	return b, bs.Set(b)
}

func UpdateBalanceWithWithdrawal(bs storage.BalanceStorage, b models.Balance, w models.Withdrawal) (models.Balance, error) {
	b.Withdrawn += w.Sum
	b.Current -= w.Sum
	return b, bs.Set(b)
}
