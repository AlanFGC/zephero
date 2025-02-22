package world

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	worldRepo "zephero/database/sqlite_world_repo"

	_ "github.com/mattn/go-sqlite3"
)

const WORLD_ID_ENV string = "WORLD_ID"
const DEFAULT_DB_NAME string = "world.db"
const PATH = "database/sqliteDB/"
const FAILED = -1

// this function is created with the intent to be run when world has to be loaded from an sql DB
func RunWorld(rows int, cols int, chunkLen int, defaultDbName string) int {
	ctx := context.Background()
	path := PATH + DEFAULT_DB_NAME
	if len(defaultDbName) > 0 {
		path = PATH + defaultDbName
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error: failed to open DB: %v", err)
		return FAILED
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
		return FAILED
	}
	// Set environment variable
	err = os.Setenv(WORLD_ID_ENV, strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatalf("Error: failed to set environment variable: %v", err)
		return FAILED
	}
	fmt.Println("WORLD_ID set to", os.Getenv(WORLD_ID_ENV))

	// Create a new chunked world
	w, err := NewChunkedWorld(rows, cols, chunkLen)
	if err != nil {
		log.Fatalf("Error: failed to create new world: %v", err)
		return FAILED
	}
	if w == nil {
		log.Fatal("Error: null reference to world")
		return FAILED
	}
	chunks, err := w.GetChunkData()
	if err != nil {
		log.Fatalf("Error: failed to get chunk data: %v", err)
		return FAILED
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			chunk := chunks[i][j]
			binaryData, err := SerializeChunkData(&chunk)
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

	return int(id)
}
