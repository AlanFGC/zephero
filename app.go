package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	game "zephero/shared"

	table "github.com/olekukonko/tablewriter"
)

var counter uint64

func main() {
	dao := game.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		fmt.Print(err)
	}
	dao.CloseDb()
}


