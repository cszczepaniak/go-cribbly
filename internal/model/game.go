package model

type GameKind string

const (
	PrelimGame     GameKind = `prelim`
	TournamentGame GameKind = `tournament`
)

type Game struct {
	ID      string   `json:"id,omitempty"`
	TeamIDs []string `json:"teams,omitempty"`
	Kind    GameKind `json:"kind,omitempty"`
}

type GameResult struct {
	ID         string `json:"id,omitempty"`
	GameID     string `json:"game_id,omitempty"`
	Winner     string `json:"winner,omitempty"`
	Loser      string `json:"loser,omitempty"`
	LoserScore int    `json:"loser_score,omitempty"`
}
