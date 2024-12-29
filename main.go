package main

import (
	"better-fantasy/api"
	"better-fantasy/insights"
	"better-fantasy/store"
	"flag"
	"fmt"
)

func main() {
	managerID := flag.Int("manager", 0, "manager id")
	dump := flag.Bool("dump", false, "for saving data to file")
	flag.Parse()

	store := store.NewStore()

	hasImported, err := store.HasImported()
	if err != nil {
		panic(err)
	}

	if !hasImported {
		fmt.Println("This may take several minutes...")
		data, err := api.FetchData(api.FetchOptions{
			ManagerID: *managerID,
		})
		if err != nil {
			panic(err)
		}

		err = store.StoreData(data, *dump)
		if err != nil {
			panic(err)
		}
	}

	insights := insights.NewInsights(&store)
	err = insights.Analyse()
	if err != nil {
		panic(err)
	}
}
