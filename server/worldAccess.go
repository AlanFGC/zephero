package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	worldRepo "zephero/database/sqlite_world_repo"
	gameWorld "zephero/world"
)

type WorldAccess struct {
	World   *gameWorld.ChunkedWorld
	worldId int64
	lock    sync.RWMutex
}

type PlayerView struct {
	Chunk01 gameWorld.WorldChunk `json:"Chunk01"`
	Chunk02 gameWorld.WorldChunk `json:"Chunk02"`
	Chunk03 gameWorld.WorldChunk `json:"Chunk03"`
	Chunk04 gameWorld.WorldChunk `json:"Chunk04"`
	Chunk05 gameWorld.WorldChunk `json:"Chunk05"`
	Chunk06 gameWorld.WorldChunk `json:"Chunk06"`
	Chunk07 gameWorld.WorldChunk `json:"Chunk07"`
	Chunk08 gameWorld.WorldChunk `json:"Chunk08"`
	Chunk09 gameWorld.WorldChunk `json:"Chunk09"`
}

// TODO it's probably a better idea for this function to recieve a new world, instead of instantiating it
func (wa *WorldAccess) Preload(ctx context.Context, path string, worldId int) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	wa.worldId = int64(worldId)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to load from sql")
		}
	}()

	worldQueries := worldRepo.New(db)
	sqlWorld, err := worldQueries.GetWorld(ctx, int64(worldId))
	if err != nil {
		return err
	}

	world, err := gameWorld.NewChunkedWorld(int(sqlWorld.RowLength),
		int(sqlWorld.ColumnLength),
		int(sqlWorld.ChunkLength))
	if err != nil || world == nil {
		return err
	}
	wa.World = world

	sqlChunks, err := worldQueries.GetWorldChunkByWorldId(ctx, int64(worldId))
	if err != nil {
		return err
	}

	if len(sqlChunks) != int(sqlWorld.RowLength*sqlWorld.ColumnLength) {
		return fmt.Errorf("world %d has %d chunks", worldId, len(sqlChunks))
	}

	for i := 0; i < len(sqlChunks); i++ {
		chunk, err := gameWorld.DeserializeChunkData(sqlChunks[i].Chunk)
		if err != nil {
			return err
		}
		err = world.SetChunk(int(sqlChunks[i].RowID), int(sqlChunks[i].ColID), chunk)
		if err != nil {
			return err
		}
	}
	log.Println(fmt.Sprint("Successfully preloaded world ", worldId))
	return nil
}

func (wa *WorldAccess) Save(ctx context.Context, path string) error {
	if wa.World == nil {
		return fmt.Errorf("Error: world is nil")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	worldQueries := worldRepo.New(db)
	rows, cols := wa.World.GetSize()
	err = worldQueries.UpdateWorld(ctx, worldRepo.UpdateWorldParams{
		RowLength:    int64(rows),
		ColumnLength: int64(cols),
		ChunkLength:  int64(wa.World.ChunkSize),
		WorldID:      wa.worldId,
	})
	if err != nil {
		return err
	}

	chunkData, err := wa.World.GetChunkData()
	if err != nil {
		return err
	}

	for i := 0; i < len(chunkData); i++ {
		for j := 0; j < len(chunkData[0]); j++ {
			chunk, err := gameWorld.SerializeChunkData(&chunkData[i][j])
			if err != nil {
				return err
			}
			err = worldQueries.UpdateWorldChunk(ctx, worldRepo.UpdateWorldChunkParams{
				Locked:  false,
				WorldID: wa.worldId,
				Chunk:   chunk,
				RowID:   int64(chunkData[i][j].Row),
				ColID:   int64(chunkData[i][j].Col),
			})
			if err != nil {
				log.Fatal("Failed to update chunk")
				return err
			}
		}
	}
	fmt.Println(fmt.Sprint("Successfully saved world ", wa.worldId))
	return nil
}

func (wa *WorldAccess) write(id uint64, terrainId uint64, row int, col int) {
	wa.lock.Lock()
	err := wa.World.SetSpace(id, terrainId, row, col)
	if err != nil {
		log.Fatal(err)
	}
	wa.lock.Unlock()
}

func (wa *WorldAccess) playerView(row int, col int) PlayerView {
	wa.lock.RLock()
	chunk, err := wa.World.GetPlayerViewByCellCoordinate(row, col)
	if err != nil {
		log.Fatal("Failed to load player view")
	}
	playerResponse := PlayerView{
		Chunk01: chunk[0],
		Chunk02: chunk[1],
		Chunk03: chunk[2],
		Chunk04: chunk[3],
		Chunk05: chunk[4],
		Chunk06: chunk[5],
		Chunk07: chunk[6],
		Chunk08: chunk[7],
		Chunk09: chunk[8],
	}
	wa.lock.RUnlock()
	return playerResponse
}
