package server

import "fmt"

type PlayerEvent struct {
	PlayerId string `json:"playerId"`
	EventId  string `json:"eventId"`
}

type GameManager struct {
	events chan []PlayerEvent
}

func NewGameManager(eventBatchSize int) *GameManager {
	return &GameManager{
		events: make(chan []PlayerEvent, eventBatchSize),
	}
}

func (game *GameManager) Configure() {

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

func (game *GameManager) processEvent(event *PlayerEvent) {
	fmt.Println(fmt.Sprintf("EVENT: %s : %s", event.PlayerId, event.EventId))
}
