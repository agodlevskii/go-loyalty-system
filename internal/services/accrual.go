package services

import (
	"context"
	"encoding/json"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"net/http"
	"time"
)

type AccrualClient struct {
	client *http.Client
	url    string
}

func NewAccrualClient(url string) AccrualClient {
	return AccrualClient{
		client: &http.Client{},
		url:    url,
	}
}

func (a AccrualClient) GetAccrual(order string) (models.AccrualOrder, *aerror.AppError) {
	accrual := models.AccrualOrder{
		Order:   order,
		Status:  models.StatusNew,
		Accrual: 0,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url+"/api/orders/"+order, nil)
	if err != nil {
		return accrual, aerror.NewError(aerror.AccrualGet, err)
	}

	res, err := a.client.Do(req)
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
