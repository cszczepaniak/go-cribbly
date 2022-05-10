package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func newTestServer(t *testing.T) *httptest.Server {
	pcfg := persistence.NewMemoryConfig()
	handler := handlers.NewRequestHandler(pcfg)
	server := NewTestServer(handler)

	s := httptest.NewServer(server)
	t.Cleanup(s.Close)

	return s
}

func TestPing(t *testing.T) {
	s := newTestServer(t)

	resp, err := http.DefaultClient.Get(s.URL + `/ping`)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
