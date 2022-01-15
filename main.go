package main

import (
	"log"
	"net/http"

	"github.com/cszczepaniak/go-cribbly/network"
)

func main() {
	router := network.SetupRouter()

	s := http.Server{
		Addr:    `:8080`,
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
