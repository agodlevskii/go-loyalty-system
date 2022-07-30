package utils

import (
	"database/sql"
	"errors"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
)

func UpdateBalanceWithAccrual(bs storage.BalanceStorage, user string, accrual float64) (models.Balance, error) {
	b, err := bs.Get(user)
	if errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
		err = bs.Set(b)
	}

	if err != nil {
		return b, err
	}

	b.Current += accrual
	return b, bs.Set(b)
}
