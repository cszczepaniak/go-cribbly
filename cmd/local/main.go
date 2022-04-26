package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/cszczepaniak/go-cribbly/internal/awscfg"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/s3"
	"github.com/cszczepaniak/go-cribbly/server"
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
	s3Client := s3.NewS3Client(bucket, awsSession)
	s := server.NewServer(s3Client)
	err = http.ListenAndServe(`:8080`, s)
	if err != nil {
		log.Fatal(err)
	}
}
