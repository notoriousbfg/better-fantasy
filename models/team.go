package models

type TeamConfig struct {
	Players   []StartingPlayer
	BankValue float32
}

type TeamID int

type Team struct {
	ID        TeamID
	Name      string
	ShortName string
	Players   []Player
	Fixtures  []Fixture
}

type StartingPlayer struct {
	Player       Player
	Fixture      Fixture
	OpposingTeam Team
	OverallRank  string
	TypeRank     string
}

type StartingEleven map[string][]StartingPlayer

func (se StartingEleven) PlayerCount() int {
	count := 0
	for position := range se {
		for range position {
			count++
		}
	}
	return count
}
