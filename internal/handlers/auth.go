package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("my_secret_key0")

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	User string
	jwt.StandardClaims
}

func login(db interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// ToDo: Validate user

		if token, err := getToken(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write([]byte(token))
		}
	}
}

func register(db interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// ToDo: Save user
		user := User{
			Login:    "test",
			Password: "testpwd",
		}

		if token, err := getToken(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write([]byte(token))
		}
	}
}

func getToken(user User) (string, error) {
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
