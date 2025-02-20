package utils

import (
	"fmt"
	table "github.com/olekukonko/tablewriter"
	"math/rand"
	"os"
	"strconv"
	"time"
	game "zephero/world"
)

func PrintWorld(w game.World) {
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

func GenerateTimeBasedID() uint64 {
	now := time.Now()
	millis := now.UnixNano() / int64(time.Millisecond)
	return uint64(millis)
}

func Chance(probability float64) bool {
	return rand.Float64() < probability
}

func SetRandomIds(w game.World) error {
	rows, cols := w.GetSize()
	fmt.Println(rows, cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if Chance(0.30) {
				id := GenerateTimeBasedID()
				err := w.SetSpace(id, 0, i, j)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
