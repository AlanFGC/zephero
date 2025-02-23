package main

import (
	"context"
	"fmt"
	"zephero/core"
	"zephero/core/world"
)

func main() {
	worldId := world.RunWorld(10, 10, 16, "world.db")
	if worldId == -1 {
		panic(fmt.Errorf("failed to create and save a new world"))
	}
	gameManager := core.NewGameManager(100)
	ctx := context.Background()
	dbPath := "database/sqliteDB/" + "world.db"
	err := gameManager.Configure(ctx, nil, dbPath, worldId)
	if err != nil {
		panic(err)
	}
	go core.RunWebSocketsServer(gameManager)
	gameManager.Run()
}
