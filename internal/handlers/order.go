package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils"
	"go-loyalty-system/user"
	"io"
	"net/http"
)

var errToStat = map[string]int{
	models.ErrSameUser:  http.StatusOK,
	models.ErrOtherUser: http.StatusConflict,
}

func GetOrders(accrualURL string, os storage.OrderStorage, bs storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		orders, err := os.FindAll(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				res[i] = upd
			} else {
				res[i] = o
			}
		}

		w.Header().Set(`Content-Type`, `application/json`)
		if err = json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func RegisterOrder(accrualURL string, os storage.OrderStorage, bs storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		id, err := io.ReadAll(r.Body)
		if err != nil || id == nil || len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		order := string(id)
		if !luhn.Validate(order) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err = utils.CheckExistingOrder(os, order, u); err != nil {
			if code, ok := errToStat[err.Error()]; ok {
				w.WriteHeader(code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if _, err = utils.AddOrderFromAccrual(os, bs, accrualURL, order, u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set(`Content-Type`, `application/json`)
		w.WriteHeader(http.StatusAccepted)
	}
}
