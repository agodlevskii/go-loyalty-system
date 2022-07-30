package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/internal/storage"
	"go-loyalty-system/user/auth"
)

var nonValidatedRoutes = []string{`/api/user/login`, `/api/user/register`}

func NewRouter(cfg *configs.Config, db storage.Repo) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.Middleware(nonValidatedRoutes))

	r.Route(`/api/user`, func(r chi.Router) {
		r.Post(`/login`, auth.Login(db.User))
		r.Post(`/register`, auth.Register(db.User))

		r.Route(`/orders`, func(r chi.Router) {
			r.Get(`/`, GetOrders(cfg.AccrualURL, db.Order))
			r.Post(`/`, RegisterOrder(cfg.AccrualURL, db.Order, db.Balance))
		})

		r.Route(`/balance`, func(r chi.Router) {
			r.Get(`/`, GetBalance(db.Balance))
			r.Post(`/withdraw`, Withdraw(db.Balance, db.Withdrawal))
			r.Get(`/withdrawals`, GetWithdrawals(db.Withdrawal))
		})
	})

	return r
}
