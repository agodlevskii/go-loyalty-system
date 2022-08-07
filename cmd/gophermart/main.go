package main

import (
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/internal/handlers"
	"go-loyalty-system/internal/jobs"
	"go-loyalty-system/internal/logs"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
	"log"
	"net/http"
)

func main() {
	logs.InitLogger()

	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewDBRepo(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	accrual := services.NewAccrualClient(cfg.AccrualURL)
	if err = jobs.OrderStatusUpdate(db, accrual); err != nil {
		log.Fatal(err)
	}

	r := handlers.NewRouter(db, accrual)
	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}
