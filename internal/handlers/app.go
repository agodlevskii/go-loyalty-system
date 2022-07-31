package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

var nonValidatedRoutes = []string{`/api/user/login`, `/api/user/register`}

func NewRouter(cfg *configs.Config, db storage.Repo) *chi.Mux {
	r := chi.NewRouter()
	r.Use(AuthMiddleware(nonValidatedRoutes))

	r.Route(`/api/user`, func(r chi.Router) {
		r.Post(`/login`, Login(db.User))
		r.Post(`/register`, Register(db.User))
		r.Get(`/withdrawals`, GetWithdrawals(db.Withdrawal))

		r.Route(`/orders`, func(r chi.Router) {
			r.Get(`/`, GetOrders(cfg.AccrualURL, db.Order, db.Balance))
			r.Post(`/`, RegisterOrder(cfg.AccrualURL, db.Order, db.Balance))
		})

		r.Route(`/balance`, func(r chi.Router) {
			r.Get(`/`, GetBalance(db.Balance))
			r.Post(`/withdraw`, Withdraw(db.Balance, db.Withdrawal))
		})
	})

	return r
}

func HandleHTTPError(w http.ResponseWriter, err *aerror.AppError, code int) {
	zap.Error(err)
	w.WriteHeader(code)
}
