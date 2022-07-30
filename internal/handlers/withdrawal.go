package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils"
	"go-loyalty-system/user"
	"log"
	"net/http"
)

func GetWithdrawals(db storage.WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
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

func Withdraw(bs storage.BalanceStorage, ws storage.WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wr models.Withdrawal
		u := r.Context().Value(user.Key).(string)
		if err := json.NewDecoder(r.Body).Decode(&wr); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !luhn.Validate(wr.Order) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		b, err := utils.GetBalance(bs, u)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if b.Current < wr.Sum {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		if err = ws.Add(models.NewWithdrawalFromWithdrawal(wr, u)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = utils.UpdateBalanceWithWithdrawal(bs, b, wr); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
