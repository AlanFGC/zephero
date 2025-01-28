package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"zephero/shared"
)

type GameManager struct {
	events        chan []PlayerEvent
	world         *shared.ChunkedWorld
	activePlayers map[string]PlayerState
	access        WorldAccess
}

type PlayerState struct {
	userName  string
	entityId  uint64
	lastEvent time.Time
}

func NewGameManager(eventBatchSize int) *GameManager {
	return &GameManager{
		events:        make(chan []PlayerEvent, eventBatchSize),
		activePlayers: make(map[string]PlayerState),
	}
}

func (game *GameManager) Configure(ctx context.Context, world *shared.ChunkedWorld, dbPath string, worldId int) error {
	if world != nil {
		game.world = world
		game.access.World = world
	} else if len(dbPath) > 0 {
		err := game.access.Preload(ctx, dbPath, worldId)
		if err != nil {
			log.Fatalf(err.Error())
			return err
		}
		game.world = game.access.World
	} else {
		return errors.New(fmt.Sprintf("Invalid parameters for GameManager"))
	}
	return nil
}

func (game *GameManager) Run(ctx context.Context, dbPath string) {
	for {
		select {
		case eventBatch, ok := <-game.events:
			if !ok {
				log.Println("GameManager event channel closed")
				continue
			}

			for _, event := range eventBatch {
				game.processEvent(&event)
				if event.GameEvent.EventId == E_EXIT {
					log.Println("Shutting down...")
					err := game.access.Save(ctx, dbPath)
					if err != nil {
						log.Fatalf(err.Error())
					} else {
						return
					}
				}
			}
		default:
			// No events in the channel; proceed with other tasks
		}
		game.timeOutActivePlayers()
	}
}

func (game *GameManager) SendEvent(event PlayerEvent) {
	game.events <- []PlayerEvent{event}
}

func (game *GameManager) registerPlayer(event *PlayerEvent) error {
	log.Println("Registering new player: ", event.PlayerId)
	if game.activePlayers == nil {
		game.activePlayers = make(map[string]PlayerState)
	}
	_, exists := game.activePlayers[event.PlayerId]
	if !exists {
		game.activePlayers[event.PlayerId] = PlayerState{
			userName:  event.PlayerId,
			lastEvent: time.Now(),
		}
	}
	return nil
}

func (game *GameManager) unregisterPlayer(username string) error {
	if game.activePlayers == nil || len(game.activePlayers) == 0 {
		return errors.New("Failed to unregister player")
	}
	_, exists := game.activePlayers[username]
	if exists == false {
		log.Println("Warning: Failed to unregister player ", username)
		return nil
	}
	delete(game.activePlayers, username)
	return nil
}

const TIME_OUT_TIME = time.Second * 3

func (game *GameManager) timeOutActivePlayers() {
	currentTime := time.Now()
	for _, player := range game.activePlayers {
		currTime := currentTime.Sub(player.lastEvent)
		if currTime > TIME_OUT_TIME {
			log.Println("Removing player ", player.userName)
			err := game.unregisterPlayer(player.userName)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (game *GameManager) processEvent(event *PlayerEvent) string {
	username := event.PlayerId
	player, exists := game.activePlayers[username]
	if !exists {
		err := game.registerPlayer(event)
		if err != nil {
			return err.Error()
		}
	} else {
		player.lastEvent = time.Now()
		game.activePlayers[username] = player
	}
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
