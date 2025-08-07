package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"zephero/core/world"
)

type GameManager struct {
	events        chan []PlayerEvent
	world         *world.ChunkedWorld
	activePlayers map[string]PlayerState
	access        WorldAccess
	tickCount     int
	lastTick      time.Time
}

type PlayerState struct {
	userName       string
	lastEvent      time.Time
	onUpdateCb     func(view *PlayerView)
	removePlayerCb func()
}

func NewGameManager(eventBatchSize int) *GameManager {
	return &GameManager{
		events:        make(chan []PlayerEvent, eventBatchSize),
		activePlayers: make(map[string]PlayerState),
	}
}

func (game *GameManager) Configure(ctx context.Context, world *world.ChunkedWorld, dbPath string, worldId int) error {
	if world != nil {
		game.world = world
		game.access.World = world
	} else if len(dbPath) > 0 {
		err := game.access.Preload(ctx, dbPath, worldId)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		game.world = game.access.World
	} else {
		return fmt.Errorf("invalid parameters for GameManager")
	}
	game.lastTick = time.Now()
	game.tickCount = 0
	return nil
}

func (game *GameManager) updatePlayerLastEvent(username string) {
	if player, exists := game.activePlayers[username]; exists {
		player.lastEvent = time.Now()
		game.activePlayers[username] = player
	}
}

func (game *GameManager) sendPlayerUpdates() {
	for _, state := range game.activePlayers {
		res := game.access.playerView(0, 0)
		state.onUpdateCb(&res)
	}
}

const TicksPerSecond = 60

func (game *GameManager) tick() {
	now := time.Now()
	if now.Sub(game.lastTick) >= time.Second {
		if game.tickCount < TicksPerSecond {
			fmt.Println("Warning failed to reach minimum tick rate", TicksPerSecond)
		}
		game.lastTick = now
		game.tickCount = 0
	}
	if game.tickCount < TicksPerSecond {
		// fmt.Println("Tick")
	}
	game.tickCount++
}

func (game *GameManager) Run() {
	for {
		select {
		case eventBatch, ok := <-game.events:
			if !ok {
				log.Println("GameManager event channel closed")
				continue
			}

			for _, event := range eventBatch {
				username := event.PlayerId
				game.updatePlayerLastEvent(username)
				err := game.processEvent(&event)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		default:
			// No events in the channel; proceed with other tasks
		}
		game.tick()
		game.sendPlayerUpdates()
		game.timeOutActivePlayers()
	}
}

func (game *GameManager) SendEvent(event PlayerEvent) {
	game.events <- []PlayerEvent{event}
}

func (game *GameManager) registerPlayer(username string, onUpdate func(view *PlayerView), onConnectionEnded func()) {
	if game.activePlayers == nil {
		game.activePlayers = make(map[string]PlayerState)
	}

	_, exists := game.activePlayers[username]
	if !exists {
		game.activePlayers[username] = PlayerState{
			userName:       username,
			lastEvent:      time.Now(),
			onUpdateCb:     onUpdate,
			removePlayerCb: onConnectionEnded,
		}
	}
}

func (game *GameManager) unregisterPlayer(username string) error {
	if len(game.activePlayers) == 0 {
		return errors.New("failed to unregister player")
	}
	playerState, exists := game.activePlayers[username]
	if !exists {
		log.Println("Warning: Failed to unregister player ", username)
		return nil
	}
	playerState.removePlayerCb()
	delete(game.activePlayers, username)
	return nil
}

const TimeOutTime = time.Second * 30

func (game *GameManager) timeOutActivePlayers() {
	currentTime := time.Now()
	for _, player := range game.activePlayers {
		currTime := currentTime.Sub(player.lastEvent)
		if currTime > TimeOutTime {
			log.Println("Removing player ", player.userName)
			err := game.unregisterPlayer(player.userName)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func (game *GameManager) processEvent(event *PlayerEvent) error {
	return nil
}
