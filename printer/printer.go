package printer

import "fmt"

type Listable interface {
	ToListItem() ListItem
}

type List struct {
	Title string
	Items []ListItem
}

type ListItem struct {
	Format string
	Values []interface{}
}

func PrintList(list List) {
	if list.Title != "" {
		fmt.Println(list.Title)
	}
	for _, item := range list.Items {
		fmt.Printf(item.Format+"\n", item.Values...)
	}
}
