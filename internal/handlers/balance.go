package handlers

import (
	"encoding/json"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils"
	"go-loyalty-system/user"
	"net/http"
)

func GetBalance(db storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		b, err := utils.GetBalance(db, u)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
		}

		w.Header().Set(`Content-Type`, `application/json`)
		if encerr := json.NewEncoder(w).Encode(b); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.BalanceGet, encerr), http.StatusInternalServerError)
		}
	}
}
