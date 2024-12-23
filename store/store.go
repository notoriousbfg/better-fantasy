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
	StorePlayer(player models.Player) error
	StorePlayerType(playerType models.PlayerType) error
	StoreTeam(team models.Team) error
	StoreGameweek(gameweek models.Gameweek) error
	StoreFixture(fixture models.Fixture) error
}

type ReadData interface {
	GetPlayer(playerID models.PlayerID) (models.Player, error)
}

func StoreData(data *api.Data, gameweekInt int, nuke bool, dumpData bool) error {
	store := DataStore{
		GameweekID: gameweekInt,
	}
	store.Setup(nuke)

	for _, playerType := range data.PlayerTypes {
		if err := store.StorePlayerType(playerType); err != nil {
			return err
		}
	}

	for _, team := range data.Teams {
		if err := store.StoreTeam(*team); err != nil {
			return err
		}
	}

	for _, player := range data.Players {
		if err := store.StorePlayer(player); err != nil {
			return err
		}
	}

	for _, gameweek := range data.Gameweeks {
		if err := store.StoreGameweek(gameweek); err != nil {
			return err
		}
	}

	for _, fixture := range data.Fixtures {
		if fixture.Gameweek.ID != models.GameweekID(store.GameweekID) {
			continue
		}

		if err := store.StoreFixture(*fixture); err != nil {
			return err
		}
	}

	if dumpData {
		if err := store.Dump(); err != nil {
			return err
		}
	}

	return nil
}

type DataStore struct {
	GameweekID int
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

func (p *DataStore) Dump() error {
	// this pattern isn't working
	_, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	exportDir := fmt.Sprintf("./exports/gw_%d", p.GameweekID)
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

func (p *DataStore) Setup(nukeDatabase bool) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	if nukeDatabase {
		_, err = db.Exec(`
			DROP TABLE IF EXISTS players;
			DROP TABLE IF EXISTS player_types;
			DROP TABLE IF EXISTS teams;
			DROP TABLE IF EXISTS gameweeks;
			DROP TABLE IF EXISTS fixtures;
		`)
		if err != nil {
			return err
		}
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
		gameweek_player_id VARCHAR PRIMARY KEY,
		id INTEGER,
		gameweek_id INTEGER,
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

	return nil
}

func (p *DataStore) StorePlayer(player models.Player) error {
	db, err := p.Connect()
	if err != nil {
		return err
	}
	defer p.Close()

	query := `
		INSERT OR IGNORE INTO players (gameweek_player_id, id, gameweek_id, name, form, points_per_game, total_points, cost, raw_cost, team_id, type_id, minutes, goals, assists, conceded, clean_sheets, yellow_cards, red_cards, bonus, starts, average_starts, matches_played, ict_index, ict_index_rank, most_captained, picked_percentage)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, fmt.Sprintf("%d_%d", p.GameweekID, player.ID), player.ID, p.GameweekID, player.Name, player.Form, player.PointsPerGame, player.TotalPoints, player.Cost, player.RawCost, player.Team.ID, player.Type.ID, player.Stats.Minutes, player.Stats.Goals, player.Stats.Assists, player.Stats.Conceded, player.Stats.CleanSheets, player.Stats.YellowCards, player.Stats.RedCards, player.Stats.Bonus, player.Stats.Starts, player.Stats.AverageStarts, player.Stats.MatchesPlayed, player.Stats.ICTIndex, player.Stats.ICTIndexRank, player.MostCaptained, player.PickedPercentage)

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

	fmt.Printf("player: %+v", player)
	return models.Player{}, nil
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
