package common

import (
	"log"
	"net/http"

	"github.com/cszczepaniak/go-cribbly/config"
	"github.com/cszczepaniak/go-cribbly/internal/awscfg"
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/server"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func Start(listenAndServeFunc func(string, http.Handler) error) {
	config.Init()

	awsSession, err := awscfg.Connect()
	if err != nil {
		log.Fatal(err)
	}

	pcfg := persistence.NewS3Config(awsSession, *config.DataBucket, *config.ByteStoreTimeout)
	handler := handlers.NewRequestHandler(pcfg)

	s := server.NewServer(handler)
	err = listenAndServeFunc(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
