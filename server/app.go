package main

import (
	"fmt"
	table "github.com/olekukonko/tablewriter"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	game "zephero/shared"
)

var counter uint64

func generateTimeBasedID() uint64 {
	timestamp := uint64(time.Now().UnixNano())
	counterValue := atomic.AddUint64(&counter, 1)
	return (timestamp << 16) | (counterValue & 0xFFFF)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	w, err := game.NewChunkedWorld(5, 5, 16)
	if err != nil {
		fmt.Print("FAILED TO CREATE:", err)
		return
	}
	setRandomUUIDs(w)
	printWorld(w)
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
				strMatrix[i][j] = string(str[19])
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
