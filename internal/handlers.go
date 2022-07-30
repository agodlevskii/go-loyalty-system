package internal

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/balance"
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/order"
	"go-loyalty-system/user/auth"
)

func NewRouter(cfg *configs.Config, db Repo) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.Middleware([]string{`/api/user/login`, `/api/user/register`}))

	r.Route(`/api/user`, func(r chi.Router) {
		r.Post(`/login`, auth.Login(db.User))
		r.Post(`/register`, auth.Register(db.User))

		r.Route(`/orders`, func(r chi.Router) {
			r.Get(`/`, order.GetOrders(cfg.AccrualURL, db.Order))
			r.Post(`/`, order.UpdateOrders(cfg.AccrualURL, db.Order))
		})

		r.Route(`/balance`, func(r chi.Router) {
			r.Get(`/`, balance.GetAccount(db.Account))
			r.Post(`/withdraw`, balance.Withdraw(db.Account, db.Withdrawal))
			r.Get(`/withdrawals`, balance.GetWithdrawals(db.Withdrawal))
		})
	})

	return r
}
