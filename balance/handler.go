package balance

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/user"
	"log"
	"net/http"
)

func GetAccount(db AccountStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		a, err := getAccount(db, u)
		if err == nil {
			w.Header().Set(`Content-Type`, `application/json`)
			err = json.NewEncoder(w).Encode(a)
		}

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func GetWithdrawals(db WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		ws, err := db.FindAll(u)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set(`Content-Type`, `application/json`)
		if err = json.NewEncoder(w).Encode(ws); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func Withdraw(as AccountStorage, ws WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wr Withdrawal
		u := r.Context().Value(user.Key).(string)
		if err := json.NewDecoder(r.Body).Decode(&wr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !luhn.Validate(wr.Order) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		a, err := getAccount(as, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if a.Current < wr.Sum {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		if err = ws.Add(NewWithdrawalFromRequest(wr, u)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.Withdrawn += wr.Sum
		a.Current -= wr.Sum
		if err = as.Set(a); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getAccount(db AccountStorage, user string) (Account, error) {
	b, err := db.Get(user)
	if errors.Is(err, sql.ErrNoRows) {
		b = NewAccount(user)
		err = db.Set(b)
	}
	return b, err
}
