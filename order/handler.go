package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/user"
	"io"
	"log"
	"net/http"
)

var errToStat = map[string]int{
	ErrSameUser:  http.StatusOK,
	ErrOtherUser: http.StatusConflict,
}

func GetOrders(accrualURL string, db Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		orders, err := db.FindAll(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		res := make([]Order, len(orders))
		for i, o := range orders {
			if o.Status == StatusNew || o.Status == StatusProcessing {
				upd, err := updateOrder(o, db, accrualURL)
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

func UpdateOrders(accrualURL string, db Storage) func(http.ResponseWriter, *http.Request) {
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

		if err = handleExistingOrder(db, order, u); err != nil {
			if code, ok := errToStat[err.Error()]; ok {
				w.WriteHeader(code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		accrual, err := getAccrual(accrualURL, order)
		if err != nil {
			log.Println(`ERROR`, err)
		}

		if _, err = db.Add(NewOrderFromAccrual(accrual, u)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set(`Content-Type`, `application/json`)
		w.WriteHeader(http.StatusAccepted)
	}
}

func updateOrder(o Order, db Storage, accrualURL string) (Order, error) {
	accrual, err := getAccrual(accrualURL, o.Number)
	if err != nil {
		fmt.Println(err)
		return o, err
	}

	upd := NewOrderFromOrderAndAccrual(o, accrual)
	if upd.Status != o.Status {
		if _, err = db.Update(upd); err != nil {
			return o, err
		}
	}

	return upd, nil
}

func handleExistingOrder(db Storage, order string, u string) error {
	o, err := db.Find(order)
	if err != nil || o.Number == `` {
		if errors.Is(err, sql.ErrNoRows) || o.Number == `` {
			return nil
		}
		return err
	}

	errStr := ErrOtherUser
	if o.User == u {
		errStr = ErrSameUser
	}

	return errors.New(errStr)
}
