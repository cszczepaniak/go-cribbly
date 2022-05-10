package standings

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func newResult(winnerID, loserID string, loserScore int) model.GameResult {
	return model.GameResult{
		GameID:     random.UUID(),
		Winner:     winnerID,
		Loser:      loserID,
		LoserScore: loserScore,
	}
}

func TestGetStandings(t *testing.T) {
	pcfg := persistence.NewMemoryConfig()
	gameResultStore := pcfg.GameResultStore
	standingService := &StandingsService{
		pcfg: pcfg,
	}

	t1, t2, t3, t4 := random.UUID(), random.UUID(), random.UUID(), random.UUID()

	// We'll create the following state:
	// Team 1: 3-0, 363 total points
	// Team 2: 2-1, 260 total points
	// Team 3: 2-1, 250 total points
	// Team 4: 1-2, 250 total points
	// The standings should come out in that order; we'll manufacture bogus team IDs to get the scores and wins where we want them

	_, err := gameResultStore.Create(newResult(t1, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(t1, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(t1, random.UUID(), 1))
	require.NoError(t, err)

	_, err = gameResultStore.Create(newResult(t2, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(t2, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(random.UUID(), t2, 18))
	require.NoError(t, err)

	_, err = gameResultStore.Create(newResult(t3, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(t3, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(random.UUID(), t3, 8))
	require.NoError(t, err)

	_, err = gameResultStore.Create(newResult(t4, random.UUID(), 1))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(random.UUID(), t4, 100))
	require.NoError(t, err)
	_, err = gameResultStore.Create(newResult(random.UUID(), t4, 29))
	require.NoError(t, err)

	results, err := standingService.GetStandings()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 4)

	assert.Equal(t, t1, results[0].TeamID)
	assert.Equal(t, 363, results[0].TotalScore)
	assert.Equal(t, 3, results[0].Wins)
	assert.Equal(t, 0, results[0].Losses)

	assert.Equal(t, t2, results[1].TeamID)
	assert.Equal(t, 260, results[1].TotalScore)
	assert.Equal(t, 2, results[1].Wins)
	assert.Equal(t, 1, results[1].Losses)

	assert.Equal(t, t3, results[2].TeamID)
	assert.Equal(t, 250, results[2].TotalScore)
	assert.Equal(t, 2, results[2].Wins)
	assert.Equal(t, 1, results[2].Losses)

	assert.Equal(t, t4, results[3].TeamID)
	assert.Equal(t, 250, results[3].TotalScore)
	assert.Equal(t, 1, results[3].Wins)
	assert.Equal(t, 2, results[3].Losses)
}
