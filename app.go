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

func generateTimeBasedID() uint64 {
	timestamp := uint64(time.Now().UnixNano())
	counterValue := atomic.AddUint64(&counter, 1)
	return (timestamp << 16) | (counterValue & 0xFFFF)
}

func main() {
	dao := game.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		fmt.Print(err)
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	w, err := game.NewChunkedWorld(200, 200, 16)
	if err != nil {
		fmt.Print("FAILED TO CREATE:", err)
		return
	}
	setRandomUUIDs(w)
	printWorld(w)
	w.Save(dao)
	err = dao.CloseDb()
	if err != nil {
		fmt.Print(err)
	}
}

func chance(probability float64) bool {
	return rand.Float64() < probability
}

func printWorld(w game.World) {
	rows, cols := w.GetSize()
	t := table.NewWriter(os.Stdout)
	strMatrix := make([][]string, rows)
	for i := 0; i < rows; i++ {
		strMatrix[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			node, _ := w.GetSpace(i, j)
			if node != nil {
				str := strconv.FormatUint(node.GetId(), 10)
				strMatrix[i][j] = string(str[0])
			} else {
				strMatrix[i][j] = " "
			}
		}
	}
	for _, row := range strMatrix {
		t.Append(row)
	}
	t.Render()
}

func setRandomUUIDs(w game.World) {
	rows, cols := w.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if chance(.25) {
				n := game.NewGNode(generateTimeBasedID())
				w.SetSpace(n, i, j)
			}
		}
	}
}
