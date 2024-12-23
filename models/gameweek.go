package models

type GameweekID int

type Gameweek struct {
	ID              GameweekID
	Name            string
	Deadline        string
	IsCurrent       bool
	IsNext          bool
	Finished        bool
	MostCaptainedID PlayerID
}

type FixtureID int

type Fixture struct {
	ID                 FixtureID
	Gameweek           *Gameweek
	HomeTeam           *Team
	AwayTeam           *Team
	HomeTeamDifficulty int
	AwayTeamDifficulty int
	DifficultyMajority int
}

func (f *Fixture) Players() []Player {
	players := make([]Player, 0)
	players = append(players, f.HomeTeam.Players...)
	players = append(players, f.AwayTeam.Players...)
	return players
}
