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

func (i *Insights) Analyse() error {
	playersForGameweek, err := i.Store.GetPlayers()
	if err != nil {
		return err
	}
	playersSlice := make([]models.Player, 0)
	for _, player := range playersForGameweek {
		playersSlice = append(playersSlice, player)
	}
	highestRankedDefenders := highestRankedDefenders(playersSlice)
	highestDefenders := printer.List{
		Title: "Top 20 defenders:",
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range highestRankedDefenders {
		highestDefenders.Items = append(highestDefenders.Items, printer.ListItem{
			Format: "%s (%s) (%d pts)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.TotalPoints,
			},
		})
	}
	printer.PrintList(highestDefenders)
	fmt.Println()
	mostAttackingDefenders := mostAttackingDefenders(playersSlice)
	attackingList := printer.List{
		Title: "Most attacking defenders (total points g/a):",
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range mostAttackingDefenders {
		attackingList.Items = append(attackingList.Items, printer.ListItem{
			Format: "%s (%s) (%.0f)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.AttackingPoints(),
			},
		})
	}
	printer.PrintList(attackingList)
	fmt.Println()
	defendersWithMostCleanSheetPoints := defendersWithMostCleanSheetPoints(playersSlice)
	cleanSheetList := printer.List{
		Title: "Defenders with most clean sheets:",
		Items: make([]printer.ListItem, 0),
	}
	for _, player := range defendersWithMostCleanSheetPoints {
		cleanSheetList.Items = append(cleanSheetList.Items, printer.ListItem{
			Format: "%s (%s) (%.0f)",
			Values: []interface{}{
				player.Name,
				player.Cost,
				player.CleanSheets(),
			},
		})
	}
	printer.PrintList(cleanSheetList)
	fmt.Println()
	// bestValueDefenders := bestValueDefenders(playersSlice)
	// valueList := printer.List{
	// 	Title: "Best value defenders:",
	// 	Items: make([]printer.ListItem, 0),
	// }
	// for _, player := range bestValueDefenders {
	// 	valueList.Items = append(valueList.Items, printer.ListItem{
	// 		Format: "%s (%s) (%.2f)",
	// 		Values: []interface{}{
	// 			player.Name,
	// 			player.Cost,
	// 			player.AttackingPoints(),
	// 		},
	// 	})
	// }
	// printer.PrintList(valueList)
	return nil
}

// func sortPlayersByFormValueDesc(players []models.Player) []models.Player {
// 	tmp := make([]models.Player, len(players))
// 	copy(tmp, players)
// 	sort.Slice(tmp, func(i, j int) bool {
// 		return tmp[i].FormOverCost() > tmp[j].FormOverCost()
// 	})
// 	return tmp
// }

// func sortPlayersByPointsValueDesc(players []models.Player) []models.Player {
// 	tmp := make([]models.Player, len(players))
// 	copy(tmp, players)
// 	sort.Slice(tmp, func(i, j int) bool {
// 		return tmp[i].PointsOverCost() > tmp[j].PointsOverCost()
// 	})
// 	return tmp
// }

// func sortPlayersByBonus(players []models.Player) []models.Player {
// 	tmp := make([]models.Player, len(players))
// 	copy(tmp, players)
// 	sort.Slice(tmp, func(i, j int) bool {
// 		return tmp[i].Stats.Bonus > tmp[j].Stats.Bonus
// 	})
// 	return tmp
// }

func highestRankedDefenders(players []models.Player) []models.Player {
	defenders := make([]models.Player, 0)
	for _, player := range players {
		if player.Type.ID == models.PTDefender {
			defenders = append(defenders, player)
		}
	}
	sort.Slice(defenders, func(i, j int) bool {
		return float32(defenders[i].TotalPoints) > float32(defenders[j].TotalPoints)
	})
	return defenders[:20]
}

func mostAttackingDefenders(players []models.Player) []models.Player {
	defenders := make([]models.Player, 0)
	for _, player := range players {
		if player.Type.ID == models.PTDefender {
			defenders = append(defenders, player)
		}
	}
	sort.Slice(defenders, func(i, j int) bool {
		return defenders[i].AttackingPoints() > defenders[j].AttackingPoints()
	})
	return defenders[:20]
}

func defendersWithMostCleanSheetPoints(players []models.Player) []models.Player {
	defenders := make([]models.Player, 0)
	for _, player := range players {
		if player.Type.ID == models.PTDefender {
			defenders = append(defenders, player)
		}
	}
	sort.Slice(defenders, func(i, j int) bool {
		return defenders[i].CleanSheets() > defenders[j].CleanSheets()
	})
	return defenders[:20]
}

// func bestValueDefenders(players []models.Player) []models.Player {
// 	defenders := make([]models.Player, 0)
// 	for _, player := range players {
// 		if player.Type.ID == models.PTDefender {
// 			defenders = append(defenders, player)
// 		}
// 	}
// 	sort.Slice(defenders, func(i, j int) bool {
// 		return defenders[i].PointsOverCost() > defenders[j].PointsOverCost()
// 	})
// 	return defenders[:20]
// }

// func printPlayerCostsList(title string, players []models.Player) {
// 	list := printer.List{
// 		Title: title,
// 		Items: make([]printer.ListItem, 0),
// 	}
// 	for _, player := range players {
// 		list.Items = append(list.Items, printer.ListItem{
// 			Format: "%s (%s) (%.2f)",
// 			Values: []interface{}{
// 				player.Name,
// 				player.Cost,
// 				player.AttackingForm(4),
// 			},
// 		})
// 	}
// 	printer.PrintList(list)
// }

// func printPlayerBonusList(title string, players []models.Player) {
// 	list := printer.List{
// 		Title: title,
// 		Items: make([]printer.ListItem, 0),
// 	}
// 	for _, player := range players {
// 		list.Items = append(list.Items, printer.ListItem{
// 			Format: "%s (%d)",
// 			Values: []interface{}{
// 				player.Name,
// 				player.Stats.Bonus,
// 			},
// 		})
// 	}
// 	printer.PrintList(list)
// }
