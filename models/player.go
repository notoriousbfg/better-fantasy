package models

type PlayerTypeID int

type PlayerType struct {
	ID               PlayerTypeID
	Name             string
	PluralName       string
	ShortName        string
	TeamPlayerCount  int
	TeamMinPlayCount int
	TeamMaxPlayCount int
}

type PlayerStats struct {
	Minutes       int
	Goals         int
	Assists       int
	Conceded      int
	CleanSheets   int
	YellowCards   int
	RedCards      int
	Bonus         int
	Starts        int
	AverageStarts float32
	MatchesPlayed float32
	ICTIndex      float32
	ICTIndexRank  int
}

type PlayerHistory struct {
	Fixture *Fixture
	Minutes int
}

type PlayerRoundProbability map[GameweekID]float32

type PlayerID int

type Player struct {
	ID               PlayerID
	Name             string
	Form             float32
	PointsPerGame    float32
	TotalPoints      int
	Cost             string
	RawCost          float32
	Team             *Team
	Type             PlayerType
	Stats            PlayerStats
	History          map[FixtureID]PlayerFixture
	ChanceOfPlaying  PlayerRoundProbability
	MostCaptained    bool
	PickedPercentage float32
}

func (p *Player) FormOverCost() float32 {
	if p.Form <= 0 || p.RawCost == 0 {
		return 0
	}
	return p.Form / p.RawCost
}

func (p *Player) PointsOverCost() float32 {
	if p.TotalPoints == 0 || p.RawCost == 0 {
		return 0
	}
	return float32(p.TotalPoints) / p.RawCost
}

type PlayerFixture struct {
	FixtureID FixtureID
	PlayerID  PlayerID
	Minutes   int
	Played    bool
	Points    int
}
