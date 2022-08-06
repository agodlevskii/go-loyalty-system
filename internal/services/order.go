package services

import (
	"database/sql"
	"errors"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
)

func CheckExistingOrder(db storage.OrderStorage, order string, user string) *aerror.AppError {
	o, err := db.Find(order)
	if err != nil || o.Number == "" {
		if err != nil && errors.Is(err, sql.ErrNoRows) || o.Number == "" {
			return nil
		}
		return err
	}

	errStr := aerror.OrderExistsOtherUser
	if o.User == user {
		errStr = aerror.OrderExistsSameUser
	}

	return aerror.NewError(errStr, err)
}

func AddOrderFromAccrual(oStorage storage.OrderStorage, bStorage storage.BalanceStorage, client AccrualClient, order, user string) (models.Order, *aerror.AppError) {
	accrual, err := client.GetAccrual(order)
	if err != nil {
		return models.Order{}, err
	}

	o, err := oStorage.Add(models.NewOrderFromAccrual(accrual, user))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		if dbOrder, err := oStorage.Find(o.Number); err == nil {
			if dbOrder.User == o.User {
				return dbOrder, aerror.NewError(aerror.OrderExistsSameUser, nil)
			}

			return models.Order{}, aerror.NewError(aerror.OrderExistsOtherUser, nil)
		}
	}

	_, err = updateBalanceWithAccrual(bStorage, user, accrual.Accrual)
	return o, err
}

func UpdateOrderWithAccrual(o models.Order, oStorage storage.OrderStorage, bStorage storage.BalanceStorage, client AccrualClient, user string) (models.Order, *aerror.AppError) {
	accrual, err := client.GetAccrual(o.Number)
	if err != nil {
		return o, err
	}

	upd := models.NewOrderFromOrderAndAccrual(o, accrual)
	if upd.Status != o.Status {
		if upd, err = oStorage.Update(upd); err != nil {
			return o, err
		}

		if upd.Accrual > 0 {
			if _, err := updateBalanceWithAccrual(bStorage, user, accrual.Accrual); err != nil {
				return o, err
			}
		}
	}

	return upd, nil
}
