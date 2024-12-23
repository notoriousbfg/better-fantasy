package main

import (
	"better-fantasy/api"
	"better-fantasy/insights"
	"better-fantasy/store"
	"flag"
)

func main() {
	gameWeekInt := flag.Int("gameweek", 0, "for specifying the gameweek")
	dump := flag.Bool("dump", false, "for saving data to file")
	flag.Parse()

	store := store.NewStore(*gameWeekInt)

	hasImported, err := store.HasImported()
	if err != nil {
		panic(err)
	}

	if !hasImported {
		data, err := api.FetchData()
		if err != nil {
			panic(err)
		}

		err = store.StoreData(data, *dump)
		if err != nil {
			panic(err)
		}
	}

	insights := insights.NewInsights(*gameWeekInt, &store)
	err = insights.Analyse()
	if err != nil {
		panic(err)
	}
}
