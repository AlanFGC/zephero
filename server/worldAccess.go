package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	worldRepo "zephero/database/sqlite_world_repo"
	gameWorld "zephero/shared"
)

type WorldAccess struct {
	World   *gameWorld.ChunkedWorld
	worldId int64
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
			log.Fatal(err)
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
