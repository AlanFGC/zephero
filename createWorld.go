package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	world "zephero/shared"
	utils "zephero/utils"
)

const WORLD_ID_ENV string = "WORLD_ID"

// this function is created with the intent to be run when a new world has to be created.
func main() {
	rows := *flag.Int("r", 100, "Optional: Number of rows (default: 10)")
	cols := *flag.Int("c", 100, "Optional: Number of columns (default: 10)")
	chunkLen := *flag.Int("len", 32, "Optional: Chunk length (default: 100)")

	dao := world.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		_ = errors.New("failed to open DB")
		return
	}
	id, err := dao.InsertNewWorld(rows, cols, chunkLen)
	fmt.Println("new world ID: ", id)
	if err != nil {
		_ = errors.New("failed to insert new world")
	}
	err = os.Setenv(WORLD_ID_ENV, strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatal("Failed to set environment variable.")
	} else {
		fmt.Println("WORLD_ID set to " + os.Getenv(WORLD_ID_ENV))
	}
	w, err := world.NewChunkedWorld(rows, cols, chunkLen)
	if err != nil || w == nil {
		_ = errors.New("failed to create new world")
	}
	setRandomUUIDs(w)
	err = w.Save(dao)
	if err != nil {
		_ = errors.New("failed to save world")
	}
	fmt.Println("World was created successfully")
	err = dao.CloseDb()
	if err != nil {
		_ = errors.New("failed to gracefully close DB")
		return
	}
}

func setRandomUUIDs(w world.World) {
	rows, cols := w.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if utils.Chance(0.30) {
				id := utils.GenerateTimeBasedID()
				err := w.SetSpace(id, 0, i, j)
				if err != nil {
					_ = errors.New("failed to set random ids")
				}
			}
		}
	}
}
