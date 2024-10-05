package main

import (
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
	rows := flag.Int("r", 10, "Optional: Number of rows (default: 10)")
	cols := flag.Int("c", 10, "Optional: Number of columns (default: 10)")
	chunkLen := flag.Int("len", 32, "Optional: Chunk length (default: 100)")
	flag.Parse()

	// Initialize the DAO
	dao := world.NewSqliteDAO("world")
	err := dao.OpenDb("world.db")
	if err != nil {
		log.Fatalf("Error: failed to open DB: %v", err)
		return
	}
	defer func() {
		if err := dao.CloseDb(); err != nil {
			log.Printf("Warning: failed to close DB: %v", err)
		}
	}()

	// Insert new world
	id, err := dao.InsertNewWorld(*rows, *cols, *chunkLen)
	if err != nil {
		log.Fatalf("Error: failed to insert new world: %v", err)
		return
	}
	fmt.Println("New world ID:", id)

	// Set environment variable
	err = os.Setenv(WORLD_ID_ENV, strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatalf("Error: failed to set environment variable: %v", err)
		return
	}
	fmt.Println("WORLD_ID set to", os.Getenv(WORLD_ID_ENV))

	// Create a new chunked world
	w, err := world.NewChunkedWorld(*rows, *cols, *chunkLen)
	if err != nil {
		log.Fatalf("Error: failed to create new world: %v", err)
		return
	}
	if w == nil {
		log.Fatal("Error: null reference to world")
		return
	}

	// Set random UUIDs
	err = setRandomUUIDs(w)
	if err != nil {
		log.Printf("Warning: failed to set random UUIDs: %v", err)
		return
	}

	// Save the world
	err = w.SaveWorld(dao)
	if err != nil {
		log.Fatalf("Error: failed to save world: %v", err)
		return
	}
	fmt.Println("World was created successfully")
}

func setRandomUUIDs(w world.World) error {
	rows, cols := w.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if utils.Chance(0.30) {
				id := utils.GenerateTimeBasedID()
				err := w.SetSpace(id, 0, i, j)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
