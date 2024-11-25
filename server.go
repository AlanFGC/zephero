package main

import (
	"context"
	"fmt"
	"zephero/server"
)

func main() {
	worldId := createWorld(100, 100, 16, "world.db")
	if worldId == -1 {
		panic(fmt.Errorf("Failed to create and save a new world"))
	}
	gameManager := server.NewGameManager(100)
	ctx := context.Background()
	err := gameManager.Configure(ctx, nil, "database/sqliteDB/"+"world.db", 15)
	if err != nil {
		panic(err)
	}
	go server.RunWebSocketsServer(gameManager)
	gameManager.Run()
}
