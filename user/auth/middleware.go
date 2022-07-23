package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

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
}
