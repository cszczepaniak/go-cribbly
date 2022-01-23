package model

type Player struct {
	ID string

	// TODO add other fields as necessary
}

func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
	}
}
