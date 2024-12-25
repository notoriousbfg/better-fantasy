package models

import (
	"sort"
)

type PlayerTypeID int

const (
	PTGoalkeeper PlayerTypeID = iota + 1
	PTDefender
	PTMidfielder
	PTForward
)

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

// goals & assists
func (p *Player) AttackingForm(weeks int) float32 {
	if weeks == 0 || len(p.History) == 0 {
		return 0
	}
	history := make([]PlayerFixture, 0)
	for _, fixture := range p.History {
		history = append(history, fixture)
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].FixtureID > history[j].FixtureID
	})
	var lastPlayed []PlayerFixture
	if weeks > len(p.History) {
		lastPlayed = history
	} else {
		lastPlayed = history[len(history)-weeks:]
	}
	goalPoints := 0
	assistPoints := 0
	switch p.Type.ID {
	case PTGoalkeeper:
		for _, fixture := range lastPlayed {
			goalPoints += fixture.GoalsScored * 10
			assistPoints += fixture.Assists * 3
		}
	case PTDefender:
		for _, fixture := range lastPlayed {
			goalPoints += fixture.GoalsScored * 6
			assistPoints += fixture.Assists * 3
		}
	case PTMidfielder:
		for _, fixture := range lastPlayed {
			goalPoints += fixture.GoalsScored * 5
			assistPoints += fixture.Assists * 3
		}
	case PTForward:
		for _, fixture := range lastPlayed {
			goalPoints += fixture.GoalsScored * 4
			assistPoints += fixture.Assists * 3
		}
	}
	return (float32(goalPoints) + float32(assistPoints)) / float32(weeks)
}

// clean sheets
func (p *Player) DefendingForm(weeks int) float32 {
	return 0
}

type PlayerFixture struct {
	FixtureID   FixtureID
	PlayerID    PlayerID
	Minutes     int
	Played      bool
	Points      int
	GoalsScored int
	Assists     int
	YellowCards int
	RedCards    int
	Bonus       int
	WasHome     bool
}
