package model

type Standing struct {
	TeamID     string `json:"team_id,omitempty"`
	Wins       int    `json:"wins,omitempty"`
	Losses     int    `json:"losses,omitempty"`
	TotalScore int    `json:"total_score,omitempty"`
}
