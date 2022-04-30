package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cszczepaniak/go-cribbly/internal/awscfg"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/games"
	"github.com/cszczepaniak/go-cribbly/server"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func main() {
	awsSession, err := awscfg.Connect()
	if err != nil {
		log.Fatal(err)
	}
	bucket := os.Getenv(`CRIBBLY_DATA_BUCKET`)
	if bucket == `` {
		log.Fatal(errors.New(`bucket not set`))
	}
	byteStore := bytestore.NewS3ByteStore(bucket, awsSession, time.Second)
	gameStore := games.NewS3GameStore(byteStore)
	handler := handlers.NewRequestHandler(gameStore)
	s := server.NewServer(handler)
	err = http.ListenAndServe(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
