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

func NewInsights(store *store.DataStore) *Insights {
	return &Insights{
		Gameweek: store.CurrentGameweek(),
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
	playersForGameweek, err := i.Store.GetPlayers()
	if err != nil {
		return err
	}
	playersSlice := make([]models.Player, 0)
	for _, player := range playersForGameweek {
		playersSlice = append(playersSlice, player)
	}
	// bestFormValuePlayers := sortPlayersByFormValueDesc(playersSlice)[:10]
	// printPlayerCostsList("The best value players (form):", bestFormValuePlayers)
	// fmt.Println()
	// bestPointsValuePlayers := sortPlayersByPointsValueDesc(playersSlice)[:10]
	// printPlayerCostsList("The best value players (total):", bestPointsValuePlayers)
	// fmt.Println()
	// highestBonus := sortPlayersByBonus(playersSlice)[:10]
	// printPlayerBonusList("Players with the highest numbers of bonus points:", highestBonus)
	highestAttackingFormPlayersFourWeeks := sortPlayersByAttackingForm(playersSlice, 4)[:10]
	list := printer.List{
		Title: "The most attacking players (form over 4 weeks):",
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range highestAttackingFormPlayersFourWeeks {
		list.Items = append(list.Items, printer.ListItem{
			Format: "%s (%s) (%.2f)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.AttackingForm(4),
			},
		})
	}
	printer.PrintList(list)
	fmt.Println()
	highestAttackingFormPlayersTwoWeeks := sortPlayersByAttackingForm(playersSlice, 2)[:10]
	list = printer.List{
		Title: "Attacking players on the rise (form over 2 weeks):",
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range highestAttackingFormPlayersTwoWeeks {
		list.Items = append(list.Items, printer.ListItem{
			Format: "%s (%s) (%.2f)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.AttackingForm(2),
			},
		})
	}
	printer.PrintList(list)
	return nil
}

func sortPlayersByFormValueDesc(players []models.Player) []models.Player {
	tmp := make([]models.Player, len(players))
	copy(tmp, players)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].FormOverCost() > tmp[j].FormOverCost()
	})
	return tmp
}

func sortPlayersByPointsValueDesc(players []models.Player) []models.Player {
	tmp := make([]models.Player, len(players))
	copy(tmp, players)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].PointsOverCost() > tmp[j].PointsOverCost()
	})
	return tmp
}

func sortPlayersByBonus(players []models.Player) []models.Player {
	tmp := make([]models.Player, len(players))
	copy(tmp, players)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Stats.Bonus > tmp[j].Stats.Bonus
	})
	return tmp
}

func sortPlayersByAttackingForm(players []models.Player, weeks int) []models.Player {
	tmp := make([]models.Player, len(players))
	copy(tmp, players)
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].AttackingForm(weeks) > tmp[j].AttackingForm(weeks)
	})
	return tmp
}

func printPlayerCostsList(title string, players []models.Player) {
	list := printer.List{
		Title: title,
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range players {
		list.Items = append(list.Items, printer.ListItem{
			Format: "%s (%s) (%.2f)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.AttackingForm(4),
			},
		})
	}
	printer.PrintList(list)
}

func printPlayerBonusList(title string, players []models.Player) {
	list := printer.List{
		Title: title,
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range players {
		list.Items = append(list.Items, printer.ListItem{
			Format: "%s (%d)",
			Values: []interface{}{
				player.Name,
				player.Stats.Bonus,
			},
		})
	}
	printer.PrintList(list)
}
