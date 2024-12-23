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

type PlayerFixture struct {
	FixtureID FixtureID
	PlayerID  PlayerID
	Minutes   int
	Played    bool
	Points    int
}
