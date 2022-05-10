package persistence

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/gameresults"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/games"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/teams"
)

type Config struct {
	GameStore games.GameStore
	gameresults.GameResultStore
	TeamStore teams.TeamStore
}

func newConfig(byteStore bytestore.ByteStore) *Config {
	return &Config{
		GameStore:       games.NewS3GameStore(byteStore),
		GameResultStore: gameresults.NewS3GameResultStore(byteStore),
		TeamStore:       teams.NewS3TeamStore(byteStore),
	}
}

func NewS3Config(awsSession *session.Session, bucket string, timeout time.Duration) *Config {
	return newConfig(bytestore.NewS3ByteStore(bucket, awsSession, timeout))
}

func NewMemoryConfig() *Config {
	return newConfig(bytestore.NewMemoryByteStore())
}
