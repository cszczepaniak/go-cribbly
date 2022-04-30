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
