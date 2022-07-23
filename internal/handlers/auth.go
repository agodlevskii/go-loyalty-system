package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-loyalty-system/internal/models"
	"go-loyalty-system/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("my_secret_key0")

type Claims struct {
	User string
	jwt.StandardClaims
}

func login(db storage.UserStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dbUser, err := db.GetUser(user.Login)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hashUser, err := hashPassword(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(hashUser.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if token, err := getToken(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write([]byte(token))
		}
	}
}

func register(db storage.UserStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Password = string(hash)
		if err = db.AddUser(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if token, err := getToken(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write([]byte(token))
		}
	}
}

func getToken(user models.User) (string, error) {
	exp := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		User: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func hashPassword(user models.User) (models.User, error) {
	if hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		return models.User{}, err
	} else {
		return models.User{
			Login:    user.Login,
			Password: string(hash),
		}, nil
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, `/user/login`) || strings.Contains(r.URL.Path, `/user/register`) {
			next.ServeHTTP(w, r)
			return
		}

		tknHdr := r.Header.Get(`Authorization`)
		if tknHdr == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknStr := strings.Split(tknHdr, `Bearer `)[1]
		claims := &Claims{}
		if tkn, err := jwt.ParseWithClaims(tknStr, claims, keyFn); err != nil || !tkn.Valid {
			if errors.Is(err, jwt.ErrSignatureInvalid) || !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func keyFn(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}
