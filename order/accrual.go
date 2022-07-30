package order

import (
	"encoding/json"
	"errors"
	"net/http"
)

const URL = `/api/orders`

func getAccrual(accrualURL string, order string) (AccrualOrder, error) {
	accrual := AccrualOrder{
		Order:   order,
		Status:  StatusNew,
		Accrual: 0,
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, accrualURL+URL+`/`+order, nil)
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

	return accrual, json.NewDecoder(res.Body).Decode(&accrual)
}
