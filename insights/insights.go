package insights

import (
	"better-fantasy/models"
	"better-fantasy/printer"
	"better-fantasy/store"
	"fmt"
	"sort"
)

type Insights struct {
	Gameweek int
	Store    *store.DataStore
}

func NewInsights(gameweek int, store *store.DataStore) *Insights {
	return &Insights{
		Gameweek: gameweek,
		Store:    store,
	}
}

// the purpose of this function is to:
// - identify good value players (commanding higher points than their value suggests) and within this:
// - players in good form, exc. penalty goals (4 week average)
// - players on the rise (2 week average)
// - players with easy upcoming fixtures
// - players with the most bonus
func (i *Insights) Analyse() error {
	playersForGameweek, err := i.Store.PlayersForGameweek(i.Gameweek)
	if err != nil {
		return err
	}
	bestFormValuePlayers := sortPlayersByFormValueDesc(playersForGameweek)[:10]
	printPlayerList("The best value players (form):", bestFormValuePlayers)
	fmt.Println()
	bestPointsValuePlayers := sortPlayersByPointsValueDesc(playersForGameweek)[:10]
	printPlayerList("The best value players (total):", bestPointsValuePlayers)
	return nil
}

func sortPlayersByFormValueDesc(players []models.Player) []models.Player {
	sort.Slice(players, func(i, j int) bool {
		return players[i].FormOverCost() > players[j].FormOverCost()
	})
	return players
}

func sortPlayersByPointsValueDesc(players []models.Player) []models.Player {
	sort.Slice(players, func(i, j int) bool {
		return players[i].PointsOverCost() > players[j].PointsOverCost()
	})
	return players
}

func printPlayerList(title string, players []models.Player) {
	list := printer.List{
		Title: title,
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range players {
		list.Items = append(list.Items, player.ToListItem())
	}
	printer.PrintList(list)
}
