package game

const (
	WinningScore = 121
)

type Game struct {
	ID string
	// Scores is a mapping from team ID to score. It can be assumed that Scores will either be empty or
	// have exactly two elements in it, and that exactly one element will be WinningScore.
	Scores map[string]int `json:"scores,omitempty"`
}

func NewGame(id, teamA, teamB string) *Game {
	return &Game{
		ID: id,
		Scores: map[string]int{
			teamA: 0,
			teamB: 0,
		},
	}
}

// Winner returns the ID of the team with WinningScore, or an empty string if there is no winner yet.
func (g *Game) Winner() string {
	for teamID, sc := range g.Scores {
		if sc == WinningScore {
			return teamID
		}
	}
	return ``
}
