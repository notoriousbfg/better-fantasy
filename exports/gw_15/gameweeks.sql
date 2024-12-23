PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE gameweeks (
		id INT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		deadline DATETIME NOT NULL,
		is_current BOOLEAN NOT NULL,
		is_next BOOLEAN NOT NULL,
		finished BOOLEAN NOT NULL,
		most_captained_id INT
	);
INSERT INTO gameweeks VALUES(1,'Gameweek 1','16 Aug 17:30',0,0,1,351);
INSERT INTO gameweeks VALUES(2,'Gameweek 2','24 Aug 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(3,'Gameweek 3','31 Aug 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(4,'Gameweek 4','14 Sep 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(5,'Gameweek 5','21 Sep 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(6,'Gameweek 6','28 Sep 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(7,'Gameweek 7','05 Oct 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(8,'Gameweek 8','19 Oct 10:00',0,0,1,351);
INSERT INTO gameweeks VALUES(9,'Gameweek 9','25 Oct 17:30',0,0,1,351);
INSERT INTO gameweeks VALUES(10,'Gameweek 10','02 Nov 11:00',0,0,1,351);
INSERT INTO gameweeks VALUES(11,'Gameweek 11','09 Nov 13:30',0,0,1,351);
INSERT INTO gameweeks VALUES(12,'Gameweek 12','23 Nov 11:00',0,0,1,351);
INSERT INTO gameweeks VALUES(13,'Gameweek 13','29 Nov 18:30',0,0,1,328);
INSERT INTO gameweeks VALUES(14,'Gameweek 14','03 Dec 18:00',0,0,1,351);
INSERT INTO gameweeks VALUES(15,'Gameweek 15','07 Dec 11:00',0,0,1,328);
INSERT INTO gameweeks VALUES(16,'Gameweek 16','14 Dec 13:30',0,0,1,328);
INSERT INTO gameweeks VALUES(17,'Gameweek 17','21 Dec 11:00',1,0,1,328);
INSERT INTO gameweeks VALUES(18,'Gameweek 18','26 Dec 11:00',0,1,0,0);
INSERT INTO gameweeks VALUES(19,'Gameweek 19','29 Dec 13:00',0,0,0,0);
INSERT INTO gameweeks VALUES(20,'Gameweek 20','04 Jan 11:00',0,0,0,0);
INSERT INTO gameweeks VALUES(21,'Gameweek 21','14 Jan 18:00',0,0,0,0);
INSERT INTO gameweeks VALUES(22,'Gameweek 22','18 Jan 11:00',0,0,0,0);
INSERT INTO gameweeks VALUES(23,'Gameweek 23','25 Jan 13:30',0,0,0,0);
INSERT INTO gameweeks VALUES(24,'Gameweek 24','01 Feb 11:00',0,0,0,0);
INSERT INTO gameweeks VALUES(25,'Gameweek 25','14 Feb 18:30',0,0,0,0);
INSERT INTO gameweeks VALUES(26,'Gameweek 26','21 Feb 18:30',0,0,0,0);
INSERT INTO gameweeks VALUES(27,'Gameweek 27','25 Feb 18:15',0,0,0,0);
INSERT INTO gameweeks VALUES(28,'Gameweek 28','08 Mar 13:30',0,0,0,0);
INSERT INTO gameweeks VALUES(29,'Gameweek 29','15 Mar 13:30',0,0,0,0);
INSERT INTO gameweeks VALUES(30,'Gameweek 30','01 Apr 17:15',0,0,0,0);
INSERT INTO gameweeks VALUES(31,'Gameweek 31','05 Apr 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(32,'Gameweek 32','12 Apr 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(33,'Gameweek 33','19 Apr 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(34,'Gameweek 34','26 Apr 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(35,'Gameweek 35','03 May 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(36,'Gameweek 36','10 May 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(37,'Gameweek 37','18 May 12:30',0,0,0,0);
INSERT INTO gameweeks VALUES(38,'Gameweek 38','25 May 13:30',0,0,0,0);
COMMIT;
