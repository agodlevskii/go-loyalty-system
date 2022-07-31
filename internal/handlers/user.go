package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/internal/utils/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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

		if equal, err := auth.CompareHashes(reqUser.Password, dbUser.Password); err != nil || !equal {
			code := http.StatusInternalServerError
			if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) && !equal {
				code = http.StatusUnauthorized
			}
			HandleHTTPError(w, err, code)
			return
		}

		if token, err := auth.GetTokenFromUser(reqUser); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, auth.GetBearer(token))
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

		hash, err := auth.HashPassword(reqUser.Password)
		if err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		reqUser.Password = hash
		if err = db.Add(reqUser); err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, sql.ErrNoRows) {
				code = http.StatusConflict
			}
			HandleHTTPError(w, err, code)
			return
		}

		if token, err := auth.GetTokenFromUser(reqUser); err != nil {
			HandleHTTPError(w, err, http.StatusInternalServerError)
		} else {
			w.Header().Set(`Authorization`, auth.GetBearer(token))
			w.Write([]byte(token))
		}
	}
}

func AuthMiddleware(excludedPath []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludedPath {
				if strings.Contains(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			tknHdr := r.Header.Get(`Authorization`)
			if tknHdr == `` {
				HandleHTTPError(w, aerror.NewError(aerror.UserTokenIncorrect, nil), http.StatusUnauthorized)
				return
			}

			tknStr, err := auth.GetTokenFromBearer(tknHdr)
			if err != nil {
				HandleHTTPError(w, err, http.StatusUnauthorized)
				return
			}

			if valid, err := auth.IsTokenValid(tknStr); err != nil || !valid {
				code := http.StatusBadRequest
				if err != nil && errors.Is(err, jwt.ErrSignatureInvalid) || !valid {
					code = http.StatusUnauthorized
				}
				HandleHTTPError(w, err, code)
				return
			}

			usr, err := auth.GetUserFromToken(tknStr)
			if err != nil {
				HandleHTTPError(w, err, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), models.UserKey, usr.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
