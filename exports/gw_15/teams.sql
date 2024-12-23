PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE teams (
		id INTEGER PRIMARY KEY,
		name VARCHAR,
		short_name VARCHAR
	);
INSERT INTO teams VALUES(1,'Arsenal','ARS');
INSERT INTO teams VALUES(2,'Aston Villa','AVL');
INSERT INTO teams VALUES(3,'Bournemouth','BOU');
INSERT INTO teams VALUES(4,'Brentford','BRE');
INSERT INTO teams VALUES(5,'Brighton','BHA');
INSERT INTO teams VALUES(6,'Chelsea','CHE');
INSERT INTO teams VALUES(7,'Crystal Palace','CRY');
INSERT INTO teams VALUES(8,'Everton','EVE');
INSERT INTO teams VALUES(9,'Fulham','FUL');
INSERT INTO teams VALUES(10,'Ipswich','IPS');
INSERT INTO teams VALUES(11,'Leicester','LEI');
INSERT INTO teams VALUES(12,'Liverpool','LIV');
INSERT INTO teams VALUES(13,'Man City','MCI');
INSERT INTO teams VALUES(14,'Man Utd','MUN');
INSERT INTO teams VALUES(15,'Newcastle','NEW');
INSERT INTO teams VALUES(16,'Nott''m Forest','NFO');
INSERT INTO teams VALUES(17,'Southampton','SOU');
INSERT INTO teams VALUES(18,'Spurs','TOT');
INSERT INTO teams VALUES(19,'West Ham','WHU');
INSERT INTO teams VALUES(20,'Wolves','WOL');
COMMIT;
