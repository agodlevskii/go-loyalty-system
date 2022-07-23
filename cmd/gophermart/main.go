package main

import (
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/internal/handlers"
	"go-loyalty-system/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewDBRepo(cfg.DBURL, ``)
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.NewRouter(db)
	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}
