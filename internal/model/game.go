package model

type GameKind int

const (
	PrelimGame     GameKind = 1
	TournamentGame GameKind = 2
)

type Game struct {
	ID      string   `json:"id,omitempty"`
	TeamIDs []string `json:"teams,omitempty"`
	Kind    GameKind `json:"kind,omitempty"`
}
