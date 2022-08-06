package handlers

import (
	"encoding/json"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
	"net/http"
)

func GetBalance(db storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(models.UserKey).(string)
		if !ok {
			HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusInternalServerError)
			return
		}

		b, err := services.GetBalance(db, u)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if encerr := json.NewEncoder(w).Encode(b); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.BalanceGet, encerr), http.StatusInternalServerError)
		}
	}
}
