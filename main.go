package main

import (
	"IMessageSearch/Search"
	"fmt"
)

func main()  {

	query := Search.Query{Db: "/Users/stephane/Documents/chat.db"}
	messages := query.Search("rate")
	last := messages[0]
	query.Enrich(&last)
	fmt.Println(last)
	}


