package main

import (
	"go-loyalty-system/internal"
	"go-loyalty-system/internal/configs"
	"log"
	"net/http"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.NewDBRepo(cfg.DBURL, ``)
	if err != nil {
		log.Fatal(err)
	}

	r := internal.NewRouter(db)
	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}
