package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"better-fantasy/models"

	"github.com/hashicorp/go-multierror"
)

const (
	fixturesApi       = "https://fantasy.premierleague.com/api/fixtures/"
	statsApi          = "https://fantasy.premierleague.com/api/bootstrap-static/"
	playerFixturesApi = "https://fantasy.premierleague.com/api/element-summary/"
)

type apiTeam struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type apiEvent struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Deadline        time.Time `json:"deadline_time"`
	IsCurrent       bool      `json:"is_current"`
	IsNext          bool      `json:"is_next"`
	Finished        bool      `json:"finished"`
	MostCaptainedID int       `json:"most_captained"`
}

type apiElement struct {
	ID                       int     `json:"id"`
	Name                     string  `json:"web_name"`
	Form                     string  `json:"form"`
	PointsPerGame            string  `json:"points_per_game"`
	TotalPoints              int     `json:"total_points"`
	Cost                     int     `json:"now_cost"`
	TypeID                   int     `json:"element_type"`
	TeamID                   int     `json:"team"`
	Minutes                  int     `json:"minutes"`
	Goals                    int     `json:"goals_scored"`
	Assists                  int     `json:"assists"`
	Conceded                 int     `json:"goals_conceded"`
	CleanSheets              int     `json:"clean_sheets"`
	YellowCards              int     `json:"yellow_cards"`
	RedCards                 int     `json:"red_cards"`
	Bonus                    int     `json:"bonus"`
	Starts                   int     `json:"starts"`
	StartsPerNinety          float32 `json:"starts_per_90"`
	ICTIndex                 string  `json:"ict_index"`
	ICTIndexRank             int     `json:"ict_index_rank"`
	News                     string  `json:"news"`
	ChanceOfPlayingThisRound *int    `json:"chance_of_playing_this_round"`
	ChanceOfPlayingNextRound *int    `json:"chance_of_playing_next_round"`
	SelectedByPercent        string  `json:"selected_by_percent"`
}

type apiElementType struct {
	ID           int    `json:"id"`
	Name         string `json:"singular_name"`
	PluralName   string `json:"plural_name"`
	ShortName    string `json:"singular_name_short"`
	PlayerCount  int    `json:"squad_select"`
	SquadMinPlay int    `json:"squad_min_play"`
	SquadMaxPlay int    `json:"squad_max_play"`
}

type apiStats struct {
	Teams        []apiTeam        `json:"teams"`
	Events       []apiEvent       `json:"events"`
	Elements     []apiElement     `json:"elements"`
	ElementTypes []apiElementType `json:"element_types"`
}

type apiPlayerFixturesAndHistory struct {
	History []apiPlayerHistory `json:"history"`
}

type apiPlayerHistory struct {
	ElementID   int  `json:"element"`
	FixtureID   int  `json:"fixture"`
	Minutes     int  `json:"minutes"`
	TotalPoints int  `json:"total_points"`
	GoalsScored int  `json:"goals_scored"`
	Assists     int  `json:"assists"`
	YellowCards int  `json:"yellow_cards"`
	RedCards    int  `json:"red_cards"`
	Bonus       int  `json:"bonus"`
	WasHome     bool `json:"was_home"`
}

type apiFixture struct {
	ID                 int `json:"id"`
	AwayTeamID         int `json:"team_a"`
	HomeTeamID         int `json:"team_h"`
	EventID            int `json:"event"`
	AwayTeamDifficulty int `json:"team_a_difficulty"`
	HomeTeamDifficulty int `json:"team_h_difficulty"`
}

type apiFixtures []apiFixture

type apiPicks struct {
	Picks        []apiPick       `json:"picks"`
	EntryHistory apiEntryHistory `json:"entry_history"`
}

type apiPick struct {
	Element   int  `json:"element"`
	IsCaptain bool `json:"is_captain"`
}

type apiEntryHistory struct {
	Bank float32 `json:"bank"`
}

type Data struct {
	PlayerTypes []models.PlayerType
	Gameweeks   []models.Gameweek
	Fixtures    []*models.Fixture
	Teams       []*models.Team
	Players     []models.Player
}

func (d *Data) FixturesByGameWeek(gameweek int) []models.Fixture {
	fixtures := make([]models.Fixture, 0)
	for _, fixture := range d.Fixtures {
		if models.GameweekID(gameweek) == fixture.Gameweek.ID {
			fixtures = append(fixtures, *fixture)
		}
	}
	return fixtures
}

func (d *Data) Gameweek(gw int) *models.Gameweek {
	for _, gameweek := range d.Gameweeks {
		if gameweek.ID == models.GameweekID(gw) {
			return &gameweek
		}
	}
	return nil
}

func (d *Data) CurrentGameweek() *models.Gameweek {
	for _, gameweek := range d.Gameweeks {
		if gameweek.IsCurrent {
			return &gameweek
		}
	}
	return nil
}

func (d *Data) PlayerType(pt string) *models.PlayerType {
	for _, playerType := range d.PlayerTypes {
		if playerType.Name == pt {
			return &playerType
		}
	}
	return nil
}

func (d *Data) GameweekPlayers(gameweek int) []models.StartingPlayer {
	gameweekPlayers := make([]models.StartingPlayer, 0)
	for _, fixture := range d.FixturesByGameWeek(gameweek) {
		for _, player := range fixture.HomeTeam.Players {
			gameweekPlayers = append(gameweekPlayers, models.StartingPlayer{
				Player:       player,
				Fixture:      fixture,
				OpposingTeam: *fixture.AwayTeam,
			})
		}
		for _, player := range fixture.AwayTeam.Players {
			gameweekPlayers = append(gameweekPlayers, models.StartingPlayer{
				Player:       player,
				Fixture:      fixture,
				OpposingTeam: *fixture.HomeTeam,
			})
		}
	}
	return gameweekPlayers
}

func (d *Data) GameweekPlayerSet(gameweek models.GameweekID) map[models.PlayerID]models.StartingPlayer {
	playerSet := make(map[models.PlayerID]models.StartingPlayer, 0)
	for _, player := range d.GameweekPlayers(int(gameweek)) {
		playerSet[player.Player.ID] = player
	}
	return playerSet
}

func (d *Data) RequestManagerPicks(managerID int) models.TeamConfig {
	endpoint := fmt.Sprintf("https://fantasy.premierleague.com/api/entry/%d/event/%d/picks/", managerID, d.CurrentGameweek().ID)

	teamBody, err := getJsonBody(endpoint)
	if err != nil {
		panic(err)
	}

	var apiPicks apiPicks
	if err := json.Unmarshal(teamBody, &apiPicks); err != nil {
		panic(err)
	}

	gameweekPlayerSet := d.GameweekPlayerSet(d.CurrentGameweek().ID)

	players := make([]models.StartingPlayer, 0)
	for _, pick := range apiPicks.Picks {
		thisPlayer, ok := gameweekPlayerSet[models.PlayerID(pick.Element)]
		if ok {
			players = append(players, thisPlayer)
		}
	}

	return models.TeamConfig{
		Players:   players,
		BankValue: apiPicks.EntryHistory.Bank,
	}
}

func FetchData() (*Data, error) {
	data := &Data{}

	statsApiBody, err := getJsonBody(statsApi)
	if err != nil {
		panic(err)
	}
	var statsResp apiStats
	if err := json.Unmarshal(statsApiBody, &statsResp); err != nil {
		panic(err)
	}

	// var currentGameweekID models.GameweekID

	gameweeksByID := make(map[models.GameweekID]*models.Gameweek, 0)
	for _, apiEvent := range statsResp.Events {
		gameweekID := models.GameweekID(apiEvent.ID)
		// if apiEvent.IsCurrent {
		// 	currentGameweekID = gameweekID
		// }
		gameweek := &models.Gameweek{
			ID:              gameweekID,
			Name:            apiEvent.Name,
			Deadline:        apiEvent.Deadline.Format("02 Jan 15:04"),
			IsCurrent:       apiEvent.IsCurrent,
			IsNext:          apiEvent.IsNext,
			Finished:        apiEvent.Finished,
			MostCaptainedID: models.PlayerID(apiEvent.MostCaptainedID),
		}
		gameweeksByID[gameweekID] = gameweek
		data.Gameweeks = append(data.Gameweeks, *gameweek)
	}

	var teams []*models.Team
	teamsByID := make(map[models.TeamID]*models.Team, 0)
	for _, apiTeam := range statsResp.Teams {
		newTeam := models.Team{
			ID:        models.TeamID(apiTeam.ID),
			Name:      apiTeam.Name,
			ShortName: apiTeam.ShortName,
		}
		teams = append(teams, &newTeam)
		teamsByID[newTeam.ID] = &newTeam
	}
	data.Teams = teams

	playerTypesByID := make(map[models.PlayerTypeID]models.PlayerType, 0)
	for _, apiElementType := range statsResp.ElementTypes {
		newType := models.PlayerType{
			ID:               models.PlayerTypeID(apiElementType.ID),
			Name:             apiElementType.Name,
			PluralName:       apiElementType.PluralName,
			ShortName:        apiElementType.ShortName,
			TeamPlayerCount:  apiElementType.PlayerCount,
			TeamMinPlayCount: apiElementType.SquadMinPlay,
			TeamMaxPlayCount: apiElementType.SquadMaxPlay,
		}
		playerTypeID := models.PlayerTypeID(apiElementType.ID)
		playerTypesByID[playerTypeID] = newType
		data.PlayerTypes = append(data.PlayerTypes, newType)
	}

	playerCount := len(statsResp.Elements)
	playersChannel := make(chan models.Player, playerCount)
	errorsChannel := make(chan error, playerCount)

	teamPlayersByID := make(map[models.TeamID][]models.Player, 0)
	allPlayers := make([]models.Player, 0)

	for _, apiPlayer := range statsResp.Elements {
		go func() {
			newPlayer, err := newPlayer(apiPlayer, teamsByID, playerTypesByID)
			if err != nil {
				errorsChannel <- err
				return
			}
			playersChannel <- newPlayer
		}()
	}

	var errors error
	for i := 0; i < playerCount; i++ {
		select {
		case player := <-playersChannel:
			allPlayers = append(allPlayers, player)
			teamPlayersByID[player.Team.ID] = append(
				teamPlayersByID[models.TeamID(player.Team.ID)],
				player,
			)
		case err := <-errorsChannel:
			errors = multierror.Append(errors, err)
		}
	}
	if errors != nil {
		return &Data{}, fmt.Errorf("there was a problem building players: %w", errors)
	}

	close(playersChannel)
	close(errorsChannel)

	data.Players = allPlayers

	for _, team := range teams {
		team.Players = teamPlayersByID[team.ID]
	}

	fixturesBody, err := getJsonBody(fixturesApi)
	if err != nil {
		panic(err)
	}

	var apiFixtures apiFixtures
	if err := json.Unmarshal(fixturesBody, &apiFixtures); err != nil {
		panic(err)
	}

	fixtures := make([]*models.Fixture, 0)
	for _, apiFixture := range apiFixtures {
		homeTeam := teamsByID[models.TeamID(apiFixture.HomeTeamID)]
		awayTeam := teamsByID[models.TeamID(apiFixture.AwayTeamID)]

		gameweek, ok := gameweeksByID[models.GameweekID(apiFixture.EventID)]
		if !ok {
			continue
		}

		newFixture := models.Fixture{
			ID:                 models.FixtureID(apiFixture.ID),
			Gameweek:           gameweek,
			HomeTeam:           homeTeam,
			AwayTeam:           awayTeam,
			HomeTeamDifficulty: apiFixture.HomeTeamDifficulty,
			AwayTeamDifficulty: apiFixture.AwayTeamDifficulty,
			DifficultyMajority: abs(apiFixture.HomeTeamDifficulty - apiFixture.AwayTeamDifficulty),
		}
		fixtures = append(fixtures, &newFixture)

		if team, ok := teamsByID[models.TeamID(apiFixture.HomeTeamID)]; ok {
			team.Fixtures = append(team.Fixtures, newFixture)
		}

		if team, ok := teamsByID[models.TeamID(apiFixture.AwayTeamID)]; ok {
			team.Fixtures = append(team.Fixtures, newFixture)
		}
	}
	data.Fixtures = fixtures

	return data, nil
}

func newPlayer(
	apiPlayer apiElement,
	teamsByID map[models.TeamID]*models.Team,
	playerTypesByID map[models.PlayerTypeID]models.PlayerType,
) (models.Player, error) {
	playerForm, err := strconv.ParseFloat(apiPlayer.Form, 32)
	if err != nil {
		return models.Player{}, err
	}

	playerPointsPerGame, err := strconv.ParseFloat(apiPlayer.PointsPerGame, 32)
	if err != nil {
		return models.Player{}, err
	}

	playerTeam, ok := teamsByID[models.TeamID(apiPlayer.TeamID)]
	if !ok {
		return models.Player{}, fmt.Errorf("missing team ID '%d'", apiPlayer.TeamID)
	}

	playerType, ok := playerTypesByID[models.PlayerTypeID(apiPlayer.TypeID)]
	if !ok {
		return models.Player{}, fmt.Errorf("missing player type ID '%d'", apiPlayer.TypeID)
	}

	ictIndex, err := strconv.ParseFloat(apiPlayer.ICTIndex, 32)
	if err != nil {
		return models.Player{}, err
	}

	formattedCost := fmt.Sprintf("£%.1fm", float32(apiPlayer.Cost)/float32(10))

	// var chanceOfPlayingThisRound float32
	// if apiPlayer.ChanceOfPlayingThisRound == nil {
	// 	chanceOfPlayingThisRound = 1
	// } else {
	// 	chanceOfPlayingThisRound = float32(*apiPlayer.ChanceOfPlayingThisRound) / 100
	// }

	// var chanceOfPlayingNextRound float32
	// if apiPlayer.ChanceOfPlayingNextRound == nil {
	// 	chanceOfPlayingNextRound = 1
	// } else {
	// 	chanceOfPlayingNextRound = float32(*apiPlayer.ChanceOfPlayingNextRound) / 100
	// }

	// chanceOfPlaying := map[models.GameweekID]float32{
	// 	currentGameweekID:     chanceOfPlayingThisRound,
	// 	currentGameweekID + 1: chanceOfPlayingNextRound, // assumes next round is gameweek ID + 1
	// }

	pickedPercentage, err := strconv.ParseFloat(apiPlayer.SelectedByPercent, 32)
	if err != nil {
		return models.Player{}, err
	}

	newPlayer := models.Player{
		ID:            models.PlayerID(apiPlayer.ID),
		Name:          apiPlayer.Name,
		Form:          float32(playerForm),
		PointsPerGame: float32(playerPointsPerGame),
		TotalPoints:   apiPlayer.TotalPoints,
		Cost:          formattedCost,
		RawCost:       float32(apiPlayer.Cost) / float32(10),
		Team:          playerTeam,
		Type:          playerType,
		Stats: models.PlayerStats{
			Minutes:       apiPlayer.Minutes,
			Goals:         apiPlayer.Goals,
			Assists:       apiPlayer.Assists,
			Conceded:      apiPlayer.Conceded,
			CleanSheets:   apiPlayer.CleanSheets,
			YellowCards:   apiPlayer.YellowCards,
			RedCards:      apiPlayer.RedCards,
			Bonus:         apiPlayer.Bonus,
			Starts:        apiPlayer.Starts,
			AverageStarts: apiPlayer.StartsPerNinety,
			ICTIndex:      float32(ictIndex),
			ICTIndexRank:  apiPlayer.ICTIndexRank,
		},
		// ChanceOfPlaying:  chanceOfPlaying,
		PickedPercentage: float32(pickedPercentage),
	}

	history, err := requestPlayerHistory(int(newPlayer.ID))
	if err != nil {
		return models.Player{}, err
	}
	newPlayer.History = history

	return newPlayer, nil
}

func requestPlayerHistory(apiPlayerID int) (map[models.FixtureID]models.PlayerFixture, error) {
	fixturesAndHistoryApiBody, err := getJsonBody(fmt.Sprintf("%s/%d", playerFixturesApi, apiPlayerID))
	if err != nil {
		return nil, err
	}
	var fixturesAndHistory apiPlayerFixturesAndHistory
	if err := json.Unmarshal(fixturesAndHistoryApiBody, &fixturesAndHistory); err != nil {
		return nil, err
	}
	fixturesToPlayerFixtures := make(map[models.FixtureID]models.PlayerFixture, 0)
	for _, fixture := range fixturesAndHistory.History {
		fixturesToPlayerFixtures[models.FixtureID(fixture.FixtureID)] = models.PlayerFixture{
			FixtureID:   models.FixtureID(fixture.FixtureID),
			PlayerID:    models.PlayerID(fixture.ElementID),
			Minutes:     fixture.Minutes,
			Played:      fixture.Minutes > 0,
			Points:      fixture.TotalPoints,
			GoalsScored: fixture.GoalsScored,
			Assists:     fixture.Assists,
			YellowCards: fixture.YellowCards,
			RedCards:    fixture.RedCards,
			Bonus:       fixture.Bonus,
			WasHome:     fixture.WasHome,
		}
	}

	return fixturesToPlayerFixtures, nil
}

func backoff(f func() error, retries int, baseInterval time.Duration) error {
	var err error
	var fib1, fib2 int = 0, 1
	for i := 0; i <= retries; i++ {
		err = f()
		if err == nil {
			return nil
		}
		nextInterval := baseInterval * time.Duration(fib1)
		fib1, fib2 = fib2, fib1+fib2
		time.Sleep(nextInterval)
	}
	return err
}

func getJsonBody(endpoint string) ([]byte, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	// to circumvent max retries errors
	if resp.StatusCode != http.StatusOK {
		backoff(func() error {
			resp, err = http.Get(endpoint)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("status was not ok (status: %d)", resp.StatusCode)
			}
			return nil
		}, 10, time.Second*1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return body, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
