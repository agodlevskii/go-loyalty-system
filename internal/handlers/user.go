package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-loyalty-system/internal/aerror"
	auth2 "go-loyalty-system/internal/auth"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Login(db storage.UserStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.User
		if encerr := json.NewDecoder(r.Body).Decode(&reqUser); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.UserRequestIncorrect, encerr), http.StatusBadRequest)
			return
		}

		dbUser, err := db.Find(reqUser.Login)
		if err != nil {
			HandleHTTPError(w, err, http.StatusUnauthorized)
			return
		}

		if equal, err := auth2.CompareHashes(reqUser.Password, dbUser.Password); err != nil || !equal {
			code := http.StatusInternalServerError
			if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) && !equal {
				code = http.StatusUnauthorized
			}
			HandleHTTPError(w, err, code)
			return
		}

		if token, err := auth2.GetTokenFromUser(reqUser); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, auth2.GetBearer(token))
			w.Write([]byte(token))
		}
	}
}

func Register(db storage.UserStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.User
		if encerr := json.NewDecoder(r.Body).Decode(&reqUser); encerr != nil {
			HandleHTTPError(w, aerror.NewError(aerror.UserRequestIncorrect, encerr), http.StatusBadRequest)
			return
		}

		hash, err := auth2.HashPassword(reqUser.Password)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		reqUser.Password = hash
		if err = db.Add(reqUser); err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, aerror.NewError(aerror.UserAdd, sql.ErrNoRows)) {
				code = http.StatusConflict
			}
			HandleHTTPError(w, err, code)
			return
		}

		if token, err := auth2.GetTokenFromUser(reqUser); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, auth2.GetBearer(token))
			w.Write([]byte(token))
		}
	}
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tknHdr := r.Header.Get(`Authorization`)
			if tknHdr == "" {
				HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusUnauthorized)
				return
			}

			tknStr, err := auth2.GetTokenFromBearer(tknHdr)
			if err != nil {
				HandleHTTPError(w, err, http.StatusUnauthorized)
				return
			}

			if valid, err := auth2.IsTokenValid(tknStr); err != nil || !valid {
				code := http.StatusBadRequest
				if err != nil && errors.Is(err, jwt.ErrSignatureInvalid) || !valid {
					code = http.StatusUnauthorized
				}
				HandleHTTPError(w, err, code)
				return
			}

			usr, err := auth2.GetUserFromToken(tknStr)
			if err != nil {
				HandleHTTPError(w, err, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), models.UserKey, usr.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
