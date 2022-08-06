package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/logs"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
	"net/http"
)

func NewRouter(db storage.Repo, accrual services.AccrualClient) *chi.Mux {
	r := chi.NewRouter()
	r.Use(AuthMiddleware())

	r.Route("/api/user", func(r chi.Router) {
		withAuthRouter := r.With(AuthMiddleware())

		r.Post("/login", Login(db.User))
		r.Post("/register", Register(db.User))
		withAuthRouter.Get("/withdrawals", GetWithdrawals(db.Withdrawal))

		withAuthRouter.Route("/orders", func(r chi.Router) {
			withAuthRouter.Get("/", GetOrders(accrual, db.Order, db.Balance))
			withAuthRouter.Post("/", RegisterOrder(accrual, db.Order, db.Balance))
		})

		withAuthRouter.Route("/balance", func(r chi.Router) {
			withAuthRouter.Get("/", GetBalance(db.Balance))
			withAuthRouter.Post("/withdraw", Withdraw(db.Balance, db.Withdrawal))
		})
	})

	return r
}

func HandleHTTPError(w http.ResponseWriter, err *aerror.AppError, code int) {
	logs.Logger.Error(err.Error())
	w.WriteHeader(code)
}
