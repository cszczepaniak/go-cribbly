package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWinner(t *testing.T) {
	tests := []struct {
		game      Game
		expWinner string
	}{{
		game: Game{
			Scores: map[string]int{
				`a`: 0,
				`b`: 0,
			},
		},
		expWinner: ``,
	}, {
		game:      Game{},
		expWinner: ``,
	}, {
		game: Game{
			Scores: map[string]int{
				`a`: 0,
				`b`: 121,
			},
		},
		expWinner: `b`,
	}}

	for _, tc := range tests {
		got := tc.game.Winner()
		require.Equal(t, tc.expWinner, got)
	}
}
