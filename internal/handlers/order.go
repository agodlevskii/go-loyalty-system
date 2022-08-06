package handlers

import (
	"encoding/json"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
	"io"
	"net/http"
)

var errToStat = map[string]int{
	aerror.OrderExistsSameUser:  http.StatusOK,
	aerror.OrderExistsOtherUser: http.StatusConflict,
}

func GetOrders(accrual services.AccrualClient, oStorage storage.OrderStorage, bStorage storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(models.UserKey).(string)
		if !ok {
			HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusInternalServerError)
			return
		}

		orders, err := oStorage.FindAll(u)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if encerr := json.NewEncoder(w).Encode(orders); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.OrderFindAll, encerr), http.StatusInternalServerError)
		}

		for _, o := range orders {
			go func(o models.Order) {
				if o.Status == models.StatusNew || o.Status == models.StatusProcessing {
					services.UpdateOrderWithAccrual(o, oStorage, bStorage, accrual, u)
				}
			}(o)
		}
	}
}

func RegisterOrder(accrual services.AccrualClient, oStorage storage.OrderStorage, bStorage storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(models.UserKey).(string)
		if !ok {
			HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusInternalServerError)
			return
		}

		id, rerr := io.ReadAll(r.Body)
		if rerr != nil || id == nil || len(id) == 0 {
			HandleHTTPError(w, aerror.NewError(aerror.OrderNumberInvalid, rerr), http.StatusBadRequest)
			return
		}

		order := string(id)
		if !luhn.Validate(order) {
			HandleHTTPError(w, aerror.NewError(aerror.OrderNumberInvalid, nil), http.StatusUnprocessableEntity)
			return
		}

		if err := services.CheckExistingOrder(oStorage, order, u); err != nil {
			code, ok := errToStat[err.Label]
			if !ok {
				code = http.StatusInternalServerError
			}
			HandleHTTPError(w, err, code)
			return
		}

		if _, err := services.AddOrderFromAccrual(oStorage, bStorage, accrual, order, u); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
}
