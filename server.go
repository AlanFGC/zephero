package main

import (
	"zephero/server"
)

func main() {
	//createWorld(10, 10, 32, "worldDB")
	gameManager := server.NewGameManager(100)
	go server.RunWebSocketsServer(gameManager)
	gameManager.Run()
}
