package services

import (
	"database/sql"
	"errors"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
)

func GetBalance(bs storage.BalanceStorage, user string) (models.Balance, *aerror.AppError) {
	b, err := bs.Get(user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
		err = bs.Set(b)
	}
	return b, err
}

func updateBalanceWithAccrual(bs storage.BalanceStorage, user string, accrual float64) (models.Balance, *aerror.AppError) {
	b, err := bs.Get(user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
	} else if err != nil {
		return b, err
	}

	b.Current += accrual
	return b, bs.Set(b)
}
