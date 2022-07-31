package utils

import (
	"encoding/json"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"net/http"
)

func GetAccrual(accrualURL string, order string) (models.AccrualOrder, *aerror.AppError) {
	accrual := models.AccrualOrder{
		Order:   order,
		Status:  models.StatusNew,
		Accrual: 0,
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, accrualURL+`/api/orders/`+order, nil)
	if err != nil {
		return accrual, aerror.NewError(aerror.AccrualGet, err)
	}

	res, err := client.Do(req)
	if err != nil {
		return accrual, aerror.NewError(aerror.AccrualGet, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return accrual, nil
	}
	if err = json.NewDecoder(res.Body).Decode(&accrual); err != nil {
		return accrual, aerror.NewError(aerror.AccrualGet, err)
	}

	return accrual, nil
}
