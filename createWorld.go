package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	worldRepo "zephero/database/sqlite_world_repo"
	world "zephero/shared"
	"zephero/utils"
)

const WORLD_ID_ENV string = "WORLD_ID"
const DEFAULT_DB_NAME string = "world.db"
const PATH = "database/sqliteDB/"

// this function is created with the intent to be run when a new world has to be created.
func createWorld(rows int, cols int, chunkLen int, defaultDbName string) {
	ctx := context.Background()
	path := PATH + DEFAULT_DB_NAME
	if len(defaultDbName) > 0 {
		path = PATH + defaultDbName
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error: failed to open DB: %v", err)
		return
	}
	worldQueries := worldRepo.New(db)
	defer db.Close()

	// Insert new world
	id, err := worldQueries.InsertWorld(ctx, worldRepo.InsertWorldParams{
		RowLength:    int64(rows),
		ColumnLength: int64(cols),
		ChunkLength:  int64(chunkLen),
	})
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
	w, err := world.NewChunkedWorld(rows, cols, chunkLen)
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

	chunks, err := w.GetChunkData()
	if err != nil {
		log.Fatalf("Error: failed to get chunk data: %v", err)
	}
	for i := 0; i < int(rows); i++ {
		for j := 0; j < int(cols); j++ {
			chunk := chunks[i][j]
			binaryData, err := world.SerializeChunkData(&chunk)
			if err != nil {
				log.Fatalf("Error: failed to serialize chunk data: %v", err)
			}
			chunkId, err := worldQueries.InsertWorldChunk(ctx, worldRepo.InsertWorldChunkParams{
				WorldID: id,
				RowID:   int64(chunk.Row),
				ColID:   int64(chunk.Col),
				Locked:  false,
				Chunk:   binaryData,
			})
			if err != nil {
				log.Fatalf("Error: failed to insert new chunk: %v", err)
			}
			log.Printf("New chunk Saved to sql database: %d", chunkId)
		}
	}
}

func setRandomUUIDs(w world.World) error {
	rows, cols := w.GetSize()
	fmt.Println(rows, cols)
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
