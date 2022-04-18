package main

import (
	"log"
	"os"

	"github.com/cszczepaniak/go-cribbly/server"
)

func main() {
	l := log.New(os.Stdout, ``, log.Flags())
	s := server.NewServer(l)
	err := s.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
