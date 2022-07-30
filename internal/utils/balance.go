package utils

import (
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
)

func UpdateBalanceWithAccrual(bs storage.BalanceStorage, user string, accrual float64) (models.Balance, error) {
	b, err := bs.Get(user)
	if err != nil {
		return b, err
	}

	b.Current += accrual
	return b, bs.Set(b)
}
