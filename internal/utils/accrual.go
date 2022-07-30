package utils

import (
	"encoding/json"
	"errors"
	"go-loyalty-system/internal/models"
	"net/http"
)

func GetAccrual(accrualURL string, order string) (models.AccrualOrder, error) {
	accrual := models.AccrualOrder{
		Order:   order,
		Status:  models.StatusNew,
		Accrual: 0,
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, accrualURL+`/api/orders/`+order, nil)
	if err != nil {
		return accrual, err
	}

	res, err := client.Do(req)
	if err != nil {
		return accrual, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		return accrual, errors.New(http.StatusText(res.StatusCode))
	}

	if res.StatusCode == http.StatusNoContent {
		return accrual, nil
	}

	return accrual, json.NewDecoder(res.Body).Decode(&accrual)
}
