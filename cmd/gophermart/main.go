package main

import (
	"go-loyalty-system/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := handlers.NewRouter(nil)
	log.Fatal(http.ListenAndServe(`localhost:8080`, r))
}
