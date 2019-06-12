package main

import (
	"IMessageSearch/Search"
	"fmt"
)

const db string = "Users/stephane/Documents/chat.db"
func main()  {

	query := Search.Query{Db: ""}
	messages := query.Search("rate")
	last := messages[0]
	query.Enrich(&last)
	fmt.Println(last)
	}


