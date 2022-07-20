package main

import (
	"go-loyalty-system/internal/configs"
	"go-loyalty-system/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := handlers.NewRouter(nil)
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}
