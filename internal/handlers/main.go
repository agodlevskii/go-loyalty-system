package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/internal/storage"
)

func NewRouter(db storage.Repo) *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth)

	r.Route(`/api/user`, func(r chi.Router) {
		r.Post(`/login`, login(db.User))
		r.Post(`/register`, register(db.User))

		r.Route(`/orders`, func(r chi.Router) {
			r.Get(`/`, getOrders(db))
			r.Post(`/`, updateOrders(db))
		})

		r.Route(`/balance`, func(r chi.Router) {
			r.Get(`/`, getBalance(db))
			r.Post(`/withdraw`, withdraw(db))
			r.Get(`/withdrawals`, getWithdrawals(db))
		})
	})

	return r
}
