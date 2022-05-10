package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
	"github.com/cszczepaniak/go-cribbly/service/standings"
)

func newTestServer(t *testing.T) (*httptest.Server, *persistence.Config) {
	pcfg := persistence.NewMemoryConfig()
	standingsService := standings.NewStandingsService(pcfg)
	handler := handlers.NewRequestHandler(pcfg, standingsService)
	server := NewTestServer(handler)

	s := httptest.NewServer(server)
	t.Cleanup(s.Close)

	return s, pcfg
}

func TestPing(t *testing.T) {
	s, _ := newTestServer(t)

	resp, err := http.DefaultClient.Get(s.URL + `/ping`)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
