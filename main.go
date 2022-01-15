package main

import (
	"log"
	"net/http"
	"os"

	gamerepo "github.com/cszczepaniak/go-cribbly/data/game/repository"
	"github.com/cszczepaniak/go-cribbly/data/persistence"
	"github.com/cszczepaniak/go-cribbly/network"
)

func main() {
	gameRepo := gamerepo.NewMemory()
	pcfg := &persistence.Config{
		GameRepository: gameRepo,
	}

	logger := log.New(os.Stdout, ``, log.LstdFlags)

	router := network.SetupRouter(logger, pcfg)

	s := http.Server{
		Addr:    `:8080`,
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
