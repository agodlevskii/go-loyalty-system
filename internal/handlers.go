package internal

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/balance"
	"go-loyalty-system/order"
	"go-loyalty-system/user/auth"
	"go-loyalty-system/withdrawal"
)

func NewRouter(db Repo) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.Middleware([]string{`/api/user/login`, `/api/user/register`}))

	r.Route(`/api/user`, func(r chi.Router) {
		r.Post(`/login`, auth.Login(db.User))
		r.Post(`/register`, auth.Register(db.User))

		r.Route(`/orders`, func(r chi.Router) {
			r.Get(`/`, order.GetOrders(db.Order))
			r.Post(`/`, order.UpdateOrders(db.Order))
		})

		r.Route(`/balance`, func(r chi.Router) {
			r.Get(`/`, balance.GetBalance(db.Balance))
			r.Post(`/withdraw`, withdrawal.Withdraw(db.Withdrawal))
			r.Get(`/withdrawals`, withdrawal.GetWithdrawals(db.Withdrawal))
		})
	})

	return r
}
