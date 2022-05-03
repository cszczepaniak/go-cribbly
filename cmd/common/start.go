package common

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/cszczepaniak/go-cribbly/config"
	"github.com/cszczepaniak/go-cribbly/internal/awscfg"
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/server"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func Start(listenAndServeFunc func(string, http.Handler) error) {
	flag.Parse()

	awsSession, err := awscfg.Connect()
	if err != nil {
		log.Fatal(err)
	}

	pcfg := persistence.NewS3Config(*config.DataBucket, awsSession, time.Second)
	handler := handlers.NewRequestHandler(pcfg)

	s := server.NewServer(handler)
	err = listenAndServeFunc(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
