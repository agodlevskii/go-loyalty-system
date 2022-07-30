package order

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		usr, ok := r.Context().Value(user.Key).(string)
		if !ok || usr == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		orders, err := db.FindAll(usr)
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
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func UpdateOrders(accrualURL string, db Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, ok := r.Context().Value(user.Key).(string)
		if !ok || usr == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := io.ReadAll(r.Body)
		if err != nil || id == nil || len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !validateOrderNumber(string(id)) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err = handleExistingOrder(db, string(id), usr); err != nil {
			if code, ok := errToStat[err.Error()]; ok {
				w.WriteHeader(code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		accrual, err := getAccrual(accrualURL, string(id))
		if err != nil {
			log.Println(`ERROR`, err)
		}

		if _, err = db.Add(getOrderFromAccrual(accrual, usr)); err != nil {
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
		return o, err
	}

	upd := combineOrderAndAccrual(o, accrual)
	if upd.Status != o.Status {
		if _, err = db.Add(upd); err != nil {
			return o, err
		}
	}

	return upd, nil
}

func handleExistingOrder(db Storage, order string, usr string) error {
	o, err := db.Find(order)
	if err != nil || o.Number == `` {
		if errors.Is(err, sql.ErrNoRows) || o.Number == `` {
			return nil
		}
		return err
	}

	errStr := ErrOtherUser
	if o.User == usr {
		errStr = ErrSameUser
	}

	return errors.New(errStr)
}
