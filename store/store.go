package store

import (
	"better-fantasy/api"
	"better-fantasy/models"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName = "./data/data.sqlite"
)

type WriteData interface {
	MarkImported(gameweekID int) error
	StorePlayer(player models.Player) error
	StorePlayerType(playerType models.PlayerType) error
	StoreTeam(team models.Team) error
	StoreGameweek(gameweek models.Gameweek) error
	StoreFixture(fixture models.Fixture) error
}

type ReadData interface {
	GetPlayer(playerID models.PlayerID) (models.Player, error)
}

func NewStore() DataStore {
	store := DataStore{}
	store.Setup()
	return store
}

type DataStore struct {
	Connection *sql.DB
}

func (p *DataStore) Connect() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	p.Connection = conn
	return conn, nil
}

func (p *DataStore) Close() error {
	return p.Connection.Close()
}

func (p *DataStore) Nuke() error {
	_, err := p.Connection.Exec(`
		DROP TABLE IF EXISTS players;
		DROP TABLE IF EXISTS player_types;
		DROP TABLE IF EXISTS teams;
		DROP TABLE IF EXISTS gameweeks;
		DROP TABLE IF EXISTS fixtures;
		DROP TABLE IF EXISTS imports;
	`)
	if err != nil {
		return err
	}
	return nil
}

func (p *DataStore) Dump() error {
	// this pattern isn't working
	_, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	exportDir := fmt.Sprintf("./exports/gw_%d", p.CurrentGameweek())
	err = os.Mkdir(exportDir, os.ModePerm)
	if err != nil {
		// end silently if dir already exists
		return nil
	}

	// Get a list of tables in the database
	tables, err := p.getTableNames()
	if err != nil {
		return err
	}

	// Dump each table to a separate SQL file
	for _, table := range tables {
		if err = p.dumpTableToFile(table, exportDir); err != nil {
			return err
		}
	}

	fmt.Println()

	return nil
}

func (p *DataStore) Setup() error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS imports (
		gameweek_id INT PRIMARY KEY,
		imported BOOLEAN NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS player_types (
		id INTEGER PRIMARY KEY,
		name TEXT,
		plural_name TEXT,
		short_name TEXT,
		team_player_count INTEGER,
		team_min_play_count INTEGER,
		team_max_play_count INTEGER
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY,
		name TEXT,
		form REAL,
		points_per_game REAL,
		total_points INTEGER,
		cost TEXT,
		raw_cost REAL,
		team_id INTEGER,
		type_id INTEGER,
		minutes INTEGER,
		goals INTEGER,
		assists INTEGER,
		conceded INTEGER,
		clean_sheets INTEGER,
		yellow_cards INTEGER,
		red_cards INTEGER,
		bonus INTEGER,
		starts INTEGER,
		average_starts REAL,
		matches_played REAL,
		ict_index REAL,
		ict_index_rank INTEGER,
		most_captained BOOLEAN,
		picked_percentage REAL
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY,
		name VARCHAR,
		short_name VARCHAR
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS gameweeks (
		id INT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		deadline DATETIME NOT NULL,
		is_current BOOLEAN NOT NULL,
		is_next BOOLEAN NOT NULL,
		finished BOOLEAN NOT NULL,
		most_captained_id INT
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS fixtures (
		id INT PRIMARY KEY,
		gameweek_id INT,
		home_team_id INT NOT NULL,
		away_team_id INT NOT NULL,
		home_team_difficulty INT NOT NULL,
		away_team_difficulty INT NOT NULL,
		difficulty_majority INT NOT NULL,
		CONSTRAINT fk_gameweek FOREIGN KEY (gameweek_id) REFERENCES gameweeks(id),
		CONSTRAINT fk_home_team FOREIGN KEY (home_team_id) REFERENCES teams(id),
		CONSTRAINT fk_away_team FOREIGN KEY (away_team_id) REFERENCES teams(id)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS player_fixtures (
		fixture_id INT NOT NULL,
		player_id INT NOT NULL,
		minutes INT NOT NULL,
		played BOOLEAN NOT NULL,
		points INT NOT NULL,
		goals_scored INT NOT NULL,
		assists INT NOT NULL,
		yellow_cards INT NOT NULL,
		red_cards INT NOT NULL,
		bonus INT NOT NULL,
		clean_sheet BOOLEAN NOT NULL,
		was_home BOOLEAN NOT NULL,
		PRIMARY KEY (fixture_id, player_id)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE table IF NOT EXISTS manager_picks (
		manager_id INT NOT NULL,
		gameweek_id INT NOT NULL,
		player_id INT NOT NULL,
		is_captain BOOLEAN NOT NULL,
		is_vice_captain BOOLEAN NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) CurrentGameweek() int {
	db, err := p.Connect()
	if err != nil {
		return 0
	}
	defer p.Close()

	row := db.QueryRow("SELECT `id` FROM `gameweeks` WHERE `is_current` = 1")

	var gameweek models.Gameweek
	err = row.Scan(
		&gameweek.ID,
	)

	if err != nil {
		return 0
	}

	return int(gameweek.ID)
}

func (p *DataStore) StoreData(data *api.Data, dumpData bool) error {
	// ensures dump only contains data for specific gw
	if dumpData {
		p.Nuke()
	}

	for _, playerType := range data.PlayerTypes {
		if err := p.StorePlayerType(playerType); err != nil {
			return err
		}
	}

	for _, team := range data.Teams {
		if err := p.StoreTeam(*team); err != nil {
			return err
		}
	}

	for _, gameweek := range data.Gameweeks {
		if err := p.StoreGameweek(gameweek); err != nil {
			return err
		}
	}

	for _, fixture := range data.Fixtures {
		// ensures we see only fixtures for specific gameweek
		if dumpData && (fixture.Gameweek.ID != models.GameweekID(p.CurrentGameweek())) {
			continue
		}

		if err := p.StoreFixture(*fixture); err != nil {
			return err
		}
	}

	for _, player := range data.Players {
		if err := p.StorePlayer(player); err != nil {
			return err
		}
		for _, fixture := range player.History {
			if err := p.StorePlayerFixture(fixture); err != nil {
				return err
			}
		}
	}

	for _, pick := range data.ManagerPicks {
		if err := p.StorePick(pick); err != nil {
			return err
		}
	}

	if dumpData {
		if err := p.Dump(); err != nil {
			return err
		}
	}

	p.MarkImported(p.CurrentGameweek())

	return nil
}

func (p *DataStore) StorePlayer(player models.Player) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO players (id, name, form, points_per_game, total_points, cost, raw_cost, team_id, type_id, minutes, goals, assists, conceded, clean_sheets, yellow_cards, red_cards, bonus, starts, average_starts, matches_played, ict_index, ict_index_rank, most_captained, picked_percentage)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, player.ID, player.Name, player.Form, player.PointsPerGame, player.TotalPoints, player.Cost, player.RawCost, player.Team.ID, player.Type.ID, player.Stats.Minutes, player.Stats.Goals, player.Stats.Assists, player.Stats.Conceded, player.Stats.CleanSheets, player.Stats.YellowCards, player.Stats.RedCards, player.Stats.Bonus, player.Stats.Starts, player.Stats.AverageStarts, player.Stats.MatchesPlayed, player.Stats.ICTIndex, player.Stats.ICTIndexRank, player.MostCaptained, player.PickedPercentage)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StorePlayerFixture(fixture models.PlayerFixture) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO player_fixtures (
			fixture_id, 
			player_id, 
			minutes, 
			played, 
			points, 
			goals_scored, 
			assists, 
			yellow_cards, 
			red_cards, 
			bonus, 
			clean_sheet,
			was_home
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err = db.Exec(
		query,
		fixture.FixtureID,
		fixture.PlayerID,
		fixture.Minutes,
		fixture.Played,
		fixture.Points,
		fixture.GoalsScored,
		fixture.Assists,
		fixture.YellowCards,
		fixture.RedCards,
		fixture.Bonus,
		fixture.CleanSheet,
		fixture.WasHome,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StorePlayerType(playerType models.PlayerType) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO player_types (id, name, plural_name, short_name, team_player_count, team_min_play_count, team_max_play_count)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, playerType.ID, playerType.Name, playerType.PluralName, playerType.ShortName, playerType.TeamPlayerCount, playerType.TeamMinPlayCount, playerType.TeamMaxPlayCount)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StoreTeam(team models.Team) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO teams (id, name, short_name)
		VALUES (?, ?, ?)
	`

	_, err = db.Exec(query, team.ID, team.Name, team.ShortName)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StoreGameweek(gameweek models.Gameweek) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO gameweeks (id, name, deadline, is_current, is_next, finished, most_captained_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, gameweek.ID, gameweek.Name, gameweek.Deadline, gameweek.IsCurrent, gameweek.IsNext, gameweek.Finished, gameweek.MostCaptainedID)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StoreFixture(fixture models.Fixture) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO fixtures (
		id,
		gameweek_id,
		home_team_id,
		away_team_id,
		home_team_difficulty,
		away_team_difficulty,
		difficulty_majority
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, fixture.ID, fixture.Gameweek.ID, fixture.HomeTeam.ID, fixture.AwayTeam.ID, fixture.HomeTeamDifficulty, fixture.AwayTeamDifficulty, fixture.DifficultyMajority)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) StorePick(pick models.ManagerPick) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO manager_picks (
			manager_id,
			gameweek_id,
			player_id,
			is_captain,
			is_vice_captain
		) VALUES (?, ?, ?, ?, ?)`

	_, err = db.Exec(query, pick.ManagerID, pick.GameweekID, pick.PlayerID, pick.IsCaptain, pick.IsViceCaptain)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) HasImported() (bool, error) {
	db, err := p.Connect()
	if err != nil {
		return false, err
	}
	defer p.Close()

	query := `SELECT imported FROM imports WHERE gameweek_id = ?`

	row := db.QueryRow(query, p.CurrentGameweek())

	var rowStruct struct {
		Imported bool
	}
	err = row.Scan(
		&rowStruct.Imported,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return rowStruct.Imported, nil
}

func (p *DataStore) MarkImported(gameweekID int) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO imports (
			gameweek_id,
			imported
		) VALUES (?, ?)
	`

	_, err = db.Exec(query, gameweekID, true)

	if err != nil {
		return err
	}

	return nil
}

func (p *DataStore) GetPlayer(playerID models.PlayerID) (models.Player, error) {
	db, err := p.Connect()
	if err != nil {
		return models.Player{}, err
	}
	defer p.Close()

	row := db.QueryRow(fmt.Sprintf("SELECT * FROM `players` WHERE `id` = %d", playerID))

	var player models.Player
	err = row.Scan(
		&player.ID,
		&player.Name,
	)

	if err != nil {
		return models.Player{}, err
	}

	return player, nil
}

func (p *DataStore) GetPlayers() (map[models.PlayerID]models.Player, error) {
	db, err := p.Connect()
	if err != nil {
		return nil, err
	}
	defer p.Close()

	playerRows, err := db.Query("SELECT id, name, form, points_per_game, total_points, cost, raw_cost, team_id, type_id, minutes, goals, assists, conceded, clean_sheets, yellow_cards, red_cards, bonus, starts, average_starts, picked_percentage FROM `players`")
	if err != nil {
		return nil, err
	}
	defer playerRows.Close()

	players := make(map[models.PlayerID]models.Player, 0)
	for playerRows.Next() {
		type player struct {
			ID               int
			Name             string
			Form             float32
			PointsPerGame    float32
			TotalPoints      int
			Cost             string
			RawCost          float32
			TeamID           int
			TypeID           int
			Minutes          int
			Goals            int
			Assists          int
			Conceded         int
			CleanSheets      int
			RedCards         int
			YellowCards      int
			Bonus            int
			Starts           int
			AverageStarts    float32
			PickedPercentage float32
		}
		var playerRow player
		err := playerRows.Scan(
			&playerRow.ID,
			&playerRow.Name,
			&playerRow.Form,
			&playerRow.PointsPerGame,
			&playerRow.TotalPoints,
			&playerRow.Cost,
			&playerRow.RawCost,
			&playerRow.TeamID,
			&playerRow.TypeID,
			&playerRow.Minutes,
			&playerRow.Goals,
			&playerRow.Assists,
			&playerRow.Conceded,
			&playerRow.CleanSheets,
			&playerRow.RedCards,
			&playerRow.YellowCards,
			&playerRow.Bonus,
			&playerRow.Starts,
			&playerRow.AverageStarts,
			&playerRow.PickedPercentage,
		)
		if err != nil {
			return nil, err
		}
		players[models.PlayerID(playerRow.ID)] = models.Player{
			ID:   models.PlayerID(playerRow.ID),
			Name: playerRow.Name,
			Form: playerRow.Form,
			Type: models.PlayerType{
				ID: models.PlayerTypeID(playerRow.TypeID),
			},
			PointsPerGame: playerRow.PointsPerGame,
			TotalPoints:   playerRow.TotalPoints,
			Cost:          playerRow.Cost,
			RawCost:       playerRow.RawCost,
			Stats: models.PlayerStats{
				Bonus: playerRow.Bonus,
			},
		}
	}

	if err = playerRows.Err(); err != nil {
		return nil, err
	}

	playerFixtureRows, err := db.Query("SELECT fixture_id, player_id, minutes, goals_scored, assists, clean_sheet FROM `player_fixtures`")
	if err != nil {
		return nil, err
	}
	defer playerFixtureRows.Close()

	for playerFixtureRows.Next() {
		type playerFixture struct {
			FixtureID   int
			PlayerID    int
			Minutes     int
			GoalsScored int
			Assists     int
			CleanSheet  bool
		}
		var playerFixtureRow playerFixture
		playerFixtureRows.Scan(
			&playerFixtureRow.FixtureID,
			&playerFixtureRow.PlayerID,
			&playerFixtureRow.Minutes,
			&playerFixtureRow.GoalsScored,
			&playerFixtureRow.Assists,
			&playerFixtureRow.CleanSheet,
		)
		if player, ok := players[models.PlayerID(playerFixtureRow.PlayerID)]; ok {
			if player.History == nil {
				player.History = make(map[models.FixtureID]models.PlayerFixture, 0)
			}
			player.History[models.FixtureID(playerFixtureRow.FixtureID)] = models.PlayerFixture{
				FixtureID:   models.FixtureID(playerFixtureRow.FixtureID),
				PlayerID:    models.PlayerID(playerFixtureRow.PlayerID),
				Minutes:     playerFixtureRow.Minutes,
				GoalsScored: playerFixtureRow.GoalsScored,
				Assists:     playerFixtureRow.Assists,
				CleanSheet:  playerFixtureRow.CleanSheet,
			}
			players[models.PlayerID(playerFixtureRow.PlayerID)] = player
		}
	}

	return players, nil
}

// getTableNames retrieves a list of table names from the SQLite database.
func (p *DataStore) getTableNames() ([]string, error) {
	rows, err := p.Connection.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// dumpTableToFile dumps the contents of the specified table to a SQL file.
func (p *DataStore) dumpTableToFile(tableName, tempDir string) error {
	outputFile := filepath.Join(tempDir, tableName+".sql")

	// use the sqlite3 command-line tool to dump the table to an SQL file
	cmd := exec.Command("sqlite3", dbName, fmt.Sprintf(".dump %s", tableName))
	dumpOutput, err := cmd.Output()
	if err != nil {
		return err
	}

	// write the dump output to the SQL file
	err = os.WriteFile(outputFile, dumpOutput, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("table '%s' dumped to '%s'\n", tableName, outputFile)

	return nil
}
