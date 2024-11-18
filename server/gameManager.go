package server

import (
	"fmt"
	"zephero/shared"
)

type GameManager struct {
	events chan []PlayerEvent
	world  shared.ChunkedWorld
}

func NewGameManager(eventBatchSize int) *GameManager {
	return &GameManager{
		events: make(chan []PlayerEvent, eventBatchSize),
	}
}

func (game *GameManager) Configure(world *shared.ChunkedWorld) {
	game.world = *world
}

func (game *GameManager) Run() {
	for {
		eventBatch, ok := <-game.events
		fmt.Println("recieving data from event batch")
		if !ok {
			fmt.Println("Game event channel closed")
			return
		}
		for _, event := range eventBatch {
			game.processEvent(&event)
		}
	}
}

func (game *GameManager) SendEvent(event PlayerEvent) {
	game.events <- []PlayerEvent{event}
}

func (game *GameManager) processEvent(event *PlayerEvent) string {
	var buffer string
	switch event.GameEvent.EventId {
	case E_SPAWN:
		buffer = "spawn"
	case E_MOVE:
		buffer = "move"
	case E_DESPAWN:
		buffer = "despawn"
	default:
		buffer = "unknown"
	}

	return buffer
}
