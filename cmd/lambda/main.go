package main

import (
	"log"

	"github.com/apex/gateway"

	"github.com/cszczepaniak/go-cribbly/server"
)

func main() {
	s := server.NewServer()
	err := gateway.ListenAndServe(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
