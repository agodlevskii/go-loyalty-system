package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-loyalty-system/user"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Login(db user.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser user.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dbUser, err := db.Find(reqUser.Login)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if equal, err := compareHashes(reqUser.Password, dbUser.Password); err != nil || !equal {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) && !equal {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if token, err := getTokenFromUser(reqUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, getBearer(token))
			w.Write([]byte(token))
		}
	}
}

func Register(db user.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser user.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash, err := hashPassword(reqUser.Password)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqUser.Password = hash
		if err = db.Add(reqUser); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if token, err := getTokenFromUser(reqUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, getBearer(token))
			w.Write([]byte(token))
		}
	}
}
