package model

type Player struct {
	ID   string
	Name string
}

func NewPlayer(id, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}
