package handlers

import (
	"encoding/json"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils"
	"go-loyalty-system/user"
	"io"
	"net/http"
)

var errToStat = map[string]int{
	aerror.OrderExistsSameUser:  http.StatusOK,
	aerror.OrderExistsOtherUser: http.StatusConflict,
}

func GetOrders(accrualURL string, os storage.OrderStorage, bs storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		orders, err := os.FindAll(u)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		res := make([]models.Order, len(orders))
		for i, o := range orders {
			if o.Status == models.StatusNew || o.Status == models.StatusProcessing {
				upd, err := utils.UpdateOrderWithAccrual(o, os, bs, accrualURL, u)
				if err != nil {
					HandleHTTPError(w, err, http.StatusInternalServerError)
					return
				}

				res[i] = upd
			} else {
				res[i] = o
			}
		}

		w.Header().Set(`Content-Type`, `application/json`)
		if encerr := json.NewEncoder(w).Encode(res); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.OrderFindAll, encerr), http.StatusInternalServerError)
		}
	}
}

func RegisterOrder(accrualURL string, os storage.OrderStorage, bs storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
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

		if err := utils.CheckExistingOrder(os, order, u); err != nil {
			code, ok := errToStat[err.Error()]
			if !ok {
				code = http.StatusInternalServerError
			}
			HandleHTTPError(w, err, code)
			return
		}

		if _, err := utils.AddOrderFromAccrual(os, bs, accrualURL, order, u); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set(`Content-Type`, `application/json`)
		w.WriteHeader(http.StatusAccepted)
	}
}
