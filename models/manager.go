package models

type ManagerPick struct {
	ManagerID     int
	PlayerID      int
	GameweekID    GameweekID
	IsCaptain     bool
	IsViceCaptain bool
}
