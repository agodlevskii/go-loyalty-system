package handlers

import (
	"encoding/json"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils"
	"go-loyalty-system/user"
	"log"
	"net/http"
)

func GetBalance(db storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		b, err := utils.GetBalance(db, u)
		if err == nil {
			w.Header().Set(`Content-Type`, `application/json`)
			err = json.NewEncoder(w).Encode(b)
		}

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
