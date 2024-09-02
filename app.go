package main

import (
	"fmt"
	game "zephero/shared"
)

func main() {
	dao := game.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		fmt.Print(err)
	}
	dao.CloseDb()
}


