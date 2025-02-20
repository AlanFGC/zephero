package main

import (
	"context"
	"fmt"
	"zephero/server"
	"zephero/world"
)

func main() {
	worldId := world.RunWorld(10, 10, 16, "world.db")
	if worldId == -1 {
		panic(fmt.Errorf("failed to create and save a new world"))
	}
	gameManager := server.NewGameManager(100)
	ctx := context.Background()
	dbPath := "database/sqliteDB/" + "world.db"
	err := gameManager.Configure(ctx, nil, dbPath, worldId)
	if err != nil {
		panic(err)
	}
	go server.RunWebSocketsServer(gameManager)
	gameManager.Run(ctx, dbPath)
}
