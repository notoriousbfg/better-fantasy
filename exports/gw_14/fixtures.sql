PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE fixtures (
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
	);
COMMIT;
