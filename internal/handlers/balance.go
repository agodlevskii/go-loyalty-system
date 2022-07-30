package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/user"
	"log"
	"net/http"
)

func GetBalance(db storage.BalanceStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.Key).(string)
		b, err := getBalance(db, u)
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

func getBalance(db storage.BalanceStorage, user string) (models.Balance, error) {
	b, err := db.Get(user)
	if errors.Is(err, sql.ErrNoRows) {
		b = models.NewBalance(user)
		err = db.Set(b)
	}
	return b, err
}
