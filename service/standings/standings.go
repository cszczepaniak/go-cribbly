package standings

import (
	"sort"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
)

type StandingsService struct {
	pcfg *persistence.Config
}

func NewStandingsService(pcfg *persistence.Config) *StandingsService {
	return &StandingsService{
		pcfg: pcfg,
	}
}

func (s *StandingsService) GetStandings() ([]model.Standing, error) {
	results, err := s.pcfg.GameResultStore.GetAll()
	if err != nil {
		return nil, err
	}

	standingsMap := make(map[string]model.Standing)
	for _, r := range results {
		winner, ok := standingsMap[r.Winner]
		if !ok {
			winner = model.Standing{
				TeamID: r.Winner,
			}
		}
		winner.Wins++
		winner.TotalScore += 121
		standingsMap[r.Winner] = winner

		loser, ok := standingsMap[r.Loser]
		if !ok {
			loser = model.Standing{
				TeamID: r.Loser,
			}
		}
		loser.Losses++
		loser.TotalScore += r.LoserScore
		standingsMap[r.Loser] = loser
	}

	res := make([]model.Standing, 0, len(standingsMap))
	for _, s := range standingsMap {
		res = append(res, s)
	}
	sort.Slice(res, func(i, j int) bool {
		// sort by wins, then total score, both descending
		if res[i].Wins == res[j].Wins {
			return res[i].TotalScore > res[j].TotalScore
		}
		return res[i].Wins > res[j].Wins
	})
	return res, nil
}
