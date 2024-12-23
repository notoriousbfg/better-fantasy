package main

import (
	"better-fantasy/api"
	"better-fantasy/store"
	"flag"
)

func main() {
	gameWeekInt := flag.Int("gameweek", 0, "for specifying the gameweek")
	save := flag.Bool("save", false, "for storing data")
	dump := flag.Bool("dump", false, "for saving data to file")
	nuke := flag.Bool("nuke", false, "nuke the database")
	flag.Parse()

	data, err := api.FetchData()
	if err != nil {
		panic(err)
	}

	if *save {
		defer func() {
			err = store.StoreData(data, *gameWeekInt, *nuke, *dump)
			if err != nil {
				panic(err)
			}
		}()
	}
}
