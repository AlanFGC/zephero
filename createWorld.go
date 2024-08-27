package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	world "zephero/shared"
)


func main() {
	dao := world.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	dao.CloseDb()
}

func setRandomUUIDs(w world.chunkedWorld) {
	rows, cols := w.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Println('chunk')
		}
	}
}
