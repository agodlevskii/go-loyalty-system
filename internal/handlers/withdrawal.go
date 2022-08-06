package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vivekmurali/luhn"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
	"net/http"
	"sync"
)

var mu = sync.Mutex{}

func GetWithdrawals(db storage.WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(models.UserKey).(string)
		if !ok {
			HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusInternalServerError)
			return
		}

		ws, err := db.FindAll(u)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, aerror.NewError(aerror.WithdrawalFindAll, sql.ErrNoRows)) {
				code = http.StatusNoContent
			}

			HandleHTTPError(w, err, code)
			return
		}

		w.Header().Set(`Content-Type`, `application/json`)
		if encerr := json.NewEncoder(w).Encode(ws); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.WithdrawalFindAll, encerr), http.StatusInternalServerError)
		}
	}
}

func Withdraw(bs storage.BalanceStorage, ws storage.WithdrawalStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wr models.Withdrawal
		u, ok := r.Context().Value(models.UserKey).(string)
		if !ok {
			HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusInternalServerError)
			return
		}

		if encerr := json.NewDecoder(r.Body).Decode(&wr); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.WithdrawalRequestInvalid, encerr), http.StatusBadRequest)
			return
		}

		if !luhn.Validate(wr.Order) {
			HandleHTTPError(w, aerror.NewError(aerror.OrderNumberInvalid, nil), http.StatusUnprocessableEntity)
			return
		}

		b, err := services.GetBalance(bs, u)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		if b.Current < wr.Sum {
			HandleHTTPError(w, aerror.NewError(aerror.BalanceInsufficient, nil), http.StatusPaymentRequired)
			return
		}

		if err = ws.Add(models.NewWithdrawalFromWithdrawal(wr, u)); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if _, err = services.UpdateBalanceWithWithdrawal(bs, b, wr); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
