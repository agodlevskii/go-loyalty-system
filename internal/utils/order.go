package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"log"
)

func AddOrderFromAccrual(os storage.OrderStorage, bs storage.BalanceStorage, accrualURL, order, user string) (models.Order, error) {
	accrual, err := GetAccrual(accrualURL, order)
	if err != nil {
		log.Println(`ERROR`, err)
	}

	o, err := os.Add(models.NewOrderFromAccrual(accrual, user))
	if err == nil {
		_, err = UpdateBalanceWithAccrual(bs, user, accrual.Accrual)
	}
	return o, err
}

func CheckExistingOrder(db storage.OrderStorage, order string, user string) error {
	o, err := db.Find(order)
	if err != nil || o.Number == `` {
		if errors.Is(err, sql.ErrNoRows) || o.Number == `` {
			return nil
		}
		return err
	}

	errStr := models.ErrOtherUser
	if o.User == user {
		errStr = models.ErrSameUser
	}

	return errors.New(errStr)
}

func UpdateOrderWithAccrual(o models.Order, os storage.OrderStorage, bs storage.BalanceStorage, accrualURL, user string) (models.Order, error) {
	accrual, err := GetAccrual(accrualURL, o.Number)
	if err != nil {
		fmt.Println(err)
		return o, err
	}

	upd := models.NewOrderFromOrderAndAccrual(o, accrual)
	if upd.Status != o.Status {
		if upd, err = os.Update(upd); err != nil {
			return o, err
		}

		if upd.Accrual > 0 {
			if _, err := UpdateBalanceWithAccrual(bs, user, accrual.Accrual); err != nil {
				return o, err
			}
		}
	}

	return upd, nil
}
