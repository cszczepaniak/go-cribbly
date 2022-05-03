package model

type Player struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type Team struct {
	ID      string   `json:"id,omitempty"`
	Players []Player `json:"players,omitempty"`
}
