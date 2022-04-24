package main

import (
	"log"
	"net/http"

	"github.com/cszczepaniak/go-cribbly/server"
)

func main() {
	s := server.NewServer()
	err := http.ListenAndServe(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
